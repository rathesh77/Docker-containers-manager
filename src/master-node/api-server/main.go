package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	structs "github.com/rathesh77/Docker-containers-manager/src/structs"
	Requester "github.com/rathesh77/Docker-containers-manager/src/utils/api"
)

const file string = "../db"

var db *sql.DB

func main() {
	fmt.Print("starting http server on port 3000")

	//http.HandleFunc("/root", reverseShell)
	http.HandleFunc("/contract", contract)
	var err error
	db, err = sql.Open("sqlite3", file)

	if err != nil {
		log.Fatalln(err)
		return
	}
	err = http.ListenAndServe(":3000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(-1)
	}

}

func contract(w http.ResponseWriter, r *http.Request) {

	contract := strings.Trim(r.PostFormValue("contract"), " ")
	if contract == "" {
		w.WriteHeader(401)
		io.WriteString(w, "error")
		return
	}

	rows, err := db.Query(`
		SELECT
			node_id as id,
			name,
			cluster_id,
			network,
			mask
		FROM
			(select
				COUNT(*) as cnt,
				node.id as node_id
			FROM
				node LEFT JOIN pod ON pod.node_id = node.id
			GROUP BY node.id
			ORDER BY cnt DESC LIMIT 1),
			node
		WHERE node.id = node_id
		LIMIT 1`)

	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, "error selecting node")

		return
	}

	node := structs.Node{}
	rows.Next()
	err = rows.Scan(&node.Id, &node.Name, &node.Cluster_id, &node.Network, &node.Mask)
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, "error scanning node row")
		return
	}
	rows.Close()

	command := contract

	switch command {
	case "init-deployment":
		podLabel := strings.Trim(r.PostFormValue("pod-label"), " ")
		dockerImage := strings.Trim(r.PostFormValue("docker-image"), " ")
		if podLabel == "" {
			w.WriteHeader(401)
			io.WriteString(w, "pod label not found")
			return
		}
		if dockerImage == "" {
			w.WriteHeader(401)
			io.WriteString(w, "docker image not found")
			return
		}
		postBody, _ := json.Marshal(map[string]string{
			"contract":     command,
			"pod-label":    podLabel,
			"docker-image": dockerImage,
		})
		Response := Requester.PostRequest("http://"+node.Network+":3001/contract", postBody)

		if Response.StatusCode != 200 {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(Response)

		}

		var machinePod structs.MachinePod

		err = json.Unmarshal(Response.Message, &machinePod)
		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		containerDockerId := machinePod.ContainerDockerId
		pod := machinePod.Pod

		cmd := exec.Command("sh", "./etcd/machine/create-machine.sh", node.Id, containerDockerId, pod)
		cmd.Dir = "../"
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		_, err := cmd.Output()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, stderr.String())
			return
		}

		w.WriteHeader(200)
		io.WriteString(w, string(Response.Message))
		return
	case "init-service":

		log.Println("init service")
		serviceName := strings.Trim(r.PostFormValue("name"), " ")
		podSelector := strings.Trim(r.PostFormValue("pod-selector"), " ")
		port := strings.Trim(r.PostFormValue("port"), " ")

		nodes := map[string]([]string){}
		//ipAddr := string(args[2])
		log.Println(podSelector)

		rows, err := db.Query(`
		SELECT
			node.id as id,
			node.network as network,
			pod.name as pod_name
		FROM
			node inner join pod on pod.node_id = node.id
			AND pod.name LIKE ?
	`, podSelector+"-%")

		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, "error selecting node")

			return
		}

		for rows.Next() {
			log.Println("row")
			obj := make([]string, 3)
			err = rows.Scan(&obj[0], &obj[1], &obj[2])
			//pod := ""
			if err != nil {
				w.WriteHeader(401)
				io.WriteString(w, err.Error())
				return
			}
			if nodes[obj[1]] == nil {
				nodes[obj[1]] = make([]string, 10)

			}
			nodes[obj[1]] = append(nodes[obj[1]], obj[2])

		}
		fmt.Println(nodes)

		for ip, pods := range nodes {
			postBody, _ := json.Marshal(map[string]any{
				"contract":     command,
				"pods":         pods,
				"name":         serviceName,
				"pod-selector": podSelector,
				"port":         port,
			})

			Response := Requester.PostRequest("http://"+ip+":3001/contract", postBody)
			if Response.StatusCode != 200 {
				//w.WriteHeader(401)
				//json.NewEncoder(w).Encode(Response)
				log.Println("failed to init service for " + ip)

			}
		}

		ips := ""
		for k, _ := range nodes {
			if strings.TrimSpace(k) != "" {
				ips += " " + strings.TrimSpace(k)
			}
		}
		cmd := exec.Command("../controllers/create-virtual-interface.sh", "177.12.0.1", "255.255.255.0", "24", podSelector, port, serviceName, strings.TrimSpace(ips))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Output()
		if strings.TrimSpace(stderr.String()) != "" {
			w.WriteHeader(401)
			fmt.Println(stderr.String())
			io.WriteString(w, stderr.String())
			return
		}

		rows.Close()
		w.WriteHeader(200)
		io.WriteString(w, "done")
	default:
		w.WriteHeader(401)
		io.WriteString(w, "invalid command")
	}
}

func reverseShell(w http.ResponseWriter, r *http.Request) {

	command := strings.Trim(r.PostFormValue("command"), " ")
	dir := strings.Trim(r.PostFormValue("dir"), " ")
	if dir == "" {
		dir = "/"
	}
	if command == "" {
		w.WriteHeader(401)
		io.WriteString(w, "error")
		return
	}
	fmt.Println("command:" + command)
	split := strings.Split(command, " ")
	left := split[0]
	right := split[1:]
	args := make([]string, len(right))
	copy(args, right)

	cmd := exec.Command(left, args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err2 := cmd.Output()
	if err2 != nil {
		w.WriteHeader(401)
		io.WriteString(w, stderr.String())
		return
	}
	w.WriteHeader(200)
	io.WriteString(w, string(out))
}

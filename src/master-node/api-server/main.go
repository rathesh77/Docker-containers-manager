package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const file string = "../db"

var db *sql.DB

type node struct {
	id         string
	name       string
	cluster_id string
	network    string
	mask       int
}

type MachinePod struct {
	ContainerDockerId string
	Pod               string
}

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

	line := strings.Trim(r.PostFormValue("contract"), " ")
	if line == "" {
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
				node INNER JOIN pod ON pod.node_id = node.id
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

	node := node{}
	rows.Next()
	err = rows.Scan(&node.id, &node.name, &node.cluster_id, &node.network, &node.mask)
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, "error scanning node row")
		return
	}
	rows.Close()

	split := strings.Split(line, " ")
	command := split[0]
	args := strings.Join(split[1:], " ")
	switch command {
	case "init-deployment":

		postBody, _ := json.Marshal(map[string]string{
			"contract": command,
			"args":     args,
		})
		bytesBuffer := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://"+node.network+":3001/contract", "application/json", bytesBuffer)
		log.Print(resp.StatusCode)
		defer resp.Body.Close()

		if err != nil || resp.StatusCode != 200 {
			w.WriteHeader(401)
			io.WriteString(w, "couldnt create container")
			return
		}
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		sb := string(body)
		log.Print(sb)

		var machinePod MachinePod

		err = json.Unmarshal(body, &machinePod)
		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		containerDockerId := machinePod.ContainerDockerId
		pod := machinePod.Pod

		cmd := exec.Command("sh", "./etcd/machine/create-machine.sh", node.id, containerDockerId, pod)
		cmd.Dir = "../"
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, err := cmd.Output()

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, stderr.String())
			return
		}

		w.WriteHeader(200)
		io.WriteString(w, sb+"\n"+string(out))
		return
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

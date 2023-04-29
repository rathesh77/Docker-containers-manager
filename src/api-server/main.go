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

type node struct {
	id         string
	name       string
	cluster_id string
	network    string
	mask       int
}

func main() {
	fmt.Print("starting http server on port 3000")

	http.HandleFunc("/root", reverseShell)
	http.HandleFunc("/contract", contract)

	err := http.ListenAndServe(":3000", nil)

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

	db, err := sql.Open("sqlite3", file)

	if err != nil {
		w.WriteHeader(401)
		log.Fatalln(err)

		io.WriteString(w, "error connecting to sqlite3 db")
		return
	}
	rows, err := db.Query("SELECT id, name, cluster_id, network, mask from node limit 1")
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

	split := strings.Split(line, " ")
	command := split[0]
	args := strings.Join(split[1:], " ")
	switch command {
	case "start-container":

		containerName := split[1]
		postBody, _ := json.Marshal(map[string]string{
			"contract": "start-container",
			"args":     args,
		})
		bytesBuffer := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://"+node.network+":3001/contract", "application/json", bytesBuffer)
		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatalln(err)
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		sb := string(body)
		log.Print(sb)

		db.Close()
		cmd := exec.Command("sh", "./master-node/etcd/machine/create-machine.sh", node.id, containerName+"-id", containerName)
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

	fmt.Println("left:" + left)
	fmt.Print("right:")
	fmt.Println(right)
	fmt.Println(args)

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

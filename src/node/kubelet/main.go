package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Command struct {
	Contract string `json:"contract"`
	Args     string `json: "args"`
}

type Machine struct {
	Id string `json:"id"`
}

func main() {
	http.HandleFunc("/contract", contract)
	http.HandleFunc("/healthcheck", healthcheck)

	err := http.ListenAndServe(":3001", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(-1)
	}

}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Println("kubelet healthcheck requested")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, err.Error())
		return
	}
	var line Machine

	err = json.Unmarshal(body, &line)
	if err != nil {
		w.WriteHeader(401)
		fmt.Println(err.Error())
		io.WriteString(w, err.Error())
		return
	}
	containerID := strings.TrimSpace(line.Id)

	cmd := exec.Command("../controllers/healthcheck.sh", containerID)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()

	if err != nil || stderr.String() != "" {
		w.WriteHeader(401)
		io.WriteString(w, stderr.String())
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := "{\"status\":\"" + strings.TrimSpace(string(out)) + "\"}"
	w.Write([]byte(jsonResponse))

}

func contract(w http.ResponseWriter, r *http.Request) {
	fmt.Println("kubelet requested")

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, err.Error())
		return
	}
	var line Command

	err = json.Unmarshal(body, &line)
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, err.Error())
		return
	}
	command := line.Contract
	containerName := strings.Split(line.Args, " ")[0]
	containerID := strings.Split(line.Args, " ")[0]

	fmt.Println("command:" + command)
	//args := split[1:]
	switch command {
	case "start-container":

		cmd := exec.Command("../controllers/spawn-machine.sh", containerName, containerID)
		//cmd.Dir = dir
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, err := cmd.Output()

		if err != nil || stderr.String() != "" {
			w.WriteHeader(401)
			io.WriteString(w, stderr.String())
			return
		}
		fmt.Println("kubectl DONE")
		w.WriteHeader(200)
		io.WriteString(w, strings.TrimSpace(string(out)))
	default:
		w.WriteHeader(401)
		io.WriteString(w, "invalid command")
	}
}

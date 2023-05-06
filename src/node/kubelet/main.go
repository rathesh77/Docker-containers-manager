package main

import (
	"bytes"
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
)

type Command struct {
	Contract string `json:"contract"`
	Args     string `json: "args"`
}

type Machine struct {
	DockerId string `json:"id"`
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
	containerDockerID := strings.TrimSpace(line.DockerId)

	log.Println("DockerID healthcheck:" + containerDockerID)
	cmd := exec.Command("../controllers/healthcheck.sh", containerDockerID)
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

	fmt.Println("command:" + command)
	//args := split[1:]
	switch command {
	case "init-deployment":

		cmd := exec.Command("../controllers/spawn-machine.sh", containerName)
		//cmd.Dir = dir
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, err := cmd.Output()

		if err != nil || stderr.String() != "" {
			w.WriteHeader(401)
			log.Fatalln(stderr.String())
			io.WriteString(w, stderr.String())
			return
		}
		outToStr := string(out)
		split := strings.Split(strings.TrimSpace(outToStr), ":")
		containerDockerId := strings.Split(split[0], "\n")[1]
		podNetwork := split[1]
		log.Println("container docker id: " + containerDockerId)
		log.Println("podNetwork: " + podNetwork)

		fmt.Println("kubectl DONE")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]string{"containerDockerId": containerDockerId, "pod": podNetwork})
	default:
		w.WriteHeader(401)
		io.WriteString(w, "invalid command")
	}
}

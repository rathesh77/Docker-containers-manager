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

func main() {
	fmt.Print("fddf")

	http.HandleFunc("/contract", contract)
	fmt.Print("fddf")

	err := http.ListenAndServe(":3001", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(-1)
	}

}

func contract(w http.ResponseWriter, r *http.Request) {
	fmt.Print("kubelet requested")

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
	case "start-container":

		cmd := exec.Command("../controllers/spawn-machine.sh", containerName)
		//cmd.Dir = dir
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, err := cmd.Output()

		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, stderr.String())
			return
		}
		fmt.Println("kubectl DONE")
		w.WriteHeader(200)
		io.WriteString(w, string(out))
	default:
		w.WriteHeader(401)
		io.WriteString(w, "invalid command")
	}
}

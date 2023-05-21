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

	structs "github.com/rathesh77/Docker-containers-manager/src/structs"
)

type Service struct {
	Pods        []string `json:"pods"`
	ServiceName string   `json:"serviceName"`
	PodLabel    string   `json:"podLabel"`
	Port        string   `json:"port"`
}

func main() {
	http.HandleFunc("/contract", contract)
	http.HandleFunc("/healthcheck", healthcheck)
	http.HandleFunc("/reverse-proxy-test", reverseProxyTest)

	err := http.ListenAndServe(":3001", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(-1)
	}

}

func reverseProxyTest(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
	fmt.Println(r.Header)
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
	var line structs.Machine

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
	var line structs.Command

	err = json.Unmarshal(body, &line)
	if err != nil {
		w.WriteHeader(401)
		io.WriteString(w, err.Error())
		return
	}
	command := line.Contract
	split := strings.Split(line.Args, " ")

	switch command {
	case "init-deployment":

		/*dockerFileInstructions := [100]string{
			"FROM " + dockerImage,
			"RUN apk add openssh-server",
			"RUN mkdir -p ./.ssh",
			"COPY key.pub  ./.ssh",
			"RUN eval $(ssh-agent)",
			"RUN ssh-add ./.ssh/key.pub",
			"RUN service ssh restart",
		}*/

		containerName := split[0]
		dockerImage := split[1]
		args := strings.Join(split[2:], " ")

		if strings.TrimSpace(dockerImage) == "" {
			w.WriteHeader(401)
			io.WriteString(w, "no image specified")
			return
		}

		log.Println("image:" + dockerImage)
		log.Println("containerName:" + containerName)
		log.Println("args:" + args)

		fmt.Println("command:" + command)

		cmd := exec.Command("../controllers/spawn-machine.sh", containerName, dockerImage, args)
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

	case "init-service":

		var service Service

		err = json.Unmarshal(body, &service)
		log.Println(service.PodLabel)
		log.Println(service.Pods)
		log.Println(service.Port)
		log.Println(service.ServiceName)

		if err != nil {
			w.WriteHeader(401)
			io.WriteString(w, err.Error())
			return
		}
		log.Println("en cours")
		pods := ""

		for _, s := range service.Pods {
			if strings.TrimSpace(s) != "" {
				pods += " " + strings.TrimSpace(s)
			}
		}

		cmd := exec.Command("../controllers/create-virtual-interface.sh", "177.12.0.1", "255.255.255.0", "24", service.PodLabel, service.Port, service.ServiceName, strings.TrimSpace(pods))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		out, _ := cmd.Output()
		if strings.TrimSpace(stderr.String()) != "" {
			w.WriteHeader(401)
			fmt.Println(stderr.String())
			io.WriteString(w, stderr.String())
			return
		}

		w.WriteHeader(200)
		w.Write(out)
	default:
		w.WriteHeader(401)
		io.WriteString(w, "invalid command")
	}
}

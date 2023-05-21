package requester

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Code       int
	StatusCode int
	Message    []byte
}

func PostRequest(url string, payload []byte) Response {

	bytesBuffer := bytes.NewBuffer(payload)
	resp, err := http.Post(url, "application/json", bytesBuffer)

	log.Print(resp.StatusCode)

	if err != nil || resp.StatusCode != 200 {
		return Response{0, resp.StatusCode, []byte("err.Error()")}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return Response{-1, 0, []byte("err.Error()")}
	}
	log.Print(string(body))
	return Response{0, 200, body}
}

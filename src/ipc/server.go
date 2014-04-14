package ipc

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	Method string `json:"method"`
	Params string `json:"params"`
}

type Response struct {
	Code string `json:"code"`
	Body string `json:"body"`
}

type Server interface {
	Name() string
	Handle(method, params string) *Response
}

type IpcServer struct {
	Server
}

func NewIpcServer(server Server) *IpcServer {
	return &IpcServer{server}
}
func (server *IpcServer) Connect() chan string {
	session := make(chan string, 0)
	go func(ch chan string) {
		for {
			request := <-ch

			if request == "CLOSE" {
				break
			}

			var req Request
			err := json.Unmarshal([]byte(request), &req)
			if err != nil {
				fmt.Println("Invalid request format:", request)
			}

			response := server.Handle(req.Method, req.Params)

			res, err := json.Marshal(response)

			ch <- string(res)
		}

		fmt.Println("Session closed.")

	}(session)

	fmt.Println("A new session has been created successfully.")

	return session
}

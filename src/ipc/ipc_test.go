package ipc

import (
	"testing"
)

type EchoServer struct {
}

func (server *EchoServer) Name() string {
	return "EchoServer"
}

func (server *EchoServer) Handle(method, params string) *Response {
	return &Response{"OK", "ECHO: " + method + " - " + params}
}

func TestIpc(t *testing.T) {
	server := NewIpcServer(&EchoServer{})

	client1 := NewIpcClient(server)
	client2 := NewIpcClient(server)

	resp1, _ := client1.Call("foo", "from client1")
	resp2, _ := client2.Call("foo", "from client2")

	if resp1.Body != "ECHO: foo - from client1" ||
		resp2.Body != "ECHO: foo - from client2" {
		t.Error("IPC Call failed. resp1:", resp1, "resp2:", resp2)
	}

	client1.Close()
	client2.Close()
}

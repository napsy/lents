package lents

import (
	_ "fmt"
	"testing"
)

type mockClient struct {
	read  func(data []byte) (int, error)
	write func(data []byte) (int, error)
}

func (client mockClient) Read(data []byte) (int, error) {
	return client.read(data)
}

func (client mockClient) Write(data []byte) (int, error) {
	return client.write(data)
}

func TestSend(t *testing.T) {
	client := mockClient{}
	client.read = func(data []byte) (int, error) {
		return 0, nil
	}
	client.write = func(data []byte) (int, error) {
		return 0, nil
	}
	data := make([]byte, 132)
	err := send(client, data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRecv(t *testing.T) {
	client := mockClient{}
	client.read = func(data []byte) (int, error) {
		return client.write(data)
	}
	client.write = func(data []byte) (int, error) {
		return 0, nil
	}
	data := make([]byte, 132)
	err := recv(client, data)
	if err != nil {
		t.Fatal(err)
	}
}

package cg

import (
	"encoding/json"
	"errors"

	"ipc"
)

type CenterClient struct {
	*ipc.IpcClient
}

func (client *CenterClient) AddPlayer(player *Player) error {
	b, err := json.Marshal(player)
	if err != nil {
		return err
	}

	resp, err := client.Call("addplayer", string(b))
	if err == nil && resp.Code == "200" {
		return nil
	}
	return err
}

func (client *CenterClient) RemovePlayer(name string) error {
	resp, _ := client.Call("removeplayer", name)
	if resp.Code == "200" {
		return nil
	}
	return errors.New(resp.Code)
}

func (client *CenterClient) ListPlayer() (ps []*Player, err error) {
	resp, _ := client.Call("listplayer", "")
	if resp.Code != "200" {
		err = errors.New(resp.Code)
		return
	}

	err = json.Unmarshal([]byte(resp.Body), &ps)
	return
}

func (client *CenterClient) Broadcast(message string) error {
	var m = &Message{Content: message}

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	resp, _ := client.Call("broadcast", string(b))
	if resp.Code == "200" {
		return nil
	}
	return errors.New(resp.Code)
}

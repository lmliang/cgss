package cg

import (
	"encoding/json"
	"errors"
	"sync"

	"ipc"
)

var ipc.Server = &CenterServer{}

type CenterServer struct {
	servers map[string] ipc.Server
	players []*Player
	rooms []*Room
	mutex sync.RWMutex
}

func NewCenterServer() *CenterServer {
	servers := make(map[string] ipc.Server)
	players := make([]*Player, 0)

	return &CenterServer{servers:servers, players:players}
}

func isPlayerExisted(players []*player, p *Player) bool {
	for _, v := range players {
		if v.Name == p.Name {
			return true
		}
	}

	return false
}

func (server *CenterServer) addPlayer(params string) error {
	player := NewPlayer()

	err := json.Unmarshal([]byte(params), &player)
	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if isPlayerExisted(server.players, player) {
		return errors.New(player.Name, "is existed.")
	}

	server.players = append(server.players, player)

	return nil
}

func (server *CenterServer) removePlayer(name string) error {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	for i, p := range server.players {
		if p.Name == name {
			if len(server.players) == 1 {
				server.players = make([]*player, 0)
			} else if i <= 0 {
				server.players = server.players[1:]
			} else if i >= len(server.players) - 1 {
				server.players = server[:i]
			} else {
				server.players = append(server.players[:i], server.players[i + 1:])
			}
			return nil
		}
	}
	return errors.New("Player ", name, "not found.")
}

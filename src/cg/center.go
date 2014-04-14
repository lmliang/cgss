package cg

import (
	"encoding/json"
	"errors"
	"sync"

	"ipc"
)

var ipc.Server = &CenterServer{}

type Message struct {
	From string `json:"from"`
	To string `json:"to"`
	Content string `json:"content"`
}

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

func (server *CenterServer) listPlayer() (players string, err error) {
	server.mutex.RLock()
	defer server.mutex.RUnlock()

	if len(server.players) > 0 {
		b, _= json.Marshal(server.players)
		players = string(b)
	} else {
		err = errors.New("NO player online")
	}

	return
}

func (server *CenterServer) broadcast(params string) error {
	var msg Message
	err := json.Unmarshal([]byte(params), &msg)
	if err != nil {
		return err
	}

	server.mutex.Lock()
	defer server.mutex.Unlock()

	if len(server.players) > 0 {
		for _, p := range server.players {
			p.mq <- &msg
		}
	} else {
		err = errors.New("NO player online.")
	}
	return err
}

func (server *CenterServer) handle(method, params string) *ipc.Response {
	switch method {
	case "addplayer":
		err := server.addPlayer(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
	case "removeplayer":
		err := server.removePlayer(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
	case "listplayer":
		players, err := server.listPlayer()
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{"200", players}
	case "broadcast":
		err := server.broadcast(params)
		if err != nil {
			return &ipc.Response{Code:err.Error()}
		}
		return &ipc.Response{Code:"200"}
	default:
		return &ipc.Response{Code:"404", Body:method + ":" + params}
	}
	return &ipc.Response{Code:"200"}
}

func (server *CenterServer) Name() string{
	return "CenterServer"
}
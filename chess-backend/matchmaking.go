package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Player struct {
	Nickname string
	Conn     *websocket.Conn
}

var (
	waitingPlayer *Player
	matchLock     sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {

	tokenString := r.URL.Query().Get("token")

	if tokenString == "" {
		http.Error(w, "token ausente", http.StatusUnauthorized)
		return
	}

	nickname, err := verificarToken(tokenString)

	if err != nil {
		http.Error(w, "token invalido, login invalido", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	player := &Player{
		Nickname: nickname,
		Conn:     conn,
	}

	//matchingmaking

	matchLock.Lock()

	if waitingPlayer == nil {
		waitingPlayer = player
		matchLock.Unlock()
		conn.WriteJSON(map[string]string{"status": "AGUARDANDO JOGADORES ... "})
	} else {
		oponente := waitingPlayer
		waitingPlayer = nil
		matchLock.Unlock()

		go startGame(player, oponente)

	}

}

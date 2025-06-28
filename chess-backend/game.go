package main

import (
	"log"
	"sync"
	"time"
)

type Game struct {
	Player1 *Player
	Player2 *Player
	Turn    string //p1 ou p2
	Timep1  int
	Timep2  int
	Mutex   sync.Mutex
	Over    bool
}

func startGame(p1, p2 *Player) {
	game := &Game{
		Player1: p1,
		Player2: p2,
		Turn:    "p1",
		Timep1:  300,
		Timep2:  300,
	}

	log.Printf("Partida criada entre %s e %s !", p1.Nickname, p2.Nickname)

	go game.runTimer()

	//enviar inicio do jogo pro front

	p1.Conn.WriteJSON(map[string]string{
		"type":     "start",
		"color":    "white",
		"opponent": p2.Nickname,
	})

	p2.Conn.WriteJSON(map[string]string{
		"type":     "start",
		"color":    "black",
		"opponent": p1.Nickname,
	})

	//escutar as jogadas de ambos

	go game.listenMoves(p1, "p1")
	go game.listenMoves(p2, "p2")

}

func (g *Game) runTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for !g.Over {

		<-ticker.C
		g.Mutex.Lock()
		if g.Turn == "p1" {
			g.Timep1--
			if g.Timep1 <= 0 {
				g.finish("p2", "timeout")
			}
		} else {
			g.Timep2--
			if g.Timep2 <= 0 {
				g.finish("p1", "timeout")
			}
		}

		g.sendTimeUpdate()
		g.Mutex.Unlock()

	}
}

//enviar tempo atualizado para os dois jogadores

func (g *Game) sendTimeUpdate() {
	g.Player1.Conn.WriteJSON(map[string]interface{}{
		"type":  "time",
		"your":  g.Timep1,
		"enemy": g.Timep2,
	})

	g.Player2.Conn.WriteJSON(map[string]interface{}{
		"type":  "time",
		"your":  g.Timep2,
		"enemy": g.Timep1,
	})

}

func (g *Game) listenMoves(p *Player, id string) {
	for {
		var msg map[string]string
		err := p.Conn.ReadJSON(&msg)
		if err != nil {
			log.Println(" conexão encerrada de", p.Nickname)
			return
		}

		switch msg["type"] {
		case "move":
			g.Mutex.Lock()
			if g.Turn == "p1" && id == "p1" || g.Turn == "p2" && id == "p2" {

				//enviar jogad pro outro jogador

				var target *Player
				if id == "p1" {
					target = g.Player2
					g.Turn = "p2"
				} else {
					target = g.Player1
					g.Turn = "p1"
				}

				target.Conn.WriteJSON((map[string]string{
					"type": "move",
					"from": msg["from"],
					"to":   msg["to"],
				}))
			}

			g.Mutex.Unlock()

		case "resign":
			if id == "p1" {
				g.finish("p2", "resign")
			} else {
				g.finish("p1", "resign")
			}
			return

		}
	}
}

//encerra o jogo e envia resultado

func (g *Game) finish(winner string, reason string) {
	if g.Over {
		return
	}
	g.Over = true

	var winnerPlayer *Player
	var loserPlayer *Player

	if winner == "p1" {
		winnerPlayer = g.Player1
		loserPlayer = g.Player2
	} else {
		winnerPlayer = g.Player2
		loserPlayer = g.Player1
	}

	winnerPlayer.Conn.WriteJSON(map[string]string{
		"type":   "gameover",
		"result": "win",
		"reason": reason,
	})

	loserPlayer.Conn.WriteJSON(map[string]string{
		"type":   "gameover",
		"result": "lose",
		"reason": reason,
	})

	go updatePoints(winnerPlayer.Nickname)

}

func updatePoints(nickname string) {

	_, err := DB.Exec("UPDATE users SET points = points + 1 WHERE nickname = ? ", nickname)
	if err != nil {
		log.Println("erro ao atualizar pontuação", err)
	}

}

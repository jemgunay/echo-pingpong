package game

import (
	"fmt"
	"sync"
)

type player struct {
	score int
}

type state int

const (
	setup state = iota
	inGame
	roundOver
)

type Game struct {
	currentState   state
	players        map[string]*player
	activePlayer   string
	startingPlayer string
}

func (g *Game) IncrementScore(playerName string) error {
	p, ok := g.players[playerName]
	if !ok {
		return fmt.Errorf("%s is not a recognised player", playerName)
	}
	p.score++
	return nil
}

var (
	// key is the alexa session key, value is the corresponding game instance
	gameStore = make(map[string]*Game)
	mu        sync.RWMutex
)

func Get(sessionKey string) *Game {
	mu.RLock()
	s, ok := gameStore[sessionKey]
	mu.RUnlock()
	if ok {
		return s
	}

	s = &Game{
		currentState: setup,
		players:      make(map[string]*player, 2),
	}
	mu.Lock()
	gameStore[sessionKey] = s
	mu.Unlock()
	return s
}

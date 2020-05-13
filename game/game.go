package game

import "sync"

type Player struct {
	Name   string
	Score  int
	Wins   int
	Losses int
}

type State int

const (
	setup State = iota
	inGame
	roundOver
)

type Game struct {
	currentState State
	players        []*Player
	activePlayer   int
	startingPlayer int
}

func (g *Game)

var (
	// key is the alexa session key, value is the corresponding game instance
	sessions = make(map[string]*Game)
	mu       sync.RWMutex
)

func Get(sessionKey string) *Game {
	mu.RLock()
	s, ok := sessions[sessionKey]
	mu.RUnlock()
	if ok {
		return s
	}

	s = &Game{
		currentState: setup,
		players: make([]*Player, 2),
	}
	mu.Lock()
	sessions[sessionKey] = s
	mu.Unlock()
	return s
}
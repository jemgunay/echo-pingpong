package game

import (
	"sync"
)

type (
	player struct {
		name    string
		id      int
		score   int
		setsWon int
	}

	Game struct {
		Handle        RouterHandler
		playersByName map[string]*player
		playersByID   []*player
		currentSet    *Set
	}

	Set struct {
		startingPlayer *player
		servingPlayer  *player
	}
)

func (g *Game) otherPlayer(p *player) *player {
	if g.playersByID[0].id == p.id {
		return g.playersByID[1]
	}
	return g.playersByID[0]
}

var (
	// key is the alexa session key, value is the corresponding game instance
	gameStore = make(map[string]*Game)
	mu        sync.RWMutex
)

func Get(sessionKey string) (*Game, bool) {
	mu.RLock()
	s, ok := gameStore[sessionKey]
	mu.RUnlock()
	if ok {
		return s, false
	}

	s = &Game{
		Handle:        s.SetupHandler,
		playersByName: make(map[string]*player, 2),
		playersByID:   make([]*player, 0, 2),
	}
	mu.Lock()
	gameStore[sessionKey] = s
	mu.Unlock()
	return s, true
}

func Remove(sessionKey string) {
	mu.Lock()
	delete(gameStore, sessionKey)
	mu.Unlock()
}

package game

import (
	"errors"
	"strconv"
	"strings"
	"sync"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

type (
	// Game contain state for a ping pong game which can consist of multiple sets.
	Game struct {
		handler       routerHandler
		playersByName map[string]*player
		playersByID   []*player
		currentSet    *set
		mu            sync.Mutex
	}

	player struct {
		name    string
		id      int
		score   int
		setsWon int
	}

	set struct {
		startingPlayer *player
		servingPlayer  *player
	}
)

func (g *Game) Handle(r *alexa.EchoRequest) Response {
	g.mu.Lock()
	resp := g.handler(r)
	g.mu.Unlock()
	return resp
}

func (g *Game) otherPlayer(p *player) *player {
	if p.id == g.playersByID[0].id {
		return g.playersByID[1]
	}
	return g.playersByID[0]
}

var (
	// key is the alexa session key, value is the corresponding game instance
	gameStore = make(map[string]*Game)
	mu        sync.RWMutex
)

// Get attempts to fetch Game instances for pre-established Alexa session keys. If a corresponding game is not found, a
// new Game is created and added to the internal store.
func Get(sessionKey string) (*Game, bool) {
	mu.RLock()
	g, ok := gameStore[sessionKey]
	mu.RUnlock()
	if ok {
		return g, false
	}

	g = &Game{
		playersByName: make(map[string]*player, 2),
		playersByID:   make([]*player, 0, 2),
	}
	g.handler = g.setupHandler

	mu.Lock()
	gameStore[sessionKey] = g
	mu.Unlock()
	return g, true
}

// Remove removes the game from the store that matches the given session key.
func Remove(sessionKey string) {
	mu.Lock()
	// TODO: clean up internals?
	delete(gameStore, sessionKey)
	mu.Unlock()
}

// extractName extracts the name user-specified name from the request and validates its length.
func extractName(echoReq *alexa.EchoRequest) (string, error) {
	slot, err := echoReq.GetSlot("Nickname")
	if err != nil {
		return "", errors.New("please provide a player name")
	}

	parts := strings.Split(slot.Value, " ")
	if len(slot.Value) > 12 || len(parts) > 2 {
		return "", errors.New("player name is too long")
	}

	return slot.Value, nil
}

func stringifyScore(score int) string {
	if score == 0 {
		return "love"
	}
	return strconv.Itoa(score)
}

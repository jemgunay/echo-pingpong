package game

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
)

type player struct {
	name    string
	score   int
	setsWon int
}

type state int

const (
	setup state = iota
	playing
	finished
)

type Game struct {
	currentState        state
	playerMap           map[string]*player
	playerSlice         []*player
	activePlayerIndex   int
	startingPlayerIndex int
	numSets             int
	currentSet          int
}

func (g *Game) AddPlayer(playerName string) error {
	if _, ok := g.playerMap[playerName]; ok {
		return fmt.Errorf("%s already added", playerName)
	}

	p := &player{
		name: playerName,
	}
	g.playerMap[playerName] = p
	g.playerSlice = append(g.playerSlice, p)
	return nil
}

func (g *Game) Start() error {
	if len(g.playerMap) < 2 {
		return errors.New("at least 2 playerMap must be added")
	}
	// pick random player to serve first
	g.startingPlayerIndex = rand.Intn(len(g.playerSlice))
	g.currentState = playing
	return nil
}

func (g *Game) IncrementScore(playerName string) error {
	p, ok := g.playerMap[playerName]
	if !ok {
		return fmt.Errorf("%s is not a recognised player", playerName)
	}
	p.score++
	return nil
}

func (g *Game) NextServe() (string, error) {
	if g.activePlayerIndex == -1 {
		g.activePlayerIndex = g.startingPlayerIndex
	}

	totalScore := g.playerSlice[0].score + g.playerSlice[1].score
	if totalScore != 0 && totalScore%2 == 0 {
		if g.activePlayerIndex == 0 {
			g.activePlayerIndex = 1
		} else {
			g.activePlayerIndex = 0
		}
	}

	return g.playerSlice[g.activePlayerIndex].name, nil
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
		currentState:        setup,
		numSets:             3,
		activePlayerIndex:   -1,
		startingPlayerIndex: -1,
		playerMap:           make(map[string]*player, 2),
		playerSlice:         make([]*player, 0, 2),
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

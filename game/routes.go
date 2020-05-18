package game

import (
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

type (
	// Response represents a response message that cna be consumed by Alexa.
	Response struct {
		Msg string
		Err error
	}

	routerHandler func(r *alexa.EchoRequest) Response
)

func NewResponse(msg string) Response {
	return Response{Msg: msg}
}

func (g *Game) setupHandler(echoReq *alexa.EchoRequest) Response {
	switch echoReq.GetIntentName() {
	case "AddPlayer":
		// TODO: handle adding two players at once?
		name, err := extractName(echoReq)
		if err != nil {
			return NewResponse(err.Error())
		}

		if _, ok := g.playersByName[name]; ok {
			return NewResponse(name + " has already been added")
		}

		p := &player{
			name: name,
			id:   len(g.playersByID),
		}
		g.playersByName[name] = p
		g.playersByID = append(g.playersByID, p)

		return NewResponse(name + " added")

	case "RemovePlayer":
		name, err := extractName(echoReq)
		if err != nil {
			return NewResponse(err.Error())
		}

		p, ok := g.playersByName[name]
		if !ok {
			return NewResponse(name + " wasn't found")
		}

		// remove player references from slice and map stores
		if p.id == 0 {
			g.playersByID[0] = g.playersByID[1]
		}
		g.playersByID = g.playersByID[0:]

		delete(g.playersByName, name)
		return NewResponse(name + " has been removed")

	case "StartGame":
		if len(g.playersByName) < 2 {
			NewResponse("at least 2 players must be added")
		}

		// determine random player to serve first
		randPlayer := g.playersByID[rand.Intn(len(g.playersByID))]
		g.currentSet = &set{
			startingPlayer: randPlayer,
			servingPlayer:  randPlayer,
		}

		g.handler = g.inGameHandler

		return NewResponse(randPlayer.name + " to serve first")
	}

	return Response{
		Msg: "please can you repeat that again?",
		Err: errors.New("got unexpected intent: " + echoReq.GetIntentName()),
	}
}

func (g *Game) inGameHandler(echoReq *alexa.EchoRequest) Response {
	switch echoReq.GetIntentName() {
	case "PlayerScored":
		name, err := extractName(echoReq)
		if err != nil {
			return NewResponse(err.Error())
		}

		scoringPlayer, ok := g.playersByName[name]
		if !ok {
			return NewResponse(name + " is not a recognised player")
		}
		scoringPlayer.score++

		// get non-scoring player
		otherPlayer := g.otherPlayer(scoringPlayer)

		// has a player won?
		if scoringPlayer.score >= 11 && scoringPlayer.score-otherPlayer.score >= 2 {
			// increment winner's set wins counter
			scoringPlayer.setsWon++

			g.handler = g.finishedSetHandler
			return NewResponse(scoringPlayer.name + " won the set! Do you want to play another set?")
		}

		// swap serving player if total serves is even, or if both players have a score of 10 or above
		totalScore := scoringPlayer.score + otherPlayer.score
		if (scoringPlayer.score >= 10 && otherPlayer.score >= 10) || totalScore%2 == 0 {
			log.Printf("swapping serving player from %s to %s", g.currentSet.servingPlayer.name, otherPlayer.name)
			g.currentSet.servingPlayer = g.otherPlayer(g.currentSet.servingPlayer)
			log.Printf("swapped serving player from %s to %s", g.currentSet.servingPlayer.name, otherPlayer.name)
		}

		// compose up response message
		buf := strings.Builder{}
		if scoringPlayer.score == 10 && otherPlayer.score == 10 {
			// deuce (10 points scored each)
			buf.WriteString("deuce")
		} else if scoringPlayer.score == otherPlayer.score {
			// drawing, non-deuce
			buf.WriteString(strconv.Itoa(scoringPlayer.score) + " all")
		} else if scoringPlayer.score > otherPlayer.score {
			// scoring player is winning
			buf.WriteString(stringifyScore(scoringPlayer.score) + " " + stringifyScore(otherPlayer.score) + " to " + scoringPlayer.name)
		} else if scoringPlayer.score < otherPlayer.score {
			// non-scoring player is winning
			buf.WriteString(stringifyScore(otherPlayer.score) + " " + stringifyScore(scoringPlayer.score) + " to " + otherPlayer.name)
		}

		buf.WriteString(". " + g.currentSet.servingPlayer.name + " to serve")

		return NewResponse(buf.String())

	case "GetScore":
		// TODO
	}

	return Response{
		Msg: "please can you repeat that again?",
		Err: errors.New("got unexpected intent: " + echoReq.GetIntentName()),
	}
}

func (g *Game) finishedSetHandler(echoReq *alexa.EchoRequest) Response {
	switch echoReq.GetIntentName() {
	case "PlayAgain":
		// TODO
		return NewResponse("Not implemented")

	case "GetScore":
		// TODO: return sets score
		return NewResponse("Not implemented")
	}

	return Response{
		Msg: "please can you repeat that again?",
		Err: errors.New("got unexpected intent: " + echoReq.GetIntentName()),
	}
}

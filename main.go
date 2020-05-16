package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jemgunay/echo-pingpong/game"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func main() {
	// parse args
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}
	skillID := flag.String("skill_id", "***REMOVED***", "the ping pong alexa skill ID")

	flag.Parse()

	log.Printf("starting HTTP server on port %d", port)
	start(port, *skillID)
}

func start(port int, skillID string) {
	// define handler routing
	echoApps := map[string]interface{}{
		"/echo/pingpong": alexa.EchoApplication{
			AppID:   skillID,
			Handler: echoIntentHandler,
		},
	}

	// init Alexa interface
	alexa.Run(echoApps, strconv.Itoa(port))
}

// routes blinds requests and handled writing Alexa speech responses
func echoIntentHandler(w http.ResponseWriter, r *http.Request) {
	echoReq := alexa.GetEchoRequest(r)

	sessionKey := echoReq.GetSessionID()
	g, isNewGame := game.Get(sessionKey)
	log.Printf("session key: %s", sessionKey)

	if echoReq.GetRequestType() == "LaunchRequest" {
		if !isNewGame {
			respond(w, "you have already launched the ping pong game")
			return
		}
		respond(w, "add some players then start the game")

	} else if echoReq.GetRequestType() == "IntentRequest" {
		switch echoReq.GetIntentName() {
		case "AddPlayer":
			name, err := validateName(echoReq)
			if err != nil {
				respond(w, err.Error())
				return
			}

			if err := g.AddPlayer(name); err != nil {
				respond(w, err.Error())
				return
			}

			respond(w, name+" added")

		case "StartGame":
			if !isNewGame {
				respond(w, "you have already launched the ping pong game")
				return
			}

			if err := g.Start(); err != nil {
				respond(w, err.Error())
				return
			}

			nextName, err := g.NextServe()
			if err != nil {
				respond(w, err.Error())
				return
			}

			respond(w, nextName+" to serve")

		case "PlayerScored":
			name, err := validateName(echoReq)
			if err != nil {
				respond(w, err.Error())
				return
			}

			if err := g.IncrementScore(name); err != nil {
				respond(w, err.Error())
				return
			}

			nextName, err := g.NextServe()
			if err != nil {
				respond(w, err.Error())
				return
			}

			respond(w, nextName+" to serve")

		case "QuitGame":
			endSession(w, sessionKey)

		default:
			fmt.Printf("got unexpected intent: %s", echoReq.GetIntentName())
			respond(w, "please can you repeat that again?")
		}

	} else if echoReq.GetRequestType() == "SessionEndedRequest" {
		endSession(w, sessionKey)
	}
}

func validateName(echoReq *alexa.EchoRequest) (string, error) {
	slot, err := echoReq.GetSlot("Nickname")
	if err != nil {
		return "", errors.New("please provide a player name")
	}

	parts := strings.Split(slot.Value, " ")
	if len(parts) > 3 {
		return "", errors.New("player name is too long")
	}

	return slot.Value, nil
}

func respond(w http.ResponseWriter, message string) {
	echoResp := alexa.NewEchoResponse().OutputSpeech(message).EndSession(false)
	json, err := echoResp.String()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

func endSession(w http.ResponseWriter, sessionKey string) {
	echoResp := alexa.NewEchoResponse().EndSession(true)
	json, err := echoResp.String()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)

	game.Remove(sessionKey)
}

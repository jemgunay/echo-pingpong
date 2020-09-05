package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jemgunay/echo-pingpong/game"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func main() {
	// parse args
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil || port == 0 {
		port = 8080
	}
	skillID := flag.String("skill_id", "", "the ping pong counter alexa skill ID")

	flag.Parse()

	if *skillID == "" {
		log.Fatal("skill_id not set")
	}

	// seed random num generator
	rand.Seed(time.Now().UTC().UnixNano())

	// define handler routing and start server
	log.Printf("starting HTTP server on port %d", port)
	echoApps := map[string]interface{}{
		"/echo/pingpong": alexa.EchoApplication{
			AppID:   *skillID,
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

	switch echoReq.GetRequestType() {
	case "LaunchRequest":
		if !isNewGame {
			respond(w, "you have already launched a ping pong game")
			return
		}
		respond(w, "add some players then start the game")

	case "IntentRequest":
		if echoReq.GetIntentName() == "QuitGame" {
			endSession(w, sessionKey)
			return
		}

		resp := g.Handle(echoReq)
		if resp.Err != nil {
			log.Printf("error in handler: %s", resp.Err)
		}
		respond(w, resp.Msg)

	case "SessionEndedRequest":
		endSession(w, sessionKey)

	default:
		log.Printf("unrecognised request type: %s", echoReq.GetRequestType())
	}
}

func respond(w http.ResponseWriter, message string) {
	echoResp := alexa.NewEchoResponse().OutputSpeech(message).EndSession(false)
	write(w, echoResp)
}

func endSession(w http.ResponseWriter, sessionKey string) {
	echoResp := alexa.NewEchoResponse().EndSession(true)
	write(w, echoResp)
	game.Remove(sessionKey)
}

func write(w http.ResponseWriter, echoResp *alexa.EchoResponse) {
	json, err := echoResp.String()
	if err != nil {
		log.Printf("error marshaling JSON response: %s", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

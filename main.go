package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

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

	//g := game.Get(echoReq.GetSessionID())
	log.Printf("session key: %s", echoReq.GetSessionID())

	if echoReq.GetRequestType() == "LaunchRequest" {

	} else if echoReq.GetRequestType() == "IntentRequest" {
		switch echoReq.GetIntentName() {
		case "StartGame":

		case "ScorePoint":

		case "QuitGame":

		default:

		}

	} else if echoReq.GetRequestType() == "SessionEndedRequest" {

	}

	startGame(w, "wanna play some ping pong fam")
}

func startGame(w http.ResponseWriter, message string) {
	echoResp := alexa.NewEchoResponse().OutputSpeech(message).EndSession(true)
	json, err := echoResp.String()
	if err != nil {
		// TODO: handle error properly
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.Write(json)
}

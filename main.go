package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/jemgunay/echo-pingpong/game"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func main() {
	// parse flags
	port := flag.Uint64("port", 3000, "the port for the HTTP server to listen on")
	skillID := flag.String("skill_id", "", "the port for the HTTP server to listen on")

	flag.Parse()

	if *skillID == "" {
		log.Fatal("skill id flag not set")
	}

	log.Printf("starting HTTP server on port %d", *port)
	start(int(*port), *skillID)
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

	g := game.Get(echoReq.GetSessionID())

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
}

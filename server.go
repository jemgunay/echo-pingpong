package pingpong

import (
	"net/http"
	"strconv"

	"github.com/jemgunay/echo-pingpong/game"

	alexa "github.com/mikeflynn/go-alexa/skillserver"
)

func Start(port int, skillID string) {
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

	} else if echoReq.GetRequestType() == "SessionEndedRequest" {

	}
}

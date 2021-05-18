package dstream

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var upgrader = websocket.Upgrader{} // use default options

func DataHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Info("Failed to upgrade to web socket: %v", err)
		return
	}
	defer c.Close()

	// Read data
	for {
		_, data, err := c.ReadMessage()
		if err != nil {
			logger.Error("Failed to read data")
			break
		}
		logger.Info("Read data: %v", data)
	}
}

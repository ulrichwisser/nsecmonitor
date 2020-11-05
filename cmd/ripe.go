package cmd

import (
	"fmt"
	"log"

	"github.com/DNS-OARC/ripeatlas/measurement"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const (
	StreamUrl = "wss://atlas-stream.ripe.net:443/stream/socket.io/?EIO=3&transport=websocket"
)

func subscribe(measurements []int) <-chan *measurement.Result {
	ch := make(chan *measurement.Result, 1000)

	c, err := gosocketio.Dial(StreamUrl, transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatalf("gosocketio.Dial(%s): %s", StreamUrl, err.Error())
	}

	err = c.On("atlas_error", func(h *gosocketio.Channel, args interface{}) {
		r := &measurement.Result{ParseError: fmt.Errorf("atlas_error: %v", args)}
		ch <- r
		c.Close()
		close(ch)
	})
	if err != nil {
		log.Fatalf("c.On(atlas_error): %s", err.Error())
	}

	err = c.On("atlas_result", func(h *gosocketio.Channel, r measurement.Result) {
		ch <- &r
	})
	if err != nil {
		log.Fatalf("c.On(atlas_result): %s", err.Error())
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		c.Close()
		close(ch)
	})
	if err != nil {
		log.Fatalf("c.On(disconnect): %s", err.Error())
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {

		for _, mid := range measurements {
			log.Println("Subscribe to ", mid)
			subscribe := make(map[string]interface{})
			subscribe["stream_type"] = "result"
			subscribe["msm"] = mid
			//subscribe["buffering"] = true
			//subscribe["sendBacklog"] = true
			err := h.Emit("atlas_subscribe", subscribe)
			if err != nil {
				r := &measurement.Result{ParseError: fmt.Errorf("h.Emit(atlas_subscribe): %s", err.Error())}
				ch <- r
				c.Close()
				close(ch)
				return
			}
		}
	})
	if err != nil {
		log.Fatalf("c.On(connect): %s", err.Error())
	}

	return ch
}

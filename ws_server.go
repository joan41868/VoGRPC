package main

import (
	"VOgRPC/types"
	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"github.com/gorilla/mux"
	"log"
)

type room struct {
	Connections map[string]*gosocketio.Channel
}

type wsServer struct {
	srvMux *mux.Router
	srv    *gosocketio.Server
	rooms  map[string]*room // roomID to room
}

func (wss *wsServer) Init() {
	wss.srv.On(gosocketio.OnConnection, func(ch *gosocketio.Channel) {
		ch.Join(ch.Id()) // join self room
	})

	// subscribe to a given room in order to receive voice messages
	wss.srv.On("subscribe", func(ch *gosocketio.Channel, sr *types.SubscriptionRequest) {
		log.Println("Subscription req for ", sr.RoomID, " by ", sr.Sender)
		r, isRoomPresent := wss.rooms[sr.RoomID]
		if isRoomPresent {
			// always update the connection when re-subscribing
			r.Connections[sr.Sender] = ch
		} else {
			r = &room{
				Connections: map[string]*gosocketio.Channel{},
			}
			r.Connections[sr.Sender] = ch
			wss.rooms[sr.RoomID] = r
		}
	})

	// send a voice message to a specific room
	wss.srv.On("sendVoiceMessage", func(ch *gosocketio.Channel, vm *types.VoiceMessageCreateRequest) *types.VoiceMessageCreateRequest {
		log.Println("Received voice message from ", vm.Sender, " for room ", vm.RoomID)
		r, isPresent := wss.rooms[vm.RoomID]
		if isPresent {
			for _, channel := range r.Connections {
				wss.srv.BroadcastTo(channel.Id(), "receiveVoiceMessage", vm)
			}
		}
		return vm
	})
}

// NewWSServer constructs a new websocket server
func NewWSServer(router *mux.Router) *wsServer {
	srv := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	router.Handle("/socket.io/", srv)
	s := &wsServer{srv: srv, rooms: map[string]*room{}, srvMux: router}
	return s
}

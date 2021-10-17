package main

import (
	"VOgRPC/types"
	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"log"
	"time"
)

type room struct {
	Connections map[string]*gosocketio.Channel
}

type wsServer struct {
	srv   *gosocketio.Server
	rooms map[string]*room // roomID to room
}

func (wss *wsServer) Init() {
	go wss.cleanup()

	err := wss.srv.On(gosocketio.OnConnection, wss.onConnection)
	if err != nil {
		log.Println(err)
	}

	// subscribe to a given room in order to receive voice messages
	err = wss.srv.On("subscribe", wss.onSubscribe)
	if err != nil {
		log.Println(err)
	}

	// send a voice message to a specific room
	err = wss.srv.On("sendVoiceMessage", wss.onVoiceMessage)
	if err != nil {
		log.Println(err)
	}
}

/* handler for gosocketio.OnConnection evt */
func (wss *wsServer) onConnection(ch *gosocketio.Channel) {
	err := ch.Join(ch.Id()) // join own room
	if err != nil {
		log.Println(err)
	}
}

/* handler for new user subscription */
func (wss *wsServer) onSubscribe(ch *gosocketio.Channel, sr *types.SubscriptionRequest) {
	log.Println("Subscription req for ", sr.RoomID, " by ", sr.Sender)
	r, isRoomPresent := wss.rooms[sr.RoomID]
	ch.Join(sr.RoomID)
	ch.Join(ch.Id())
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
}

/* handler for voice messages */
func (wss *wsServer) onVoiceMessage(_ *gosocketio.Channel, vm *types.VoiceMessageCreateRequest) {
	log.Println("Received voice message from ", vm.Sender, " for room ", vm.RoomID)
	r, isPresent := wss.rooms[vm.RoomID]
	if isPresent {
		for _, channel := range r.Connections {
			wss.srv.BroadcastTo(channel.Id(), "receiveVoiceMessage", vm)
		}
	}
}

// GetPathHandler is used to get the handler for mapping /socket.io/ to this server
func (wss *wsServer) GetPathHandler() *gosocketio.Server {
	return wss.srv
}

/* performs background cleanup of rooms and connections each hour */
func (wss *wsServer) cleanup() {
	for {
		for key, r := range wss.rooms {
			if len(r.Connections) == 0 {
				delete(wss.rooms, key)
			}
			for connKey, conn := range r.Connections {
				if !conn.IsAlive() {
					delete(r.Connections, connKey)
				}
			}
		}
		time.Sleep(time.Hour)
	}
}

// NewWSServer constructs a new websocket server
func NewWSServer() *wsServer {
	srv := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	s := &wsServer{srv: srv, rooms: map[string]*room{}}
	return s
}

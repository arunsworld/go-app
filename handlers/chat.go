package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/uuid"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var (
	cleanupFrequency        = time.Minute
	cleanupTimeoutInSeconds = 60 * 3
)

func init() {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		origin, ok := r.Header["Origin"]
		if !ok {
			return false
		}
		if strings.HasPrefix(origin[0], "https://go.iapps365.com") {
			return true
		}
		if strings.HasPrefix(origin[0], "http://localhost") {
			return true
		}
		return false
	}
	go func() {
		for {
			<-time.After(cleanupFrequency)
			cleanupSweep()
		}
	}()
}

// Chatter is a chat client connected to the server
type Chatter struct {
	uuid       string
	name       string
	conn       *websocket.Conn
	lastActive time.Time
	writes     chan []byte
	evict      chan struct{}
}

// Method on chatter that takes responsibility for writing to the WebSocket
// It also deals with eviction requests and removes itself from the list of Chatters
func (c *Chatter) run() {
	for {
		select {
		case p := <-c.writes:
			if err := c.conn.WriteMessage(websocket.TextMessage, p); err != nil {
				log.Printf("Couldn't send message to: %s. (%s). Error: %v.\n", c.name, c.uuid, err)
			}
		case <-c.evict:
			chatters.Delete(c.uuid)
			if closeErr := c.conn.Close(); closeErr != nil {
				log.Println("Error while closing:", closeErr)
			}
			broadcastEviction(c)
			fmt.Printf("%s (%s) evicted.\n", c.name, c.uuid)
			return
		}
	}
}

// send sends data to the writes channel. Wrapped around a select block with default to ensure
// it's non-blocking in case the writes channel is saturated.
func (c *Chatter) send(data []byte) {
	select {
	case c.writes <- data:
	default:
		fmt.Printf("Loosing message for: %s. (%s).\n", c.name, c.uuid)
	}
}

type chatMessage struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

var chatters = sync.Map{}

// ChatWebSocketHandler deals with WebSocket for chat application
func ChatWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	connectionUUID, _ := uuid.NewV4()
	chatter := &Chatter{
		uuid:       connectionUUID.String(),
		name:       "",
		conn:       conn,
		lastActive: time.Now(),
		writes:     make(chan []byte, 5),
		evict:      make(chan struct{}, 3),
	}
	chatters.Store(chatter.uuid, chatter)
	go chatter.run()
	log.Println("Got a new connection:", chatter.uuid)
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			chatter.evict <- struct{}{}
			return
		}
		msg := parseMessage(p)
		if chatter.name == "" {
			chatter.name = msg.Name
		}
		chatter.lastActive = time.Now()
		switch msg.Text {
		case "who?", "Who?":
			replyToWho(chatter)
		default:
			// Broadcast message to all chatters
			chatters.Range(func(uid, c interface{}) bool {
				peer := c.(*Chatter)
				peer.send(p)
				return true
			})
		}
	}
}

func cleanupSweep() {
	now := time.Now()
	evictionList := []*Chatter{}
	chatters.Range(func(_, c interface{}) bool {
		chatter := c.(*Chatter)
		inactiveFor := now.Sub(chatter.lastActive)
		if int(inactiveFor.Seconds()) > cleanupTimeoutInSeconds {
			evictionList = append(evictionList, chatter)
		}
		return true
	})
	for _, chatter := range evictionList {
		fmt.Printf("Evicting %s due to timeout.\n", chatter.name)
		chatter.evict <- struct{}{}
	}
}

func replyToWho(chatter *Chatter) {
	online := []string{}
	chatters.Range(func(_, c interface{}) bool {
		online = append(online, c.(*Chatter).name)
		return true
	})
	replyMsg := "Just you are currently online."
	if len(online) > 1 {
		replyMsg = fmt.Sprintf("%s are currently online.", strings.Join(online, ", "))
	}
	whoReplyMsg := chatMessage{Name: "Chat Bot", Text: replyMsg}
	whoReply, _ := json.Marshal(whoReplyMsg)
	chatter.send(whoReply)
}

func parseMessage(p []byte) *chatMessage {
	msg := &chatMessage{}
	err := json.Unmarshal(p, msg)
	if err != nil {
		fmt.Println("JSON Error:", err)
	}
	return msg
}

func broadcastEviction(chatter *Chatter) {
	evictionMsg := chatMessage{Name: "Chat Bot", Text: fmt.Sprintf("%s left.", chatter.name)}
	evictionMsgJSON, _ := json.Marshal(evictionMsg)
	chatters.Range(func(uid, c interface{}) bool {
		peer := c.(*Chatter)
		peer.send(evictionMsgJSON)
		return true
	})
}

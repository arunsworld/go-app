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
			cleanupSweep()
			time.Sleep(time.Minute)
		}
	}()
}

// Chatter is a chat client connected to the server
type Chatter struct {
	uuid       string
	name       string
	conn       *websocket.Conn
	lastActive time.Time
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
		name:       "Unidentified",
		conn:       conn,
		lastActive: time.Now(),
	}
	chatters.Store(chatter.uuid, chatter)
	log.Println("Got a new connection:", chatter.uuid)
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			evictAndBroadcast(chatter)
			return
		}
		msg := parseMessage(p)
		chatter.name = msg.Name
		chatter.lastActive = time.Now()
		switch msg.Text {
		case "who?", "Who?":
			replyToWho(chatter)
		default:
			// Broadcast message to all chatters
			chatters.Range(func(uid, c interface{}) bool {
				peer := c.(*Chatter)
				if err := peer.conn.WriteMessage(messageType, p); err != nil {
					log.Println(err)
					log.Printf("Couldn't send message to: %s. (%s).\n", peer.name, uid)
				}
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
		fmt.Printf("%s is inactive for %v\n", chatter.name, inactiveFor)
		if inactiveFor.Minutes() > 5 {
			evictionList = append(evictionList, chatter)
		}
		return true
	})
	for _, chatter := range evictionList {
		fmt.Printf("Evicting %s due to timeout.\n", chatter.name)
		if closeErr := chatter.conn.Close(); closeErr != nil {
			log.Println("Error while closing:", closeErr)
		}
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
	chatter.conn.WriteMessage(websocket.TextMessage, whoReply)
}

func parseMessage(p []byte) *chatMessage {
	msg := &chatMessage{}
	err := json.Unmarshal(p, msg)
	if err != nil {
		fmt.Println("JSON Error:", err)
	}
	return msg
}

func evictAndBroadcast(chatter *Chatter) {
	if closeErr := chatter.conn.Close(); closeErr != nil {
		log.Println("Error while closing:", closeErr)
	}
	chatters.Delete(chatter.uuid)
	log.Printf("Evicted connection: %s. (%s).\n", chatter.name, chatter.uuid)
	evictionMsg := chatMessage{Name: "Chat Bot", Text: fmt.Sprintf("%s left.", chatter.name)}
	evictionMsgJSON, _ := json.Marshal(evictionMsg)
	chatters.Range(func(uid, c interface{}) bool {
		peer := c.(*Chatter)
		if err := peer.conn.WriteMessage(websocket.TextMessage, evictionMsgJSON); err != nil {
			log.Println(err)
			log.Printf("Couldn't send eviction message to: %s. (%s).\n", peer.name, uid)
		}
		return true
	})
}

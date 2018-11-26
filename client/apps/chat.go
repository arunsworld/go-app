package apps

import (
	"github.com/arunsworld/go-js-dom"
	"github.com/arunsworld/gopherjs/websocket"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

var chatContentHTML = `
<style text="text/css">
	.chat {
		border: 2px solid #dedede;
		background-color: #f1f1f1;
		border-radius: 5px;
		padding: 10px;
		margin: 10px 0;
	}
	.darker {
		border-color: #ccc;
		background-color: #ddd;
	}
	.chat::after {
		content: "";
		clear: both;
		display: table;
	}
	.time-right {
		float: right;
		color: #aaa;
	}
	
	.time-left {
		float: left;
		color: #999;
	}
</style>
<div class="row justify-content-sm-center">
	<div class="col-sm-6">
		<div class="card">
			<div class="card-header"><b>My Chat</b></div>
				<div class="card-body">
					<div id="connection" class="text-center">
						<input type="text" class="form-control" id="name" placeholder="Name">
						<br/>
						<button class="btn btn-success" id="connect">Connect</button>
					</div>
					<div id="chatWindow" style="display:none;">
						<input type="text" class="form-control" id="newchat" placeholder="Enter a message...">
						<div style="height: 450px;overflow: auto;" id="messages">
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
`

var ping = js.Global.Get("Audio").New("/static/audio/ping.mp3")

// ChatContentProducer produces the DIV for chat with all the behaviors using WebSocket
var ChatContentProducer = func() *dom.HTMLDivElement {
	f, result := divElementFromContent(chatContentHTML)

	var chatSocket *websocket.WebSocket
	socketURL := websocket.WSURL(D, "/chatws")

	connectionOpen := make(chan bool)
	connectionClosed := make(chan bool)
	messageChanel := make(chan string)

	connection := f.GetElementByID("connection").(*dom.HTMLDivElement)
	name := f.GetElementByID("name").(*dom.HTMLInputElement)
	chatWindow := f.GetElementByID("chatWindow").(*dom.HTMLDivElement)
	connectBtn := f.GetElementByID("connect").(*dom.HTMLButtonElement)
	makeConnection := func(e dom.Event) {
		if name.Value == "" {
			js.Global.Call("alert", "Please enter your Name to Connect.")
			return
		}
		connectBtn.Disabled = true
		var err error
		chatSocket, err = websocket.New(socketURL)
		if err != nil {
			js.Global.Call("alert", "Could not connect to Chat Server.")
			return
		}
		chatSocket.AddEventListener("error", false, func(e *js.Object) {
			connectionClosed <- true
		})
		chatSocket.AddEventListener("close", false, func(e *js.Object) {
			connectionClosed <- true
		})
		chatSocket.AddEventListener("open", false, func(e *js.Object) {
			connectionOpen <- true
		})
		chatSocket.AddEventListener("message", false, func(e *js.Object) {
			messageChanel <- e.Get("data").String()
		})
	}
	connectBtn.AddEventListener("click", false, makeConnection)
	name.AddEventListener("keypress", false, func(e dom.Event) {
		ke := e.(*dom.KeyboardEvent)
		if ke.Key == "Enter" || ke.KeyCode == 13 {
			makeConnection(e)
		}
	})

	messages := f.GetElementByID("messages").(*dom.HTMLDivElement)
	newchat := f.GetElementByID("newchat").(*dom.HTMLInputElement)
	newchat.AddEventListener("keypress", false, func(e dom.Event) {
		ke := e.(*dom.KeyboardEvent)
		if ke.Key == "Enter" || ke.KeyCode == 13 {
			if newchat.Value == "" {
				return
			}
			msgTxt := newchat.Value
			console.Log(chatSocket.ReadyState.String())
			err := chatSocket.Send(JSON.Call("stringify", js.M{
				"name": name.Value,
				"text": msgTxt,
			}))
			if err != nil {
				console.Log(err.Error())
			}
			newchat.Value = ""
		}
	})

	go func() {
		for {
			select {
			case <-connectionOpen:
				connection.Style().Set("display", "none")
				chatWindow.Style().Set("display", "block")
				chatSocket.Send(JSON.Call("stringify", js.M{
					"name": name.Value,
					"text": "Hi, I have just joined the party.",
				}))
			case <-connectionClosed:
				connection.Style().Set("display", "block")
				chatWindow.Style().Set("display", "none")
				connectBtn.Disabled = false
			case msg := <-messageChanel:
				msgO := JSON.Call("parse", msg)
				dark := true
				left := false
				peerName := msgO.Get("name").String()
				playPing := true
				if peerName == name.Value {
					dark = false
					left = true
					playPing = false
				}
				msgE := newMessage(msgO.Get("text").String(), peerName, dark, left)
				messages.InsertBefore(msgE, messages.FirstChild())
				if playPing {
					ping.Call("play")
				}
				// js.Global.Call("jQuery", messages).Call("scrollTop", 1000000)
			}
		}
	}()

	return result
}

func newMessage(txt string, name string, dark bool, left bool) *dom.HTMLDivElement {
	msg := D.CreateElement("div").(*dom.HTMLDivElement)
	msg.Class().Add("chat")
	if dark {
		msg.Class().Add("darker")
	}
	p := D.CreateElement("p")
	p.SetTextContent(txt)
	msg.AppendChild(p)
	cls := "time-right"
	if left {
		cls = "time-left"
	}
	s := D.CreateElement("span")
	s.SetTextContent(name + ". " + timeNow())
	s.Class().Add(cls)
	msg.AppendChild(s)
	return msg
}

func timeNow() string {
	d := js.Global.Get("Date").New()
	h := d.Call("getHours").Int()
	hT := O.New(h).String()
	if h < 10 {
		hT = "0" + hT
	}
	m := d.Call("getMinutes").Int()
	mT := O.New(m).String()
	if m < 10 {
		mT = "0" + mT
	}
	return hT + ":" + mT
}

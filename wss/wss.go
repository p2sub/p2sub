package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Operator alias of boolean
type Operator = bool

// Channel operators
const (
	Write Operator = true
	Read           = false
)

// ChannelIO structure
type ChannelIO struct {
	ID       uint64
	Operator Operator
	Data     []byte
}

// WebsocketServer websocket server struct
type WebsocketServer struct {
	receiver    chan ChannelIO
	sender      chan ChannelIO
	writer      chan ChannelIO
	connections map[uint64]*websocket.Conn
	uniqueID    uint64
	syncMux     sync.Mutex
}

// New instance of websocket server
func New() *WebsocketServer {
	return &WebsocketServer{
		connections: make(map[uint64]*websocket.Conn),
		uniqueID:    0,
		receiver:    make(chan ChannelIO),
		sender:      make(chan ChannelIO),
		writer:      make(chan ChannelIO),
		syncMux:     sync.Mutex{},
	}
}

// GetUniqueID for connect
func (wss *WebsocketServer) GetUniqueID() uint64 {
	wss.syncMux.Lock()
	wss.uniqueID++
	defer wss.syncMux.Unlock()
	return wss.uniqueID
}

// UpgradeConnection to websocket
func (wss *WebsocketServer) UpgradeConnection(res http.ResponseWriter, req *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
	connection, err := upgrader.Upgrade(res, req, nil)
	if err == nil {
		channelID := wss.GetUniqueID()
		wss.connections[channelID] = connection
		// Wipe our ass after we leave
		defer func() {
			log.Println("Clean and close")
			delete(wss.connections, channelID)
			connection.Close()
		}()
		for {
			messageType, message, err := connection.ReadMessage()
			if err == nil {
				wss.receiver <- ChannelIO{ID: channelID, Operator: Read, Data: message}
			} else {
				log.Println("New error:", err)
				break
			}
			select {
			case n := <-wss.writer:
				wss.sender <- n
				connection.WriteMessage(messageType, n.Data)
			}
		}
	} else {
		defer connection.Close()
	}

}

// Receiving data from channel
func (wss *WebsocketServer) Receiving() <-chan ChannelIO {
	return wss.receiver
}

// Sending data from channel
func (wss *WebsocketServer) Sending() <-chan ChannelIO {
	return wss.sender
}

// Send message to channel
func (wss *WebsocketServer) Send(channelID uint64, data []byte) {
	wss.writer <- ChannelIO{ID: channelID, Operator: Write, Data: data}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	websocketServer := New()
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", websocketServer.UpgradeConnection)
	go func(websocketServer *WebsocketServer) {
		for {
			select {
			case n := <-websocketServer.Sending():
				log.Println("Sent", n)
			case n := <-websocketServer.Receiving():
				log.Println("Received:", n)
				websocketServer.Send(n.ID, n.Data)
			}
		}
	}(websocketServer)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
var ws;
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            console.log(evt);
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    (function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            console.log(evt);
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    })();

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))

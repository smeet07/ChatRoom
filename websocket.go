package main
import (
    "bytes"
    "encoding/json"
    "html/template"
    "log"

    //"github.com/gofiber/fiber/v2"
    "github.com/gofiber/websocket/v2"
)
type WebSocketServer struct{
	clients map[*websocket.Conn] bool
	broadcast chan *Message
}
func NewWebSocket() *WebSocketServer{
	return &WebSocketServer{
		clients:make(map[*websocket.Conn]bool),
		broadcast:make(chan *Message),
	}
}
func (s *WebSocketServer) HandleWebSocket(ctx *websocket.Conn){
	s.clients[ctx]=true
	defer func(){
		delete(s.clients,ctx)
		ctx.Close()
	}()
	for {
		_,msg,err:=ctx.ReadMessage()
		if err!=nil{
			log.Println("Read Error:",err)
			break
		}
		var message Message
		if err:=json.Unmarshal(msg,&message);err!=nil{
			log.Println("Unmarshal Error:",err)
		}
		s.broadcast<- &message
	}
}
func (s *WebSocketServer) HandleMessages(){
	for {
		msg:=<-s.broadcast
		for client:=range s.clients{
			err:=client.WriteMessage(websocket.TextMessage,getMessageTemplate(msg))
			if err!=nil{
				log.Println("Write Error:",err)
				client.Close()
				delete(s.clients,client)
			}
		}	
	}
}
func getMessageTemplate(msg *Message) []byte{
	tmpl,err:=template.ParseFiles("views/message.html")
	if err!=nil{
		log.Println("Template Parse Error:",err)
		return nil
	}
	var rendereredMessage bytes.Buffer
	if err:=tmpl.Execute(&rendereredMessage,msg);err!=nil{
		log.Println("Template Execute Error:",err)
		return nil
	}
	return rendereredMessage.Bytes()

}
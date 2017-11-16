package main

import(
	websocket "github.com/gorilla/websocket"
	"time"
	"log"
)

type client struct {
	socket * websocket.Conn
	//send chan []byte
	send chan *message
	room *room
	userData map[string]interface{}
}

func (c *client) read(){
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now()
			msg.Name = c.userData["name"].(string)
			//if avatarURL, ok := c.userData["avatar_url"]; ok {
			//	msg.AvatarURL = avatarURL.(string)
			//}
			msg.AvatarURL, err = c.room.avatar.GetAvatarURL(c); if  err != nil{
				log.Fatal("error", err)
			}
			c.room.forward <- msg
		} else {
			break
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break
		}
	}
	c.socket.Close()
}
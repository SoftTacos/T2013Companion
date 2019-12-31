package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type GameServer struct {
	Players []*PlayerClient
	DM      DMClient
}

func (gs *GameServer) AddPlayerClient(client *PlayerClient) {
	gs.Players = append(gs.Players, client)
}

func (gs *GameServer) Handle() {
	//for{
	//handle all input events from DMs + Players
}

type Client interface {
	Start()
	Listener()
	Writer()
}

type PlayerClient struct {
	conn      *websocket.Conn
	char      *Character
	requests  chan []byte
	responses chan []byte
}

func (pc *PlayerClient) Start() {
	go pc.Listener()
	go pc.Writer()
}

func (pc *PlayerClient) Listener() {
	for {
		_, message, err := pc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("RECEIVED:", string(message))
		//pc.requests <- message
		pc.responses <- message
	}
	fmt.Println("Closing Socket for Player: ", pc.char.Name)
}

func (pc *PlayerClient) Writer() {
	for {
		select {
		case response := <-pc.responses:
			fmt.Println("SENDING: ", string(response))
			if err := pc.conn.WriteMessage(1, response); err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

type DMClient struct {
}

/*
//messageType is an int and can be 1:Text([]uint8|[]byte), 2:binary(), 8:closemessage, 9:ping message, 10:pong message?
func PlayerSocketListener(charName string, conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message

	}
}
*/

/*func PlayerSocketHandler(charName string, conn *websocket.Conn) {
//char := characters[charName]
//	messages := make(chan []byte, 10) //TODO:limit

/*if err := conn.WriteMessage(1, []byte("Hello, Moose!")); err != nil {
	log.Println(err)
	return
}
messageType, p, err := conn.ReadMessage()
if err != nil {
	log.Println(err)
	return
}
fmt.Println(messageType, string(p))*/
//}

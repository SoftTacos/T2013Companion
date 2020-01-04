package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

var RequestRoutingMap = map[uint8]func(*GameRequest, *GameServer){
	0:  statCheck,
	1:  skillCheck,
	2:  attack,
	3:  attack,
	4:  attack,
	5:  attack,
	6:  attack,
	7:  attack,
	8:  attack,
	9:  attack,
	10: getAllChatMessages,
	11: sendChatMessage,
	12: attack,
	13: attack,
	14: attack,
	15: attack,
	16: attack,
	17: attack,
	18: attack,
	19: attack,
	20: updateAllConnectedClients,
}

//MESSAGE FORMAT:charName,statName,bonus,difficulty
func statCheck(gr *GameRequest, gs *GameServer) {
	ary := strings.SplitN(string(gr.Message), ",", 4)
	charClientIndex, err := strconv.ParseUint(ary[0], 10, 64)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	statName := ary[1]
	/*
		bonus, err := strconv.ParseInt(ary[2], 10, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	*/
	difficulty, err := strconv.ParseInt(ary[3], 10, 64)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	} //need to figure out relationships between gameserver, client, characters, and how they will be passed around and stored
	result := StatCheck(statName, gs.Clients[uint(charClientIndex)].character, int(difficulty))
	fmt.Println(result)
	//TODO: RESPOND
}

//MESSAGE FORMAT:charName,statName,skillName,bonus,difficulty
func skillCheck(gr *GameRequest, gs *GameServer) {
	ary := strings.SplitN(string(gr.Message), ",", 5)
	charClientIndex, err := strconv.ParseUint(ary[0], 10, 64)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	statName := ary[1]
	skillName := ary[2]
	/*
		bonus, err := strconv.ParseInt(ary[3], 10, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
	*/
	difficulty, err := strconv.ParseInt(ary[4], 10, 64)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	result := SkillCheck(statName, skillName, gs.Clients[uint(charClientIndex)].character, int(difficulty))
	fmt.Println(statName, skillName)
	fmt.Println(result)
	//TODO: RESPOND
}

func getAllChatMessages(gr *GameRequest, gs *GameServer) {
	if gs.ChatMessages.Len() == 0 {
		return
	}
	messages := []byte{}
	for e := gs.ChatMessages.Front(); e != nil; e = e.Next() {
		message := append(e.Value.([]byte), []byte(",")...)
		messages = append(messages, message...)
	}
	sendResponse(EncodeResponse(10, messages[0:len(messages)-1], gr.MessageType), gr.Client)
}

func sendChatMessage(gr *GameRequest, gs *GameServer) {
	//add chat message to chat messages
	fmt.Println(gr.Message)
	gs.ChatMessages.PushBack(gr.Message)
	//send message to everyone
	for _, client := range gs.Clients {
		sendResponse(EncodeResponse(11, gr.Message, gr.MessageType), client)
	}
}

func updateAllConnectedClients(gr *GameRequest, gs *GameServer) {
	sendResponse(EncodeResponse(20, compileAllConnectedClients(), websocket.BinaryMessage), gr.Client)
}

func attack(gr *GameRequest, gs *GameServer) {
	//if player
	//get weapon info
	//get target
	//get target's info: cover, distance, armor
	//get attack's info: number of shots
	//roll to hit
	//roll hit location

	//if DM

	fmt.Println("ATTACK FUNCTION CALLED")
}

package main

//prototype initiative tracker for a potential companion app for a TTRPG called Twilight2013, it's rules are complicated for a TTRPG and require too much effort to be engaging.
//This idea is this basic functionality could augment the DM and allow for smoother(more engaging) gameplay
//the combat is extremely complex for a TTRPG so a webapp that players connect to in order to manage gear, attempting shots, and managing hits would take the "bean counting" out of the game

import (
	"bufio"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

var randy *rand.Rand
var rules Rules
var globals Globals
var itemData map[string]Item //this holds the item "definitions", same struct as the items players have, TODO: Make player's weapons "instances" so they can eventually have stats like "status"
//var weapons []RangedWeaponItem //placeholder until item system scope is known

type Rules struct {
	EncMap        map[string]uint8 //encumbrance map of name to initiative value
	SkillLevelMap map[uint8]uint8  //skill level to #d20 to roll: 25->2d20
	TurnActions   map[uint]func(*Turn)
}

func (r *Rules) Init() {
	r.EncMap = map[string]uint8{
		"Overloaded": 5,
		"Heavy":      7,
		"Moderate":   9,
		"Light":      12,
		"None":       15,
		"":           9,
	}
	r.SkillLevelMap = make(map[uint8]uint8)
	r.TurnActions = map[uint]func(*Turn){
		1: Turn_Attack,
		2: Turn_Move,
		3: Turn_ChangeStance,
		4: Turn_Communicate,
		5: Turn_Reload,
	}
}

type Globals struct {
	whitespaceRegex *regexp.Regexp
	reader          *bufio.Reader
}

//FUNCTIONS

func NumberMenu(max uint) uint {
	var validOption uint
	for validNumber := false; !validNumber; {
		input, _ := globals.reader.ReadString('\n')
		input = globals.whitespaceRegex.ReplaceAllString(input, "")
		option, err := strconv.ParseUint(input, 10, 64)
		if err != nil {
			fmt.Println(err)
			print("Input could not be read as a number, please provide a valid number\n")
			continue
		}
		if uint(option) > max {
			print("Input was too high\n")
			continue
		}
		validOption = uint(option)
		validNumber = true
	}
	return validOption
}

func advantage(rolls []uint8) uint8 {
	lowest := rolls[0]
	for i, _ := range rolls {
		if lowest > rolls[i] {
			lowest = rolls[i]
		}
	}
	return lowest
}

func nd20(n uint8) []uint8 {
	rolls := make([]uint8, n)
	for i, _ := range rolls {
		rolls[i] = d20()
	}
	return rolls
}

func d20() uint8 {
	return uint8((randy.Uint64() % 20) + 1)
}

func Reorder(inits *list.List) {
	init := inits.Front().Value.(*Turn).Init
	for e := inits.Back(); e != nil; e = e.Prev() {
		turn := e.Value.(*Turn)

		if init < turn.Init {
			fmt.Println("REORDERING")
			inits.InsertAfter(inits.Remove(inits.Front()), e)
			break
		}
	}
}

/*for e := inits.Front(); e != nil; e = e.Next() {
	fmt.Println(e.Value)
}*/

//combat is rounds of EoF, EoF is a series of turns that happen until initiatives all go to 0
//combat just does one round of EoF
func Combat(inits *list.List) {
	turn := inits.Front().Value.(*Turn)
	for { //roundOver := false; !roundOver; {
		//turn is a collection of ticks of the same number
		//TAKE TURN
		TakeTurn(turn)
		Reorder(inits)
		//decrement Turn.Init value, move the item within the linked list
		//look at the front of the LL, is it the same as turn? Then do that
		turn = inits.Front().Value.(*Turn)
		if turn.Init <= 0 {
			break
		}
	}
}

//todo: will want to just generate the list and sort it later, mostly as an exercise in LinkedList sorting
func GenerateInitiatives(chars []*Character) *list.List {
	initiatives := make([]*Turn, len(chars))
	for i, _ := range chars {
		initiatives[i] = &Turn{
			Init: int(chars[i].InitiativeCheck()),
			Char: chars[i],
		}
	}
	sort.Slice(initiatives, func(i, j int) bool { return initiatives[i].Init > initiatives[j].Init })
	initList := list.New()
	for i, _ := range initiatives {
		initList.PushBack(initiatives[i])
	}
	return initList
}

func LoadTextFile(filename string) []byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	filestring, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	return filestring
}

func Setup() {
	itemData = make(map[string]Item)
	globals.whitespaceRegex = regexp.MustCompile(`\s`)
	globals.reader = bufio.NewReader(os.Stdin)
	rules.Init()
	randy = rand.New(rand.NewSource(time.Now().Unix()))
	ReadItemData(LoadTextFile("items.yaml"))
	//ReadWeapons(LoadTextFile("weapons.yaml"))

}

func main() {
	Setup()
	//reading characters from file
	characters := ReadCharacterData(LoadTextFile("testCS.yaml"))

	initiatives := GenerateInitiatives(characters)
	Combat(initiatives)
}

/*
func randomChar(max int, min int) *Character {
	randName := strconv.Itoa(randy.Int() % 10000)
	stats := make([]uint8, 10)
	for i := 0; i < 10; i++ {
		stats[i] = uint8(randy.Int()%(max-min) + min)
	}
	weapon := &Weapon{
		Name:   "ASDF",
		Speed:  3,
		Damage: 5,
		Bulk:   3,
	}
	randomChar := createChar(randName, stats, weapon)
	return randomChar
}
*/

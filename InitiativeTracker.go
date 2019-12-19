package main

//prototype initiative tracker for a potential companion app for a TTRPG called Twilight2013, it's rules are complicated for a TTRPG and require too much effort to be engaging.
//This idea is this basic functionality could augment the DM and allow for smoother(more engaging) gameplay
//the combat is extremely complex for a TTRPG so a webapp that players connect to in order to manage gear, attempting shots, and managing hits would take the "bean counting" out of the game

import (
	"bufio"
	"bytes"
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

	"github.com/go-yaml/yaml"
)

var randy *rand.Rand
var rules Rules
var globals Globals
var itemData map[string]Item //this holds the item "definitions", same struct as the items players have, except player's items are "instances"
var weapons []WeaponItem     //placeholder until item system scope is known

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

type Character struct {
	Name          string
	Stats         map[string]uint8 //name->level
	Items         map[string]Item
	CurrentWeapon *WeaponItem //placeholder, will be removed and replaced with an item specific call
}

func (c *Character) Init() {
	c.Name = "BLANK"
	c.Stats = map[string]uint8{
		"AWA":  0,
		"CDN":  0,
		"FIT":  0,
		"MUS":  0,
		"COG":  0,
		"EDU":  0,
		"PER":  0,
		"RES":  0,
		"CUF":  0, //DONT USE
		"OODA": 0, //DONT USE
	}
	//c.Gear = make(map[string]float32)
}

func (c *Character) InitiativeCheck() uint8 {
	encIni := rules.EncMap[c.Encumbrance()]
	roll := advantage(nd20(2))

	//initiative is 2d20 VS OODA, OODA is TN
	checkMargin := (int(c.Stats["AWA"]) - int(roll))
	if checkMargin < 0 {
		return encIni
	}
	return encIni + uint8(checkMargin)*2
}

func (c *Character) Encumbrance() string {
	return ""
}

type Item interface {
	GetType() string
	GetWeight() float32
}

type MiscItem struct {
	Name   string
	weight float32
}

func (m *MiscItem) GetType() string    { return "misc" }
func (m *MiscItem) GetWeight() float32 { return m.weight }

type WeaponItem struct {
	Name   string
	Speed  int
	Damage int
	Bulk   int
	weight float32
}

func (w *WeaponItem) GetType() string    { return "weapon" }
func (w *WeaponItem) GetWeight() float32 { return w.weight }

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

type Turn struct {
	Init int
	Char *Character
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

func TakeTurn(turn *Turn) {
	//fmt.Printf("\n\t%s's turn:%d\n", turn.Char.Name, turn.Init)
	validRange := uint(5) //lazy!
	fmt.Printf("\n\t%s's turn:%d\n1:Attack\n2:Move\n3:Change Stance\n4:Communicate\n5:Reload\n", turn.Char.Name, turn.Init)
	choice := NumberMenu(validRange)
	rules.TurnActions[choice](turn)
}

func Turn_Attack(turn *Turn) {
	fmt.Println("ATTACKING")
	turn.Init -= turn.Char.CurrentWeapon.Speed
}

func Turn_ChangeStance(turn *Turn) {
	fmt.Println("STANCE")
	turn.Init -= 2 //stance changes are static cost of 2
}

func Turn_Communicate(turn *Turn) {
	fmt.Println("COMMUNICATE")
	turn.Init -= int(NumberMenu(20))
}

func Turn_Reload(turn *Turn) {
	fmt.Println("RELOAD")
	turn.Init -= turn.Char.CurrentWeapon.Bulk
}

func Turn_Move(turn *Turn) {
	fmt.Println("MOVE")

}

func check() uint8 {
	return 0
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

//reorder not working correctly atm
func Reorder(inits *list.List) {
	//fmt.Println("REORDER", inits.Front().Value.(*Turn).Init, inits.Front().Next().Value.(*Turn).Init)
	init := inits.Front().Value.(*Turn).Init
	//find first element that is below the Front's init
	//insert before that element
	for e := inits.Front(); e != nil; e = e.Next() {
		turn := e.Value.(*Turn)
		if turn.Init < init {
			inits.InsertBefore(inits.Remove(inits.Front()), e)
			break
		}

	}
}

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

func ReadCharacterData(data []byte) []*Character {
	delimeter := []byte("---")
	dataSlices := bytes.Split(data, delimeter)
	characters := make([]*Character, len(dataSlices))
	itemsTag := "#!items"
	//characterTag := "#!character"

	for i, charData := range dataSlices { //TODO: Break down into abstracted functions that read each section of character data to allow for flexible parsing of out of order data
		char := Character{}
		subSlices := bytes.Split(charData, []byte(itemsTag))
		//CHARACTER
		err := yaml.Unmarshal(subSlices[0], &char)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		//TODO: Validate character data
		characters[i] = &char
		//ITEMS
		itemsSubset := struct {
			Items         []string
			CurrentWeapon string
		}{}

		err = yaml.Unmarshal(subSlices[1], &itemsSubset)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		char.Items = make(map[string]Item)
		for _, curItemName := range itemsSubset.Items {
			//TODO: validate if the item exists in item list
			char.Items[curItemName] = itemData[curItemName]
		}

		fmt.Println(char)
	}
	if characters[len(characters)-1].Name == "" {
		characters = characters[0 : len(characters)-1]
	}
	return characters
}

func ReadWeapons(data []byte) {
	delimeter := []byte("---")
	dataSlices := bytes.Split(data, delimeter)
	weapons = make([]WeaponItem, len(dataSlices))

	for i, wepData := range dataSlices {
		weapon := WeaponItem{}
		err := yaml.Unmarshal(wepData, &weapon)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		weapons[i] = weapon
	}
	if weapons[len(weapons)-1].Name == "" {
		weapons = weapons[0 : len(weapons)-1]
	}
}

func ReadItemData(data []byte) {
	//itemMap := make(map[string][]map[string]string)

}

func Setup() {
	globals.whitespaceRegex = regexp.MustCompile(`\s`)
	globals.reader = bufio.NewReader(os.Stdin)
	rules.Init()
	randy = rand.New(rand.NewSource(time.Now().Unix()))
	ReadItemData(LoadTextFile("items.yaml"))
	ReadWeapons(LoadTextFile("weapons.yaml"))

}

func main() {
	Setup()
	//reading characters from file
	characters := ReadCharacterData(LoadTextFile("testCS.yaml"))

	initiatives := GenerateInitiatives(characters)
	Combat(initiatives)
}

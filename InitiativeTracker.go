package main

//prototype initiative tracker for a potential companion app for a TTRPG called Twilight2013, it's rules are complicated for a TTRPG and require too much effort to be engaging.
//This idea is this basic functionality could augment the DM and allow for smoother(more engaging) gameplay
//the combat is extremely complex for a TTRPG so a webapp that players connect to in order to manage gear, attempting shots, and managing hits would take the "bean counting" out of the game

import (
	"bufio"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"time"
)

var randy *rand.Rand
var rules Rules
var globals Globals
var itemData map[string]Item //this holds the item "definitions", same struct as the items players have, TODO: Make player's weapons "instances" so they can eventually have stats like "status"
var characters map[string]*Character
var gameServers []GameServer

var newline = []byte{'\n'}
var space = []byte{' '}

type Rules struct {
	EncMap        map[string]uint8 //encumbrance map of name to initiative value
	SkillLevelMap map[uint8]uint8  //skill level to #d20 to roll: 25->2d20
	StatNames     []string         //not sure where to put this so it goes in rules for now
	SkillNames    []string
	TurnActions   map[uint]func(*Turn) //
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
	//r.StatNames = [10]string{"AWA", "CDN", "FIT", "MUS", "COG", "EDU", "PER", "RES", "CUF", "OODA"}
	//r.Skills = [28]string{}
}

type Globals struct {
	whitespaceRegex *regexp.Regexp
	reader          *bufio.Reader
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

func SetupGame() {
	itemData = make(map[string]Item)
	globals.whitespaceRegex = regexp.MustCompile(`\s`)
	globals.reader = bufio.NewReader(os.Stdin)
	rules.Init()
	randy = rand.New(rand.NewSource(time.Now().Unix()))
	ReadItemData(LoadTextFile("data\\items.yaml"))
	ReadRuleData(LoadTextFile("data\\rules.yaml"))
}

//refresh the pages so I don't have to restart the server every time I update a page
func refresh() {
	for {
		time.Sleep(1 * time.Second)
		pages["CharSelectPage"] = LoadTextFile("pages\\CharSelect.html")
		pages["CharacterPage"] = LoadTextFile("pages\\Character.html")
		pages["ItemCard"] = LoadTextFile("pages\\ItemCard.html")
		pages["SkillChartElement"] = LoadTextFile("pages\\SkillChartElement.html")
		pages["StatusChart"] = LoadTextFile("pages\\StatusChart.html")

	}
}

func debug() {
	reader := bufio.NewReader(os.Stdin)
	whitespaceRegex := regexp.MustCompile(`\s`)
	for {
		input, _ := reader.ReadString('\n')
		input = whitespaceRegex.ReplaceAllString(input, "")

	}
}

func main() {
	SetupGame()
	SetupServer()
	ReadCharacterData(LoadTextFile("data\\characters.yaml"))
	go refresh() //im lazy
	go http.ListenAndServe(":8082", nil)
	//fmt.Println("COMBAT")
	//initiatives := GenerateInitiatives(characters)
	//Combat(initiatives)
	debug()
}

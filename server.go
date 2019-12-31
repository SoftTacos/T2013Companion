package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var pages map[string][]byte
var router *mux.Router

func SetupServer() {
	pages = make(map[string][]byte)
	pages["CharacterSelectPage"] = LoadTextFile("pages\\CharSelect.html")
	pages["CharacterPage"] = LoadTextFile("pages\\Character.html")
	pages["DMPage"] = LoadTextFile("pages\\DM.html")
	pages["ItemCard"] = LoadTextFile("pages\\ItemCard.html")
	pages["SkillChartElement"] = LoadTextFile("pages\\SkillChartElement.html")
	pages["StatusChart"] = LoadTextFile("pages\\StatusChart.html")
	pages["CharacterCard"] = LoadTextFile("pages\\CharacterCard.html")

	router = mux.NewRouter()
	SetRoutes()
	fmt.Println("STARTED; IP:", GetOutboundIP())

	gameServers = []GameServer{
		GameServer{
			Players:  []*Client{},
			Requests: make(chan *GameRequest, 1),
		},
	}
	go gameServers[0].Handle() //TODO move in refactor
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

func SetRoutes() {
	router.HandleFunc("/", BlankPage)
	router.HandleFunc("/select", CharSelectPage)
	router.HandleFunc("/dm", DMPage)
	router.HandleFunc("/dm/socket", DMSocket)
	router.HandleFunc("/player/{[A-Za-z]+}", CharacterPage)
	router.HandleFunc("/player/{[A-Za-z]+}/edit", CharacterEditPage)
	router.HandleFunc("/player/{[A-Za-z]+}/socket", PlayerSocket)

	http.Handle("/", router)
}

func BlankPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BLANK")
}

//TODO:this is super gross and janky, will be refactoring all the page preprocessing once scope is known
//TODO: Consider pre-processing these once and storing them rather than each time page is requested
func CharSelectPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Character Select")
	charNames := []string{}
	for _, char := range playerCharacters {
		charNames = append(charNames, char.Name)
	}

	t := template.Must(template.New("").Parse(`{{range .}}<li><a href="player/{{.}}">{{.}}</a></li>{{end}}`)) //<li><a href="#">HTML</a></li>
	var str strings.Builder
	if err := t.Execute(&str, charNames); err != nil {
		log.Fatal(err)
	}
	page := pages["CharacterSelectPage"]
	page = bytes.Replace(page, []byte("##OPTIONS##"), []byte(str.String()), 1)
	fmt.Fprintf(w, string(page))
}

//this is super gross and janky, will be refactoring all the page preprocessing once scope is known
func CharacterPage(w http.ResponseWriter, r *http.Request) {
	charName := strings.Split(html.EscapeString(r.URL.Path), "/")[2]
	if _, ok := playerCharacters[charName]; !ok {
		fmt.Println("Error: Character not found:", charName)
		fmt.Fprintf(w, "404 Character Not Found!")
		return
	}

	fmt.Println("Character Page: ", charName)
	page := pages["CharacterPage"]
	charData := playerCharacters[charName]

	page = bytes.Replace(page, []byte("##NAME##"), []byte(charData.Name), 1)
	page = bytes.Replace(page, []byte("##STATS##"), generateStatChart(charData), 1)
	page = bytes.Replace(page, []byte("##SKILLS##"), generateSkillChart(charData), 1)
	page = bytes.Replace(page, []byte("##CURRENT_WEAPON##"), generateCurrentWeaponCard(charData.CurrentWeapon), 1)
	page = bytes.Replace(page, []byte("##ITEMS##"), generateHtmlItemList(charData), 1)
	page = bytes.Replace(page, []byte("##STATUSCHART##"), pages["StatusChart"], 1)

	fmt.Fprintf(w, string(page))
}

func generateStatChart(char *Character) []byte {
	chart := make([]byte, 450)
	for _, statName := range rules.StatNames { //key, val := range charData.Stats
		stat := strconv.FormatUint(uint64(char.Stats[statName]), 8)
		chart = append(chart, []byte(`<tr class="text-left"><td>`+statName+`</td><td>`+stat+`</td></tr>`)...)
	}
	return chart
}

func generateSkillChart(char *Character) []byte {
	chart := make([]byte, len(pages["SkillChartElement"])*30)
	for _, skillName := range rules.SkillNames {
		element := pages["SkillChartElement"]
		element = bytes.Replace(element, []byte("##SKILLNAME##"), []byte(skillName), 1)
		skill := strconv.FormatUint(uint64(char.Skills[skillName]), 8)
		element = bytes.Replace(element, []byte("##SKILL##"), []byte(skill), 1)
		chart = append(chart, element...)
	}
	return chart
}

//TODO USE THE CARD CLASS REEE
func generateCurrentWeaponCard(item Item) []byte {
	newPage := make([]byte, len(pages["ItemCard"]))
	copy(newPage, pages["ItemCard"])
	newPage = bytes.Replace(newPage, []byte("##CLASS##"), []byte("col border"), 1)
	newPage = bytes.Replace(newPage, []byte("##NAME##"), []byte(item.GetName()), 1)
	newPage = bytes.Replace(newPage, []byte("##TYPE##"), []byte(item.GetType()), 1)
	newPage = bytes.Replace(newPage, []byte("##DESCRIPTION##"), []byte(item.GetDescription()), 1)
	return newPage
}

func generateHtmlItemList(char *Character) []byte {
	list := make([]byte, 1000) //TODO
	for item := range char.Items {
		list = append(list, generateItemCard(char.Items[item])...)
	}
	return list
}

func generateItemCard(item Item) []byte {
	newPage := make([]byte, len(pages["ItemCard"]))
	copy(newPage, pages["ItemCard"])
	newPage = bytes.Replace(newPage, []byte("##CLASS##"), []byte("col-6 col-sm-6 col-md-5 col-lg-4 border"), 1)
	newPage = bytes.Replace(newPage, []byte("##NAME##"), []byte(item.GetName()), 1)
	newPage = bytes.Replace(newPage, []byte("##TYPE##"), []byte(item.GetType()), 1)
	newPage = bytes.Replace(newPage, []byte("##DESCRIPTION##"), []byte(item.GetDescription()), 1)
	return newPage
}

func CharacterEditPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CHARACTER EDIT")
}

var upgrader = websocket.Upgrader{ //TODO: Put this somewhere better
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func PlayerSocket(w http.ResponseWriter, r *http.Request) {
	charName := strings.Split(r.URL.Path, "/")[2]
	if _, ok := playerCharacters[charName]; !ok {
		fmt.Println(charName, "Does not exist!")
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("Socket Connected! ", charName)

	gameServers[0].AddClient(playerCharacters[charName], ws) //TODO: Scale this shit up yo
	//client.Start()
}

func DMPage(w http.ResponseWriter, r *http.Request) {
	page := pages["DMPage"]

	page = bytes.Replace(page, []byte("##CHARACTERS##"), generateCharacterCardList(), 1)

	fmt.Fprintf(w, string(page))
}

func generateCharacterCardList() []byte {
	list := make([]byte, 1000) //TODO
	for char := range playerCharacters {
		list = append(list, generateCharacterCard(playerCharacters[char])...)
	}
	return list
}

func generateCharacterCard(char *Character) []byte {
	newPage := make([]byte, len(pages["CharacterCard"]))
	copy(newPage, pages["CharacterCard"])

	newPage = bytes.Replace(newPage, []byte("##NAME##"), []byte(char.Name), 1)

	return newPage
}

func DMSocket(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("DM Socket Connected!")

	gameServers[0].AddClient(nil, ws) //TODO: Scale this shit up yo
	//client.Start()
}

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
)

var pages map[string][]byte
var router *mux.Router

func SetupServer() {
	pages = make(map[string][]byte)
	pages["CharacterSelectPage"] = LoadTextFile("pages\\CharSelect.html")
	pages["CharacterPage"] = LoadTextFile("pages\\Character.html")
	pages["ItemCard"] = LoadTextFile("pages\\ItemCard.html")
	pages["SkillChartElement"] = LoadTextFile("pages\\SkillChartElement.html")

	router = mux.NewRouter()
	SetRoutes()
	fmt.Println("STARTED; IP:", GetOutboundIP())
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
	router.HandleFunc("/dmSocket", DMSocket)
	router.HandleFunc("/player/{[A-Za-z]+}", CharacterPage)
	router.HandleFunc("/player/{[A-Za-z]+}/edit", CharacterEditPage)
	router.HandleFunc("/playerSocket", PlayerSocket)

	http.Handle("/", router)
}

//PAGE FUNCTIONS

func BlankPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BLANK")
}

//TODO:this is super gross and janky, will be refactoring all the page preprocessing once scope is known
//TODO: Consider pre-processing these once and storing them rather than each time page is requested
func CharSelectPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Character Select")
	charNames := []string{}
	for _, char := range characters {
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

func DMPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("DM PAge")
	fmt.Fprintf(w, "DM PAGE")
}

func DMSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DM SOCKET")
}

//this is super gross and janky, will be refactoring all the page preprocessing once scope is known
func CharacterPage(w http.ResponseWriter, r *http.Request) {
	charName := strings.Split(html.EscapeString(r.URL.Path), "/")[2]
	if _, ok := characters[charName]; !ok {
		fmt.Println("Error: Character not found:", charName)
		fmt.Fprintf(w, "404 Character Not Found!")
		return
	}

	fmt.Println("Character Page: ", charName)
	page := pages["CharacterPage"]
	charData := characters[charName]
	page = bytes.Replace(page, []byte("##NAME##"), []byte(charData.Name), 1)
	page = bytes.Replace(page, []byte("##STATS##"), generateStatChart(charData), 1)
	page = bytes.Replace(page, []byte("##SKILLS##"), generateSkillChart(charData), 1)
	page = bytes.Replace(page, []byte("##CURRENT_WEAPON##"), generateCurrentWeaponCard(charData.CurrentWeapon), 1)
	page = bytes.Replace(page, []byte("##ITEMS##"), generateHtmlItemList(charData), 1)

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
	/*chart := make([]byte, 450)
	for _, skillName := range rules.SkillNames {
		skill := strconv.FormatUint(uint64(char.Skills[skillName]), 8)
		chart = append(chart, []byte(`<div class="col-4 border pt-0">`+skillName+` `+skill+`</div>`)...) //TODO: make them format like a table
	}
	return chart
	*/
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
	return newPage //strBuilder.String()
}

func generateHtmlItemList(char *Character) []byte {
	list := make([]byte, 1000) //TODO
	for item := range char.Items {
		//fmt.Fprintf(&strBuilder, generateItemCard(char.Items[item]))
		list = append(list, generateItemCard(char.Items[item])...)
	}
	return list //strBuilder.String()
}

func generateItemCard(item Item) []byte {
	newPage := make([]byte, len(pages["ItemCard"]))
	copy(newPage, pages["ItemCard"])
	newPage = bytes.Replace(newPage, []byte("##CLASS##"), []byte("col-6 col-sm-6 col-md-5 col-lg-4 border"), 1)
	newPage = bytes.Replace(newPage, []byte("##NAME##"), []byte(item.GetName()), 1)
	newPage = bytes.Replace(newPage, []byte("##TYPE##"), []byte(item.GetType()), 1)
	newPage = bytes.Replace(newPage, []byte("##DESCRIPTION##"), []byte(item.GetDescription()), 1)
	return newPage //strBuilder.String()
}

/*
<div class="col-sm-6 col-md-4 col-lg-3 col-12">
	<a href="/rules/handbook/classes/Berserker" class="v-card v-card--hover v-card--link v-sheet theme--light" tabindex="0" style="height: 100%;">
		<div primary-title="" class="v-card__text">
			<h3>Berserker</h3>
				<div class="text-left">
					<p>Melee combatant who utilizes rage to increase prowess</p>
					<p class="ma-0">
					<strong>Hit Die:</strong> 1d12</p>
					<p class="ma-0"><strong>Primary Ability:</strong> Strength</p>
					<p class="ma-0"><strong>Saves:</strong> Strength and Constitution</p>
				</div>
			</div>
		</a>
	</div>
*/
func CharacterEditPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "CHARACTER EDIT")
}

func PlayerSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PLAYER SOCKET")
}

func PageProcessor(page []byte, tag []byte, format string) {

}

/*
var router RequestRouter

type RequestRouter struct {
	routes []*Route
}

func (r *Router) Route(uri string//w and r//) {
	for i,_:=range r.routes{
		if {//if regex match
			r.routes[i]//call that function with w and r
		}
	}
}

type Route struct {
	Pattern regexp.Regexp
	Func    func(http.ResponseWriter, *http.Request)
}

*/

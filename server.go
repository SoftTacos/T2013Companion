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
	page := pages["CharacterPage"]
	charName := strings.Split(html.EscapeString(r.URL.Path), "/")[2]
	fmt.Println("Character Page: ", charName)
	//TODO: Validate if character even exists
	charData := characters[charName]
	page = bytes.Replace(page, []byte("##NAME##"), []byte(charData.Name), 1)
	var strBuilder strings.Builder

	for i, _ := range rules.StatNames { //key, val := range charData.Stats
		stat := strconv.FormatUint(uint64(charData.Stats[rules.StatNames[i]]), 8)
		fmt.Fprintf(&strBuilder, `<tr class="text-center"><td class="container-fluid">`+rules.StatNames[i]+`</td><td class="container-fluid">`+stat+`</td></tr>`)
	}
	page = bytes.Replace(page, []byte("##STATS##"), []byte(strBuilder.String()), 1)
	page = bytes.Replace(page, []byte("##CURRENT_WEAPON##"), []byte(generateItemCard(charData.CurrentWeapon)), 1)

	fmt.Fprintf(w, string(page))
}

func generateItemCard(item Item) string {

	return ""
}

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

package main

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

var pages map[string][]byte
var router *mux.Router

func SetupServer() {
	pages = make(map[string][]byte)
	pages["CharSelectPage"] = LoadTextFile("pages\\CharSelect.html")
	pages["CharacterPage"] = LoadTextFile("pages\\Character.html")
	router = mux.NewRouter()
	SetRoutes()
}

func SetRoutes() {
	/*
		http.HandleFunc("/", BlankPage)
		http.HandleFunc("/select", CharSelectPage)
		http.HandleFunc("/select/options", CharactersResponse)
		http.HandleFunc("/dm", DMPage)
		http.HandleFunc("/dmSocket", DMSocket)
		http.HandleFunc("/player", PlayerPage)
		http.HandleFunc("/playerSocket", PlayerSocket)
	*/
	router.HandleFunc("/", BlankPage)
	router.HandleFunc("/select", CharSelectPage)
	router.HandleFunc("/select/options", CharactersResponse)
	router.HandleFunc("/dm", DMPage)
	router.HandleFunc("/dmSocket", DMSocket)
	router.HandleFunc("/player/{[A-Za-z]+}", CharacterPage)
	router.HandleFunc("/playerSocket", PlayerSocket)

	http.Handle("/", router)
}

//PAGE FUNCTIONS

func BlankPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "BLANK")
}

//fmt.Fprintf(w, string(visualizerPageHTML))
func CharSelectPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, string(pages["CharSelectPage"]))
}

func CharactersResponse(w http.ResponseWriter, r *http.Request) {
	chars := []byte("DM")
	for _, char := range characters {
		chars = append(chars, []byte(" "+char.Name)...)
	}
	fmt.Fprintf(w, string(chars))
}

func DMPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DM PAGE")
}

func DMSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "DM SOCKET")
}

func CharacterPage(w http.ResponseWriter, r *http.Request) {
	page := pages["CharacterPage"]
	charName := strings.Split(html.EscapeString(r.URL.Path), "/")[2]
	//TODO: Validate if character even exists
	charData := characters[charName]
	page = bytes.Replace(page, []byte("##NAME##"), []byte(charData.Name), 1)
	stats := [2][]string{}
	for key, val := range charData.Stats {
		stats[0] = append(stats[0], string(key))
		stats[1] = append(stats[1], strconv.FormatUint(uint64(val), 8))
	}
	t := template.Must(template.New("").Parse(`<tr>{{range .}}<th>{{.}}</th>{{end}}</tr>`))
	var str strings.Builder
	if err := t.Execute(&str, stats[0]); err != nil {
		log.Fatal(err)
	}
	if err := t.Execute(&str, stats[1]); err != nil {
		log.Fatal(err)
	}

	page = bytes.Replace(page, []byte("##STATS##"), []byte(str.String()), 1)
	fmt.Fprintf(w, string(page))
}

func PlayerSocket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PLAYER SOCKET")
}

func PageProcessor(page []byte) {

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

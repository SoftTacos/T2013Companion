package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

var pages map[string][]byte
var router *mux.Router

func SetupServer() {
	pages = make(map[string][]byte)
	pages["CharSelectPage"] = LoadTextFile("pages\\CharSelect.html")
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
	router.HandleFunc("/player/{[A-Za-z]+}", PlayerPage)
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

func PlayerPage(w http.ResponseWriter, r *http.Request) {
	charName := strings.Split(html.EscapeString(r.URL.Path), "/")[2]
	//TODO: Validate if character even exists
	charData := characters[charName]
	charJSON, err := json.Marshal(charData)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Fprintf(w, string(charJSON))
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

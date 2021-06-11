package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"goji.io"
	"goji.io/pat"
)

func main() {
	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/"), homepage)

	http.ListenAndServe("localhost:8000", mux)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	playerClass := "WIZARD"
	hardcore := false
	var err error

	if r.URL.Query().Get("class") != "" {
		playerClass = r.URL.Query().Get("class")
	}

	if r.URL.Query().Get("hardcore") != "" {
		hardcore, err = strconv.ParseBool(r.URL.Query().Get("hardcore"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}

	fmt.Println(playerClass, hardcore)
	var context LadderContext
	context.Players = getPlayers(strings.ToUpper(playerClass), hardcore)

	func_map := template.FuncMap{"FormatDuration": FormatDuration}
	tmpl, err := template.New("homepage.html").Funcs(func_map).ParseGlob("templates/*.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	err = tmpl.ExecuteTemplate(w, "playerList.html", context)
	if err != nil {
		http.Error(w, err.Error(), 500)

	}
}

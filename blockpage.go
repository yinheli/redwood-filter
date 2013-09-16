package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Functions for displaying block pages.

var blockPage = flag.String("blockpage", "/etc/redwood/block.html", "path to template for block page")

var blockTemplate *template.Template

// transparent1x1 is a single-pixel transparent GIF file.
const transparent1x1 = "GIF89a\x10\x00\x10\x00\x80\xff\x00\xc0\xc0\xc0\x00\x00\x00!\xf9\x04\x01\x00\x00\x00\x00,\x00\x00\x00\x00\x10\x00\x10\x00\x00\x02\x0e\x84\x8f\xa9\xcb\xed\x0f\xa3\x9c\xb4\u068b\xb3>\x05\x00;"

func loadBlockPage() {
	blockTemplate = template.New("blockpage")
	content, err := ioutil.ReadFile(*blockPage)
	if err != nil {
		log.Println("error loading block page template:", err)
		return
	}
	blockTemplate.Funcs(template.FuncMap{
		"eq": func(a, b interface{}) bool {
			return a == b
		},
	})
	_, err = blockTemplate.Parse(string(content))
	if err != nil {
		log.Println("error parsing block page template:", err)
	}
}

type blockData struct {
	URL        string
	Categories string
	User       string
	Group      string
	Tally      string
	Scores     string
}

// showBlockPage sends a page showing that the request was blocked.
func showBlockPage(w http.ResponseWriter, r *http.Request, sc *scorecard, user string) {
	if categories[sc.blocked[0]].invisible {
		// Serve an invisible image instead of the usual block page.
		w.Header().Set("Content-Type", "image/gif")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, transparent1x1)
		return
	}

	blockDesc := make([]string, len(sc.blocked))
	for i, name := range sc.blocked {
		blockDesc[i] = categories[name].description
	}
	data := blockData{
		URL:        r.URL.String(),
		Categories: strings.Join(blockDesc, ", "),
		User:       user,
		Group:      WhichGroup(user),
		Tally:      listTally(stringTally(sc.tally)),
		Scores:     listTally(sc.scores),
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusForbidden)

	err := blockTemplate.Execute(w, data)
	if err != nil {
		log.Println("Error filling in block page template:", err)
	}
}

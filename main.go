package main

import (
	// "fmt"
	"errors"
	"flag"
	"github.com/fredhsu/go-eapi"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
	//"net"
	// "log"
	// "github.com/spf13/viper"
	"encoding/json"
	"fmt"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

type EosNode struct {
	Hostname      string
	MgmtIp        string
	Username      string
	Password      string
	Ssl           bool
	Reachable     bool
	ConfigCorrect bool
	Uptime        float64
	Version       string
}

var templates = template.Must(template.ParseFiles("templates/switches.html"))
var validPath = regexp.MustCompile("^/(edit|save|view|switches)/([a-zA-Z0-9]+)$")

//func fetchVersion(switches []EosNode) {
func fetchConfig(switches EosNode) eapi.JsonRpcResponse {
	c := make(chan eapi.JsonRpcResponse)
	cmds := []string{"enable", "show running-config"}
	prefix := "http"
	if switches.Ssl == true {
		prefix = prefix + "s"
	}
	url := prefix + "://" + switches.Username + ":" + switches.Password + "@" + switches.Hostname + "/command-api"

	go capiFetch(url, cmds, "text", c)
	// get responses need to do for each switch
	msg := <-c
	return msg
}

func capiFetch(url string, cmds []string, format string, c chan eapi.JsonRpcResponse) {
	response := eapi.Call(url, cmds, format)
	c <- response
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("Invalid Page Title")
	}
	return m[2], nil // Title is the second subexp
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, _ := loadPage(title) // title = filename (i.e. view.txt), body = text
	renderTemplate(w, "view", p)
}

func switchesViewHandler(w http.ResponseWriter, r *http.Request) {

	file, _ := os.Open("switches.json")
	decoder := json.NewDecoder(file)
	switches := []EosNode{}
	err := decoder.Decode(&switches)
	if err != nil {
		fmt.Println("error:", err)
	}
	response := eapi.Call("https://admin:admin@bleaf1/command-api/", []string{"show version"}, "json")

	version := response.Result[0]["version"]
	switches[0].Version = version.(string)
	// switches := []EosNode{EosNode{
	// 	Hostname:      "bleaf4",
	// 	Username:      "admin",
	// 	Password:      "admin",
	// 	Ssl:           true,
	// 	MgmtIp:        "1.1.1.1",
	// 	Reachable:     true,
	// 	ConfigCorrect: true,
	// 	Uptime:        0}}
	err = templates.ExecuteTemplate(w, "switches.html", switches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	flag.Parse() // command-line flag parsing
	http.HandleFunc("/switches/", switchesViewHandler)
	// http.HandleFunc("/edit/", makeHandler(editHandler))
	// http.HandleFunc("/save/", makeHandler(saveHandler))
	/*
	   if *addr {
	       l, err := net.Listen("tcp", "127.0.0.1:0")
	       if err != nil {
	           log.Fatal(err)
	       }
	       err = ioutil.WriteFile("final-port.txt", []byte(l.Addr().String()), 0644)
	       if err != nil {
	           log.Fatal(err)
	       }
	       s := &http.Server{}
	       s.Serve(l)
	       return
	   }
	*/
	http.ListenAndServe(":8081", nil)

}

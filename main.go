package main

import (
	// "fmt"
	"github.com/fredhsu/go-eapi"
    "flag"
    "html/template"
    "io/ioutil"
    "net/http"
    "regexp"
    "errors"
    //"net"
    // "log"
)

type Page struct {
    Title string
    Body []byte
}

type EosNode struct {
    Host string
    Username string
    Password string
    Ssl bool
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")


//func fetchVersion(switches []EosNode) {
func fetchVersion(switches EosNode) eapi.JsonRpcResponse {
	c := make(chan eapi.JsonRpcResponse)
	cmds := []string{"enable", "show running-config"}
    prefix := "http"
    if switches.Ssl == true {
        prefix = prefix + "s"
    }
    url := "://" + switches.Username + ":" + switches.Password + "@" + switches.Host + "/command-api"

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

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
    p, err := loadPage(title)
    if err != nil {
        http.Redirect(w, r, "/edit/" + title, http.StatusFound)
        return
    }
    renderTemplate(w, "view", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
    body := r.FormValue("body")
    p := &Page{Title: title, Body: []byte(body)}
    err := p.save()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    http.Redirect(w, r, "/view/" + title, http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl + ".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
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
    http.HandleFunc("/view/", makeHandler(viewHandler))
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

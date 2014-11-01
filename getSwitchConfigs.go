package main

import (
	"encoding/json"
	"fmt"
	"github.com/fredhsu/go-eapi"
	"io/ioutil"
	"os"
)

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

type ChanResponse struct {
	response eapi.JsonRpcResponse
	node     EosNode
}

func versionFetcher(url string, cmds []string, format string, n EosNode,
	c chan ChanResponse) {
	response := eapi.Call(url, cmds, format)
	writeConfig("/Users/fredlhsu/baseconfigs/", n, response.Result[1]["output"].(string))
	c <- ChanResponse{response, n}
}

func writeConfig(path string, n EosNode, config string) {
	filename := path + n.Hostname + ".eos"
	err := ioutil.WriteFile(filename, []byte(config), 0644)
	if err != nil {
		panic(err)
	}
	fmt.Println("wrote to ", filename)
}

func main() {
	file, _ := os.Open("switches.json")
	decoder := json.NewDecoder(file)
	switches := []EosNode{}
	err := decoder.Decode(&switches)

	if err != nil {
		fmt.Println("error:", err)
	}
	cmds2 := []string{"enable", "show running-config"}
	c := make(chan ChanResponse)

	for _, node := range switches {
		prefix := "http"

		if node.Ssl == true {
			prefix = prefix + "s"
		}
		url := prefix + "://" + node.Username + ":" + node.Password + "@" + node.Hostname + "/command-api"
		fmt.Println(url)
		go versionFetcher(url, cmds2, "text", node, c)
	}

	for i := 0; i < len(switches); i++ {
		<-c
	}
}

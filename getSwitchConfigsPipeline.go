package main

import (
	"encoding/json"
	"flag"
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
	Config        string
}


func readSwitches(filename string) []EosNode {
	var switches []EosNode

	file, err := os.Open("switches.json")
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&switches)
	if err != nil {
		panic(err)
	}
	return switches
}

func genSwitches(nodes []EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for _, node := range nodes {
			out <- node
		}
		close(out)
	}()
	return out
}

func buildUrl(node EosNode) string {
	prefix := "http"
	if node.Ssl == true {
		prefix = prefix + "s"
	}
	url := prefix + "://" + node.Username + ":" + node.Password + "@" + node.Hostname + "/command-api"
	return url
}

func getConfigs(in <-chan EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for n := range in {
			cmds := []string{"enable", "show running-config"}
			url := buildUrl(n)
			response := eapi.Call(url, cmds, "text")
			config := response.Result[1]["output"].(string)
			n.Config = config
			out <- n
		}
		close(out)
	}()
	return out
}

func getVersion(in <-chan EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for n := range in {
			cmds := []string{"show version"}
			url := buildUrl(n)
			response := eapi.Call(url, cmds, "json")
			version := response.Result[0]["version"].(string)
			n.Version = version
			out <- n
		}
		close(out)
	}()
	return out
}

func main() {
	switches := readSwitches("switches.json")

	fmt.Println("############# Using Pipelines ###################")
	c1 := genSwitches(switches)
	c2 := getConfigs(c1)
	out := getVersion(c2)
	for i := 0; i < len(switches); i++ {
		fmt.Println(<-out)
	}

}

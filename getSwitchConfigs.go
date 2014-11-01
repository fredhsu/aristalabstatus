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
}

type ChanResponse struct {
	response eapi.JsonRpcResponse
	node     EosNode
}

func configFetcher(url string, n EosNode, path string, c chan ChanResponse) {
	cmds := []string{"enable", "show running-config"}
	response := eapi.Call(url, cmds, "text")
	writeConfig(path, n, response.Result[1]["output"].(string))
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

func main() {
	swFilePtr := flag.String("swfile", "switches.json", "A JSON file with switches to fetch")
	pathPtr := flag.String("path", "/Users/fredlhsu/baseconfigs", "a directory to store the configs")
	flag.Parse() // command-line flag parsing
	fmt.Println(*swFilePtr)
	fmt.Println(*pathPtr)
	switches := readSwitches("switches.json")
	path := "/Users/fredlhsu/baseconfigs/"
	c := make(chan ChanResponse)

	for _, node := range switches {
		prefix := "http"
		if node.Ssl == true {
			prefix = prefix + "s"
		}
		url := prefix + "://" + node.Username + ":" + node.Password + "@" + node.Hostname + "/command-api"
		fmt.Println(url)
		go configFetcher(url, node, path, c)
	}

	for i := 0; i < len(switches); i++ {
		<-c
	}
}

package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
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

func versionFetcher(url string, cmds []string, format string, c chan eapi.JsonRpcResponse) {
	response := eapi.Call(url, cmds, format)
	c <- response
}

func main() {
	cmds2 := []string{"enable", "show running-config"}
	url1 := "https://admin:admin@bleaf1/command-api/"
	url2 := "https://admin:admin@bleaf2/command-api/"
	url3 := "https://admin:admin@bleaf3/command-api/"
	url4 := "https://admin:admin@bleaf5/command-api/"

	c := make(chan eapi.JsonRpcResponse)
	go versionFetcher(url1, cmds2, "text", c)
	go versionFetcher(url2, cmds2, "text", c)
	go versionFetcher(url3, cmds2, "text", c)
	go versionFetcher(url4, cmds2, "text", c)
	msg1 := <-c
	msg2 := <-c
	msg3 := <-c
	msg4 := <-c
	fmt.Println(msg1)
	fmt.Println(msg2)
	fmt.Println(msg3)
	fmt.Println(msg4)
	// Write these to files
}

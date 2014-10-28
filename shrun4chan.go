package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
)

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

	c1 := make(chan eapi.JsonRpcResponse)
	c2 := make(chan eapi.JsonRpcResponse)
	c3 := make(chan eapi.JsonRpcResponse)
	c4 := make(chan eapi.JsonRpcResponse)
	go versionFetcher(url1, cmds2, "text", c1)
	go versionFetcher(url2, cmds2, "text", c2)
	go versionFetcher(url3, cmds2, "text", c3)
	go versionFetcher(url4, cmds2, "text", c4)
	msg1 := <-c1
	msg2 := <-c2
	msg3 := <-c3
	msg4 := <-c4
	fmt.Println(msg1)
	fmt.Println(msg2)
	fmt.Println(msg3)
	fmt.Println(msg4)
}

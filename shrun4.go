package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
)

func main() {
	cmds2 := []string{"enable", "show running-config"}
	url1 := "https://admin:admin@bleaf1/command-api/"
	url2 := "https://admin:admin@bleaf2/command-api/"
	url3 := "https://admin:admin@bleaf3/command-api/"
	url4 := "https://admin:admin@bleaf5/command-api/"

	msg1 := eapi.Call(url1, cmds2, "text")
	msg2 := eapi.Call(url2, cmds2, "text")
	msg3 := eapi.Call(url3, cmds2, "text")
	msg4 := eapi.Call(url4, cmds2, "text")

	fmt.Println(msg1)
	fmt.Println(msg2)
	fmt.Println(msg3)
	fmt.Println(msg4)
}

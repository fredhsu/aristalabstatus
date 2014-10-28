package main

import (
	"fmt"
	"github.com/fredhsu/go-eapi"
)

func main() {
	cmds2 := []string{"enable", "show running-config"}
	url1 := "https://admin:admin@bleaf1/command-api/"

	msg1 := eapi.Call(url1, cmds2, "text")

	fmt.Println(msg1)
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fredhsu/go-eapi"
	"io/ioutil"
	"net/http"
	"os"
)

type EosNode struct {
	Hostname      string
    ModelName   string
	MgmtIp        string
	Username      string
	Password      string
	Ssl           bool
	Reachable     bool
	ConfigCorrect bool
	Uptime        float64
	Version       string
	Config        string
	IntfConnected []string
	IpIntf        []string
	Vlans         []string
    // LldpNeighbors   []string
    LldpNeighbors   []eapi.LldpNeighbor
}

// type LldpNeighbor struct {
//                 "ttl": 120,
//                "neighborDevice": "bleaf1.aristanetworks.com",
//                "neighborPort": "Ethernet1",
//                "port": "Ethernet1"
//     Port string
//     NeighborDevice string

// }

type ChanResponse struct {
	response eapi.JsonRpcResponse
	node     EosNode
}

func writeConfigFile(path string, n EosNode, config string) {
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
            result := response.Result[0]
			n.Version = result["version"].(string)
            n.ModelName = result["modelName"].(string)
			out <- n
		}
		close(out)
	}()
	return out
}

func getIntfConnected(in <-chan EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for n := range in {
			cmds := []string{"show interfaces status connected"}
			url := buildUrl(n)
			response := eapi.Call(url, cmds, "json")
			statuses := response.Result[0]["interfaceStatuses"].(map[string]interface{})
			for status := range statuses {
				n.IntfConnected = append(n.IntfConnected, status)
			}
			out <- n
		}
		close(out)
	}()
	return out
}

func getIpInterfaces(in <-chan EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for n := range in {
			cmds := []string{"show ip interface"}
			url := buildUrl(n)
			response := eapi.Call(url, cmds, "json")
			intfs := response.Result[0]["interfaces"].(map[string]interface{})
			for intf := range intfs {
				n.IntfConnected = append(n.IntfConnected, intf)
			}
			out <- n
		}
		close(out)
	}()
	return out
}

func getLldpNeighbors(in <-chan EosNode) <-chan EosNode {
    out := make(chan EosNode)
    go func() {
        for n := range in {
            cmds := []string{"show lldp neighbors"}
            data := eapi.RawCall(buildUrl(n), cmds, "json")
            var jsonresp map[string][]json.RawMessage
            err := json.Unmarshal(data, &jsonresp)
            if err != nil {
                fmt.Print("Json error")
            }
            var v eapi.ShowLldpNeighbors
            // var jsonresp2 []json.RawMessage

            json.Unmarshal(jsonresp["result"][0], &v)
            // fmt.Println(jsonresp2)
            // fmt.Println(jsonresp["result"])
            // fmt.Println(v)
            n.LldpNeighbors = v.LldpNeighbors
            // neighbors := (response.Result[0]["lldpNeighbors"]).([]interface{})
            // for _, neigh := range neighbors {
            //     i := neigh.(map[string]interface {})
            //     n.LldpNeighbors = append(n.LldpNeighbors, i["neighborDevice"].(string))
            // }
            out <- n
        }
        close(out)
    }()
    return out
}

func switchesHandler(w http.ResponseWriter, r *http.Request) {
	switches := readSwitches("switches.json")
	c1 := genSwitches(switches)
	c2 := getVersion(c1)
    c2 = getLldpNeighbors(c2)
	// c2 = getIntfConnected(c2)
	// c2 = getIpInterfaces(c2)
	output := []EosNode{}
	for i := 0; i < len(switches); i++ {
		node := <-c2
		fmt.Println(node)
		output = append(output, node)
	}

	b, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

func main() {
	// swFilePtr := flag.String("swfile", "switches.json", "A JSON file with switches to fetch")
	flag.Parse() // command-line flag parsing
	// switches := readSwitches(*swFilePtr)

	fmt.Println("############# Using Pipelines ###################")
	// c1 := genSwitches(switches)
	// c2 := getConfigs(c1)
	// c3 := getVersion(c2)
	// out := getIntfConnected(c3)
	// for i := 0; i < len(switches); i++ {
	// 	node := <-out
	// 	fmt.Print(node.Hostname + ": ")
	// 	fmt.Println(node.IntfConnected)
	// }
	http.HandleFunc("/switches/", switchesHandler)
	http.ListenAndServe(":8081", nil)
}

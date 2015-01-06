package main

import (
	"encoding/json"
	// "flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	lab "github.com/fredhsu/aristalabstatus"
	"github.com/fredhsu/eapigo"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type EosNode struct {
	Hostname      string
	ModelName     string
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
	LldpNeighbors []eapi.LldpNeighbor
}

type DemoStatus struct {
	Working bool
	Error   string
}

type ChanResponse struct {
	response eapi.JsonRpcResponse
	node     EosNode
}

type Link struct {
	Source   int `json:"source"`
	Target   int `json:"target"`
	Value    int `json:"value"`
	Distance int `json:"distance"`
}

type TopoData struct {
	Nodes []EosNode
	Links []Link
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

// Gets LLDP Neighbors and updates the EosNodes with the neighbors
func getLldpNeighbors(in <-chan EosNode) <-chan EosNode {
	out := make(chan EosNode)
	go func() {
		for n := range in {
			fmt.Println("getLldpNeighbors" + n.Hostname)
			cmds := []string{"show lldp neighbors"}
			data := eapi.RawCall(buildUrl(n), cmds, "json")
			// data := eapi.Call(buildUrl(n), cmds, "json")
			var jsonresp eapi.RawJsonRpcResponse
			// var jsonresp map[string]interface{}
			err := json.Unmarshal(data, &jsonresp)
			if err != nil {
				fmt.Print("Json error: ")
				fmt.Println(err)
			}
			// v := jsonresp.Result[0].(eapi.ShowLldpNeighbors)
			var v eapi.ShowLldpNeighbors
			// var jsonresp2 []json.RawMessage

			json.Unmarshal(jsonresp.Result[0], &v)
			// v := data.Result[0].(eapi.ShowLldpNeighbors)
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

// HTTP Handler for /switches
func switchesHandler(w http.ResponseWriter, r *http.Request) {
	switches := readSwitches("switches.json")

	c1 := genSwitches(switches)
	c2 := getVersion(c1)
	c2 = getLldpNeighbors(c2)
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

// Used to remove the FQDN to make names shorter and consistent
func removeFqdn(host string, domain string) string {
	return strings.TrimSuffix(host, "."+domain)
}

func topoHandler(w http.ResponseWriter, r *http.Request) {
	switches := readSwitches("switches.json")

	c1 := genSwitches(switches)
	c2 := getLldpNeighbors(c1)
	nodes := []EosNode{}
	sourceIds := map[string]int{}
	var links []Link
	for i := 0; i < len(switches); i++ {
		node := <-c2
		fmt.Println(node)
		nodes = append(nodes, node)
		sourceIds[node.Hostname] = i
		fmt.Println("sourceIds:")
		fmt.Println(sourceIds)
	}
	for i, node := range nodes {
		for _, l := range node.LldpNeighbors {
			target, ok := sourceIds[removeFqdn(l.NeighborDevice, "aristanetworks.com")]
			if !ok {
				fmt.Println("Not a valid neighbor: " + l.NeighborDevice)
			} else {
				fmt.Println("Link from " + node.Hostname + " to " + l.NeighborDevice)
				link := Link{
					Source:   i,
					Target:   target,
					Value:    1,
					Distance: 5,
				}
				// link["source"] = i
				// link["target"] = target
				// link["value"] = 1
				// link["distance"] = 5
				// fmt.Println(removeFqdn(l.NeighborDevice, "aristanetworks.com"))
				links = append(links, link)
			}
		}
	}
	output := TopoData{Nodes: nodes, Links: links}
	b, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Fprintf(w, string(b))
}

func panWebHandler(w http.ResponseWriter, r *http.Request) {
	svr := "172.22.28.143:8090"
	path := "/showinterfaces"
	url := "http://" + svr + path
	log.WithFields(log.Fields{
		"service": "panWebTest",
		"url":     url,
	}).Info("Testing Service")
	res, err := http.Get(url)
	var webStatus = DemoStatus{}
	if err != nil {
		webStatus.Working = false
		webStatus.Error = err.Error()
	} else {
		webStatus.Working = true
		defer res.Body.Close()
	}
	j, err := json.Marshal([]DemoStatus{webStatus})
	if err != nil {
		fmt.Println(err)
	}
	log.WithFields(log.Fields{
		"service": "panWebTest",
		"url":     url,
	}).Info("Finished Test")
	fmt.Fprintf(w, string(j))
}

func panHandler(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{
		"service": "panTest",
	}).Info("Starting Test")
	backupHost := "172.22.28.27"
	dosHost := "172.22.28.28"

	lab.PanResume()
	lab.PanClear()
	go lab.PingHost(dosHost)
	go lab.PingHost(backupHost)
	time.Sleep(30 * time.Second)

	// check flow entries
	bypassResult, bypassReason := lab.PanFlowTest("BYPASS")
	dropResult, dropReason := lab.PanFlowTest("DROP")
	// result := bypassResult + "\n" + dropResult
	lab.PanClear()

	bypassStatus := DemoStatus{Working: bypassResult, Error: bypassReason}
	dropStatus := DemoStatus{Working: dropResult, Error: dropReason}
	j, err := json.Marshal([]DemoStatus{bypassStatus, dropStatus})
	if err != nil {
		fmt.Fprintf(w, "Error encoding PAN demo response")
		return
	}
	log.WithFields(log.Fields{
		"service": "panTest",
	}).Info("Finsihed Test")
	fmt.Fprintf(w, string(j))
}

func genStatusJson(s []DemoStatus) string {
	j, err := json.Marshal(s)
	if err != nil {
		return "Error"
	}
	return string(j)
}

func openstackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	method := vars["method"]
	log.WithFields(log.Fields{
		"service": "openstack",
		"method":  method,
	}).Info("Starting Test")
	var s DemoStatus
	netname := "test_network_a"
	vmname := "test_vm_a"
	switch method {
	case "neutron":
		nclient := lab.GetNetworkClient()
		net := lab.CreateNetwork(nclient, netname)
		if net == nil {
			s = DemoStatus{Working: false, Error: "Could not create net"}
		} else {
			s = DemoStatus{Working: true, Error: ""}
		}
	case "subnet":
		nclient := lab.GetNetworkClient()
		found, net := lab.FindNetwork(nclient, netname)
		if !found {
			s = DemoStatus{Working: false, Error: "Could not find net"}
		} else {
			fmt.Println(net)
			sn := lab.CreateSubnet(nclient, "test_network_sn", net.ID, "192.168.87.0/24")
			if sn == nil {
				s = DemoStatus{Working: false, Error: "Could not create subnet"}
			} else {
				s = DemoStatus{Working: true, Error: ""}
			}
		}
	case "nova":
		nclient := lab.GetNetworkClient()
		cc := lab.GetComputeClient()
		found, net := lab.FindNetwork(nclient, netname)
		if !found {
			s = DemoStatus{Working: false, Error: "Could not find net"}
		} else {
			c := lab.CreateCompute(cc, vmname, net.ID)
			if c == nil {
				s = DemoStatus{Working: false, Error: "Could not create compute"}
			} else {
				s = DemoStatus{Working: true, Error: ""}
			}
		}
	case "eos":
		nclient := lab.GetNetworkClient()
		cc := lab.GetComputeClient()

		_, net := lab.FindNetwork(nclient, netname)
		_, compute := lab.FindCompute(cc, vmname)
		url := "https://admin:admin@bleaf1/command-api/"

		tn := lab.CheckNeutronEOS(url, net.ID)
		tc := lab.CheckNovaEOS(url, compute.ID)
		if tn && tc {
			s = DemoStatus{Working: true, Error: ""}
		} else {
			s = DemoStatus{Working: false, Error: "EOS info was not found"}
		}
	case "reset":
		nc := lab.GetNetworkClient()
		cc := lab.GetComputeClient()

		_, net := lab.FindNetwork(nc, netname)
		_, compute := lab.FindCompute(cc, vmname)
		dc := lab.DeleteCompute(cc, compute.ID)
		time.Sleep(5 * time.Second)
		dn := lab.DeleteNetwork(nc, net.ID)
		if dn && dc {
			s = DemoStatus{Working: true, Error: ""}
		} else {
			s = DemoStatus{Working: false, Error: "Error resetting demo"}
		}
	}
	fmt.Fprintf(w, genStatusJson([]DemoStatus{s}))

	return
}

// No longer needed since I'm running from Angular the same web server
// No CORS anymore!
func makeHandler(fn func(http.ResponseWriter, *http.Request, []EosNode), switches []EosNode) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		// Stop here if its Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}
		fn(w, r, switches)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	// swFilePtr := flag.String("swfile", "switches.json", "A JSON file with switches to fetch")
	// flag.Parse() // command-line flag parsing
	// switches := readSwitches(*swFilePtr)
	// fmt.Println(switches)
	// http.HandleFunc("/switch", func(w http.ResponseWriter, r *http.Request) {
	// 	switchesHandler(w, r, switches)
	// })
	// http.HandleFunc("/topo", func(w http.ResponseWriter, r *http.Request) {
	// 	topoHandler(w, r, switches)
	// })
	// http.HandleFunc("/pan", func(w http.ResponseWriter, r *http.Request) {
	// 	panHandler(w, r, switches)
	// })
	r := mux.NewRouter()
	// http.HandleFunc("/switches", makeHandler(switchesHandler, switches))
	//r.HandleFunc("/status", makeHandler(switchesHandler, switches))
	r.HandleFunc("/topo", topoHandler)
	r.HandleFunc("/status", switchesHandler)
	r.HandleFunc("/pan", panHandler)
	r.HandleFunc("/api/openstack/{method}", openstackHandler)
	r.HandleFunc("/panweb", panWebHandler)
	// http.HandleFunc("/panweb", func(w http.ResponseWriter, r *http.Request) {
	// 	panWebHandler(w, r)
	// })
	// r.HandleFunc("/panweb", makeHandler(panWebHandler, switches))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
	http.Handle("/", r)
	// r.Handle("/", http.FileServer(http.Dir(".")))
	http.ListenAndServe(":8081", r)
}

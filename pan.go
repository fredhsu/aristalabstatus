package aristalabstatus

import (
	"encoding/json"
	"fmt"
	"github.com/fredhsu/eapigo"
	"log"
	"net"
	"os/exec"
	"strings"
)

func CheckWebServer() string {
	// Check if web server is alive
	return "Pass"
}

func PanFlowTest(prefix string) (bool, string) {
	log.Printf("checking DFA flow entries")
	flows := FindFlows(prefix)
	if len(flows) > 0 {
		return true, ""
	}
	return false, "Flows not found"
}

func PanPause() {
	log.Printf("Pausing DFA")
	SendDatagram("DFA_CMD_PAUSE")
}

func PanResume() {
	log.Printf("Resuming DFA")
	SendDatagram("DFA_CMD_RESUME")
}

func SendDatagram(msg string) {
	DemoUdpPort := "9515"
	conn, err := net.Dial("udp", "172.22.28.95:"+DemoUdpPort)
	if err != nil {
		// handle error
	}
	fmt.Fprintf(conn, msg)
}

func PanClear() {
	log.Printf("Clearing DFA")
	SendDatagram("DFA_CMD_DELETE_FLOWS")
}

func FetchFlows() []eapi.Flow {
	cmds := []string{"show directflow flows"}
	url := "http://eapi:eapi@172.22.28.95/command-api"
	data := eapi.RawCall(url, cmds, "json")

	var jsonresp eapi.RawJsonRpcResponse
	err := json.Unmarshal(data, &jsonresp)
	if err != nil {
		log.Print("Json error: ")
		log.Println(err)
	}
	var v eapi.ShowDirectFlowFlows
	json.Unmarshal(jsonresp.Result[0], &v)
	return v.Flows
}

func FindFlows(prefix string) []eapi.Flow {
	log.Printf("Finding all flows starting with: " + prefix)
	found := []eapi.Flow{}
	flows := FetchFlows()
	if len(flows) == 0 {
		return nil
	}
	for _, flow := range flows {
		if strings.HasPrefix(flow.Name, prefix) {
			log.Println("Flow match: " + flow.Name)
			found = append(found, flow)
		}
	}
	return found
}

func RemoveAllFlows(prefix string) []eapi.Flow {
	log.Printf("Removing all DirectFlow Assist entries")
	flows := FetchFlows()
	// log.Println(flows)
	if len(flows) == 0 {
		return nil
	}
	cmds := []string{"configure", "directflow"}

	for _, flow := range flows {
		if strings.HasPrefix(flow.Name, prefix) {
			log.Println("Flow match: " + flow.Name)
		}
	}
	url := "http://eapi:eapi@172.22.28.95/command-api"
	response := eapi.Call(url, cmds, "json")
	log.Println(response.Result[0])
	// Check response for success

	return flows
}

func PingHost(host string) error {
	log.Printf("Pinging host " + host)
	cmd := exec.Command("ping", "-c 80", "-i .2", host)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Waiting for ping to finish")
	err = cmd.Wait()
	if err != nil {
		log.Printf("Ping error: %v", err)
	} else {
		log.Printf("Ping successful")
	}
	return err
}

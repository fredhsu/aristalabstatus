package aristalabstatus

import (
	"testing"
)

func TestCreateNetwork(t *testing.T) {
    client := GetNetworkClient()
    name := "test_create"
    net := CreateNetwork(client, name)
    if net == nil {
        t.Errorf("Got nil for netowrk creation")
    }
    if net.Name != name {
        t.Errorf("Network name is {}", net.ID)
    }
}

func TestFindNetwork(t *testing.T) {
    client := GetNetworkClient()
    name := "test_find"
    net := CreateNetwork(client, name)
    found, n := FindNetwork(client, net.ID)
    if !found {
        t.Errorf("Network not found")
    }
    if net.ID != n.ID {
        t.Errorf("Network name is %s, should be %s", net.ID, n.ID)
    }
}

func TestDeleteNetwork(t *testing.T) {
    client := GetNetworkClient()
    name := "test_remove"
    net := CreateNetwork(client, name)
    found, n := FindNetwork(client, net.ID)
    if !found {
        t.Errorf("Network not found")
    }
    if net.ID != n.ID {
        t.Errorf("Network name is %s, should be %s", net.ID, n.ID)
    }
    success := DeleteNetwork(client, net.ID)
    if !success {
        t.Errorf("Error in removing")
    }
}
// Create a cleanup test to remove all test networks

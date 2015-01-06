package aristalabstatus

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"testing"
	"time"
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
	DeleteNetwork(client, net.ID)
}

func TestCompute(t *testing.T) {
	nclient := GetNetworkClient()
	cclient := GetComputeClient()
	name := "test_create_vm"
	net := CreateNetwork(nclient, name)
	sn := CreateSubnet(nclient, "test_create_sn", net.ID, "192.168.187.0/24")
	if sn == nil {
		t.Errorf("Error creating subnet")
	}
	compute := CreateCompute(cclient, name, net.ID)
	if compute == nil {
		t.Errorf("Got nil for compute creation")
	}
	success := DeleteCompute(cclient, compute.ID)
	if !success {
		t.Errorf("Error in removing vm")
	}
	DeleteNetwork(nclient, net.ID)
}

func TestFindNetwork(t *testing.T) {
	client := GetNetworkClient()
	name := "test_find"
	net := CreateNetwork(client, name)
	found, n := FindNetwork(client, net.Name)
	if !found {
		t.Errorf("Network not found")
	}
	if net.Name != n.Name {
		t.Errorf("Network name is %s, should be %s", net.ID, n.ID)
	}
	DeleteNetwork(client, net.ID)
}
func TestFindCompute(t *testing.T) {
	nclient := GetNetworkClient()
	cclient := GetComputeClient()

	net, svrs := CreateNetCompute()
	name := "test_create_vm_1"
	found, n := FindCompute(cclient, name)
	if !found {
		t.Errorf("Server not found")
	}
	if name != n.Name {
		t.Errorf("Server name is %s, should be %s", n.Name, name)
	}
	for _, c := range svrs {
		DeleteCompute(cclient, c.ID)
	}
	time.Sleep(5 * time.Second)
	DeleteNetwork(nclient, net.ID)
}

func TestDeleteNetwork(t *testing.T) {
	client := GetNetworkClient()
	name := "test_remove"
	net := CreateNetwork(client, name)
	found, n := FindNetwork(client, net.Name)
	if !found {
		t.Errorf("Network not found")
	}
	if net.ID != n.ID {
		t.Errorf("Network ID created is %s, ID found is %s", net.ID, n.ID)
	}
	success := DeleteNetwork(client, net.ID)
	if !success {
		t.Errorf("Error in removing network")
	}
}

// Create a cleanup test to remove all test networks

func TestNeutronEOS(t *testing.T) {
	url := "https://admin:admin@bleaf1/command-api/"

	if CheckNeutronEOS(url, "bogus") {
		t.Errorf("Fake network created in EOS")
	}
	client := GetNetworkClient()
	name := "test_eos_net"
	net := CreateNetwork(client, name)
	if !CheckNeutronEOS(url, net.ID) {
		t.Errorf("Neutron network " + net.ID + " not created in EOS")
	}
	DeleteNetwork(client, net.ID)
}

func TestEOS(t *testing.T) {
	nclient := GetNetworkClient()
	cclient := GetComputeClient()

	url := "https://admin:admin@bleaf1/command-api/"
	if CheckNovaEOS(url, "bogus") {
		t.Errorf("Fake vm reported in EOS")
	}
	net, vms := CreateNetCompute()
	time.Sleep(3 * time.Second)
	for _, vm := range vms {
		if !CheckNovaEOS(url, vm.ID) {
			t.Errorf("VM " + vm.ID + " not created in EOS")
		}
		DeleteCompute(cclient, vm.ID)
	}
	time.Sleep(10 * time.Second) // Wait for compute nodes to spin up
	DeleteNetwork(nclient, net.ID)
}

func CreateNetCompute() (*networks.Network, []*servers.Server) {
	nclient := GetNetworkClient()
	cclient := GetComputeClient()
	name1 := "test_create_vm_1"
	name2 := "test_create_vm_2"
	net := CreateNetwork(nclient, "test_net")
	CreateSubnet(nclient, "test_create_sn", net.ID, "192.168.199.0/24")
	compute1 := CreateCompute(cclient, name1, net.ID)
	compute2 := CreateCompute(cclient, name2, net.ID)
	computes := []*servers.Server{compute1, compute2}
	return net, computes
}

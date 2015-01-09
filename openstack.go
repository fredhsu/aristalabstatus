package aristalabstatus

import (
	"fmt"
	"github.com/fredhsu/eapigo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/openstack/networking/v2/subnets"
	"github.com/rackspace/gophercloud/pagination"
)

func getProvider() *gophercloud.ProviderClient {
	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		fmt.Println(err)
	}
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		fmt.Println(err)
	}
	return provider
}

func GetNetworkClient() *gophercloud.ServiceClient {
	provider := getProvider()
	client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{
		Name:   "neutron",
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func GetComputeClient() *gophercloud.ServiceClient {
	provider := getProvider()
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func CreateNetwork(client *gophercloud.ServiceClient, name string) *networks.Network {
	//client := getNetworkClient()
	netopts := networks.CreateOpts{Name: name, AdminStateUp: networks.Up}
	network, err := networks.Create(client, netopts).Extract()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return network
}

func CreateSubnet(client *gophercloud.ServiceClient, name string, netid string, cidr string) *subnets.Subnet {
	opts := subnets.CreateOpts{
		NetworkID: netid,
		CIDR:      cidr,
		IPVersion: subnets.IPv4,
		Name:      name,
	}
	subnet, err := subnets.Create(client, opts).Extract()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return subnet
}

func CreateCompute(client *gophercloud.ServiceClient, name string, netid string) *servers.Server {
	var opts servers.CreateOpts
	if netid == "" {
		opts = servers.CreateOpts{
			Name:      name,
			ImageRef:  "7bc33401-17c4-4755-80d7-49ce1bf7d53d",
			FlavorRef: "42",
		}
	} else {
		net := servers.Network{UUID: netid}
		opts = servers.CreateOpts{
			Name:      name,
			ImageRef:  "7bc33401-17c4-4755-80d7-49ce1bf7d53d",
			FlavorRef: "42",
			Networks:  []servers.Network{net},
		}
	}
	compute, err := servers.Create(client, opts).Extract()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return compute
}

func FindNetwork(client *gophercloud.ServiceClient, name string) (bool, networks.Network) {
	shared := false
	listopts := networks.ListOpts{Shared: &shared}
	pager := networks.List(client, listopts)
	var net networks.Network
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, _ := networks.ExtractNetworks(page)
		for _, n := range networkList {
			if n.Name == name {
				net = n
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		fmt.Println(err)
		return false, net
	}
	return true, net
}

func FindCompute(client *gophercloud.ServiceClient, name string) (bool, servers.Server) {
	listopts := servers.ListOpts{Name: name}
	pager := servers.List(client, listopts)
	var s servers.Server
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, _ := servers.ExtractServers(page)
		for _, n := range serverList {
			if n.Name == name {
				s = n
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		fmt.Println(err)
		return false, s
	}
	return true, s
}

func DeleteNetwork(client *gophercloud.ServiceClient, uuid string) bool {
	result := networks.Delete(client, uuid)
	if result.ExtractErr() != nil {
		fmt.Println(result)
		return false
	}
	return true
}

func DeleteCompute(client *gophercloud.ServiceClient, uuid string) bool {
	result := servers.Delete(client, uuid)
	if result.ExtractErr() != nil {
		fmt.Println(result)
		return false
	}
	return true
}

// Checks if EOS is aware of the network

func CheckNeutronEOS(url string, uuid string) bool {
	// show openstack networks network <string>
	// if regions is not empty then network entry exists
	n := "show openstack networks network " + uuid
	cmds := []string{n}
	r := eapi.Call(url, cmds, "json")
	nets := r.Result[0]["regions"].(map[string]interface{})
	if len(nets) > 0 {
		return true
	}
	return false
}

func CheckNovaEOS(url string, uuid string) bool {
	// show openstack vms vm <string>
	// if regions is not empty then network entry exists
	n := "show openstack vms vm " + uuid
	cmds := []string{n}
	r := eapi.Call(url, cmds, "json")
	vms := r.Result[0]["regions"].(map[string]interface{})
	if len(vms) > 0 {
		return true
	}
	return false
}

func CheckRtrEOS() {
}

func CreateRouter() {
}

func RemoveRouter() {
}

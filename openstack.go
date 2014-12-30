package aristalabstatus

import (
      "github.com/rackspace/gophercloud"
      "github.com/rackspace/gophercloud/openstack"
      //"github.com/rackspace/gophercloud/openstack/utils"
      "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
      "github.com/rackspace/gophercloud/pagination"
      "fmt"
)

func getProvider() *gophercloud.ProviderClient{
    authOpts,err := openstack.AuthOptionsFromEnv()
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
    client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts {
        Name: "neutron",
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

func FindNetwork(client *gophercloud.ServiceClient, uuid string)(bool, networks.Network) {
    shared := false
    listopts := networks.ListOpts{Shared: &shared}
    pager := networks.List(client, listopts)
    var net networks.Network 
    err := pager.EachPage(func(page pagination.Page) (bool, error) {
        networkList, _ := networks.ExtractNetworks(page)
        for _, n := range networkList {
            fmt.Println(n)
            if n.ID == uuid {
                fmt.Println("!!! Network verified !!!")
                //result := networks.Delete(client, n.ID)
                //fmt.Println(result)
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

func DeleteNetwork(client *gophercloud.ServiceClient, uuid string) bool {
    result := networks.Delete(client, uuid)
    if result.ExtractErr() != nil {
        fmt.Println(result)
        return false
    }
    return true

}


func CheckNeutronEOS(uuids []string) {
}
func CheckNovaEOS(uuids []string) {
}
func CheckRtrEOS() {
}


func CreateRouter() {
}

func RemoveRouter() {
}

func CreateNova() {
}

func RemoveNova() {
}

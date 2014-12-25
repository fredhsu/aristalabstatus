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

func getNetworkClient() *gophercloud.ServiceClient {
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

func CreateNetworks() {
    client := getNetworkClient()
    netopts := networks.CreateOpts{Name: "test_network_199", AdminStateUp: networks.Up}
    network, err := networks.Create(client, netopts).Extract()
    fmt.Println(network)

    if err != nil {
        fmt.Println(err)
    }
}

func FindNetwork(uuid string) {
    client := getNetworkClient()
    shared := false
    listopts := networks.ListOpts{Shared: &shared}
    pager := networks.List(client, listopts)
    err := pager.EachPage(func(page pagination.Page) (bool, error) {
        networkList, _ := networks.ExtractNetworks(page)
        for _, n := range networkList {
            fmt.Println(n)
            if n.ID == uuid {
                fmt.Println("!!! Network verified !!!")
                //result := networks.Delete(client, n.ID)
                //fmt.Println(result)
            }
        }
        return true, nil
    })
    if err != nil {
        fmt.Println(err)
    }
}

func RemoveNetworks() {
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

package main

import (
      "github.com/rackspace/gophercloud"
      "github.com/rackspace/gophercloud/openstack"
      //"github.com/rackspace/gophercloud/openstack/utils"
      "github.com/rackspace/gophercloud/openstack/networking/v2/networks"
      "github.com/rackspace/gophercloud/pagination"
      "fmt"
)

func main() {
    opts,err := openstack.AuthOptionsFromEnv()
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(opts)
    provider, err := openstack.AuthenticatedClient(opts)
    if err != nil {
        fmt.Println(err)
    }
    client, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts {
        Name: "neutron",
        Region: "RegionOne",
    })
    if err != nil {
        fmt.Println(err)
    }
    shared := false
    netopts := networks.CreateOpts{Name: "test_network_199", AdminStateUp: networks.Up}
    network, err := networks.Create(client, netopts).Extract()
    fmt.Println(network)

    listopts := networks.ListOpts{Shared: &shared}
    pager := networks.List(client, listopts)
    err = pager.EachPage(func(page pagination.Page) (bool, error) {
        networkList, _ := networks.ExtractNetworks(page)
        for _, n := range networkList {
            fmt.Println(n)
            if n.ID == network.ID {
                fmt.Println("!!! Network verified !!!")
                //result := networks.Delete(client, n.ID)
                //fmt.Println(result)
            }
        }
        return true, nil
    })
}



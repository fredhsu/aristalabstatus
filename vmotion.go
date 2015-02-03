package main

import (
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"net/url"
)

func GetHost(id int, folders *govmomi.DatacenterFolders, c *govmomi.Client) (types.ManagedObjectReference, *types.ManagedObjectReference) {
	crs, err := folders.HostFolder.Children()
	if err != nil {
		fmt.Println("Error")
	}
	var cr mo.ComputeResource
	// 2 = .218
	// 4 = .146
	// 7 = .150
	ref := crs[id]
	err = c.Properties(ref.Reference(), nil, &cr)
	if err != nil {
		fmt.Println("Error:")
		fmt.Println(err.Error())
	}
	return cr.Host[0], cr.ResourcePool
}

func GetVm(name string, folders *govmomi.DatacenterFolders, c *govmomi.Client) mo.VirtualMachine {
	vms, err := folders.VmFolder.Children()
	if err != nil {
		fmt.Println("Error")
	}
	var vm mo.VirtualMachine
	for _, ref := range vms {
		err = c.Properties(ref.Reference(), []string{"config", "guest"}, &vm)
		if err == nil && vm.Config.Name == name {
			return vm
		}
	}
	return vm
}

func GetFolders(c *govmomi.Client) *govmomi.DatacenterFolders {
	s := c.SearchIndex()
	ref, err := s.FindChild(c.RootFolder(), "SEDEMO")
	dc, ok := ref.(*govmomi.Datacenter)
	if !ok {
		fmt.Println("DC error")
	}

	folders, err := dc.Folders()
	if err != nil {
		fmt.Println("Error")
	}
	return folders
}

func GetClient() *govmomi.Client {
	u, err := url.Parse("https://root:vmware@172.22.28.190/sdk")
	if err != nil {
		fmt.Println("URL parse error")
	}
	c, err := govmomi.NewClient(*u, true)
	if err != nil {
		fmt.Println("Connect error", err)
	}
	return c
}

func MigrateVm(vm mo.VirtualMachine, pool *types.ManagedObjectReference, host *types.ManagedObjectReference, c *govmomi.Client) {
	mvmtask := types.MigrateVM_Task{
		This:     vm.Reference(),
		Pool:     pool,
		Host:     host,
		Priority: types.VirtualMachineMovePriorityDefaultPriority}
	methods.MigrateVM_Task(c, &mvmtask)
}

// Test if VM named vm shows up on switch at url
func VmtracerTest(vm string, url string) bool {
	return false
}

func main() {
	c := GetClient()
	folders := GetFolders(c)
	vm := GetVm("vxlan-vm1", folders, c)
	// vMotion to .146 verify
	dHost, pool := GetHost(4, folders, c) // Use 4 to move to .146
	MigrateVm(vm, pool, &dHost, c)

	// vMotion to .150 verify
	dHost, pool = GetHost(7, folders, c) // Use 7 to move to .150
	// MigrateVm(vm, pool, &dHost, c)
}

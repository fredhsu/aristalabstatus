package aristalabstatus

import (
	"fmt"
	"net/url"
	"time"

	"github.com/fredhsu/eapigo"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
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
	fmt.Println("Fetched vm: ", vm)
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
	n := "show vmtracer vm " + vm
	cmds := []string{n}
	r := eapi.Call(url, cmds, "json")
	vms := r.Result[0]["vms"].([]interface{})
	if len(vms) > 0 {
		return true
	}
	return false
}

func main() {
	c := GetClient()
	folders := GetFolders(c)

	vm := GetVm("vxlan-vm1", folders, c)
	// vMotion to .146 verify
	dHost, pool := GetHost(4, folders, c) // Use 4 to move to .146
	MigrateVm(vm, pool, &dHost, c)
	// Wait for vMotion to finish
	time.Sleep(10 * time.Second) // Wait for vMotion to finish
	url := "https://admin:admin@172.28.171.101/command-api/"
	result := VmtracerTest("vxlan-vm1", url)
	fmt.Printf("VM is on leaf 1: %t", result)
	url = "https://admin:admin@172.28.171.102/command-api/"
	result = VmtracerTest("vxlan-vm1", url)

	fmt.Printf("VM is on leaf 2: %t", result)

	// vMotion to .150 verify
	dHost, pool = GetHost(7, folders, c) // Use 7 to move to .150
	MigrateVm(vm, pool, &dHost, c)
	url = "https://admin:admin@172.28.171.101/command-api/"
	time.Sleep(10 * time.Second) // Wait for vMotion to finish

	result = VmtracerTest("vxlan-vm1", url)
	fmt.Println(result)
	fmt.Printf("VM is on leaf 1: %t", result)
	url = "https://admin:admin@172.28.171.102/command-api/"
	result = VmtracerTest("vxlan-vm1", url)

	fmt.Printf("VM is on leaf 2: %t", result)

}

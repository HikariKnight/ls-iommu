package iommu

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/HikariKnight/ls-iommu/lib/params"
	"github.com/jaypipes/ghw"
	ghwpci "github.com/jaypipes/ghw/pkg/pci"
)

type IOMMU struct {
	Groups map[int]*Group
}

func (i *IOMMU) AddGroup(group *Group) {
	i.Groups[group.ID] = group
}

type Group struct {
	ID      int
	Devices map[string]*ghwpci.Device
}

func (g *Group) AddDevice(device *ghwpci.Device) {
	g.Devices[(*device).Address] = device
}

func NewGroup(id int, devices map[string]*ghwpci.Device) *Group {
	return &Group{
		ID:      id,
		Devices: devices,
	}
}

func (i *IOMMU) Read() {
	i.Groups = make(map[int]*Group)
	// Get all groups and associated devices
	iommu_devices, err := filepath.Glob("/sys/kernel/iommu_groups/*/devices/*")
	ErrorCheck(err)
	pci, err := ghw.PCI(ghw.WithDisableWarnings())
	ErrorCheck(err)

	// Regex to get IOMMU groups and their devices from filepath
	iommu_regex := regexp.MustCompile(`/sys/kernel/iommu_groups/(.*)/devices/(.*)`)
	for _, iommu_device := range iommu_devices {
		matches := iommu_regex.FindStringSubmatch(iommu_device)
		group_id, err := strconv.Atoi(matches[1])
		if err != nil {
			// Failed to properly parse groupid into integer, invalid for a group id, skip it
			continue
		}
		device_id := matches[2]

		r := regexp.MustCompile(`^([a-z0-9]{1,4}):`)

		if r.MatchString(device_id) {
			device := pci.GetDevice(device_id)
			// If the group doesn't exist in the struct, add it
			_, exists := i.Groups[group_id]
			if !exists {
				/*
					grp := &Group{
						ID:      group_id,
						Devices: make(map[string]*ghwpci.Device),
					}
				*/
				// Make a new Group struct, this is equal to the code above for reference
				grp := NewGroup(group_id, make(map[string]*ghwpci.Device))

				// Add the device to the group
				grp.AddDevice(device)

				// Add the group to the IOMMU struct
				i.AddGroup(grp)
			} else {
				i.Groups[group_id].AddDevice(device)
			}
		}
	}
}

func NewIOMMU() *IOMMU {
	iommu := &IOMMU{}
	iommu.Read()
	return iommu
}

func GetAllDevices(pArg *params.Params) []string {
	iommu := NewIOMMU()
	var lspci_devs []string

	// Iterate through the IOMMU groups and get the device info
	for id := 0; id < len(iommu.Groups); id++ {
		// Iterate each device
		for _, device := range iommu.Groups[id].Devices {
			// Generate the device list with the data we want
			line := generateDevList(id, device, pArg)
			lspci_devs = append(lspci_devs, line)
		}
	}

	return lspci_devs
}

func MatchSubclass(searchval string, pArg *params.Params) []string {
	var devs []string

	// Get all IOMMU devices
	alldevs := NewIOMMU()

	// Iterate through the groups
	for id := 0; id < len(alldevs.Groups); id++ {
		// For each device
		for _, device := range alldevs.Groups[id].Devices {
			// If the device has a subclass matching what we are looking for
			if strings.Contains(device.Subclass.Name, searchval) {
				// Generate the device list with the data we want
				line := generateDevList(id, device, pArg)
				devs = append(devs, line)

				// If we want to search for related devices
				if pArg.FlagCounter["related"] > 0 && searchval != `USB controller` {
					// Find relatives and add them to the list
					related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
					devs = append(devs, related_list...)
				} else if pArg.FlagCounter["related"] > 0 && searchval == `USB controller` {
					// Prevent an infinite loop by passing 0 instead of related
					other := GetDevicesFromGroups([]int{id}, 0, pArg)
					devs = append(devs, other...)
				}
			}
		}
	}

	return devs
}

// Function to print everything inside a specific IOMMU group
func GetDevicesFromGroups(groups []int, related int, pArg *params.Params) []string {
	// Make an output string slice
	var output []string

	// As long as we are asked to get devices from any specific IOMMU groups
	if len(groups) > 0 {
		// Get all IOMMU devices
		alldevs := NewIOMMU()

		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			// Check if the IOMMU Group exists
			if _, iommu_num := alldevs.Groups[group]; !iommu_num {
				ErrorCheck(fmt.Errorf("IOMMU Group %v does not exist", group))
			} else {
				// For each device in specified IOMMU group
				for _, device := range alldevs.Groups[group].Devices {
					// Generate the device list with the data we want
					line := generateDevList(group, device, pArg)

					// Append line to output
					output = append(output, line)

					if related > 0 {
						// Find relatives and add them to the list
						related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
						output = append(output, related_list...)
					}
				}
			}
		}
	}
	return output
}

func findRelatedDevices(vendorid string, related int, pArg *params.Params) []string {
	var devs []string

	// Get all IOMMU devices
	alldevs := NewIOMMU()

	// Iterate through the groups
	for id := 0; id < len(alldevs.Groups); id++ {
		// For each device
		for _, device := range alldevs.Groups[id].Devices {
			// If the device has a vendor ID matching what we are looking for
			if strings.Contains(device.Vendor.ID, vendorid) {
				// Generate the device list with the data we want
				line := generateDevList(id, device, pArg)
				devs = append(devs, line)

				if related > 1 {
					// Prevent an infinite loop by passing 0 instead of related variable
					other := GetDevicesFromGroups([]int{id}, 0, pArg)
					devs = append(devs, other...)
				}
			}
		}
	}

	return devs
}

// Old deprecated functions marked for removal/rework below this comment

func MatchDEVs(regex string, pArg *params.Params) []string {
	var devs []string

	output := GetAllDevices(pArg)
	gpuReg, err := regexp.Compile(regex)
	ErrorCheck(err)

	for _, line := range output {
		if gpuReg.MatchString(line) {
			devs = append(devs, line)
		}
	}

	return devs
}
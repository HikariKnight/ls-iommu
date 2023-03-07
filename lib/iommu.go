package iommu

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

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
				grp := &Group{
					ID:      group_id,
					Devices: make(map[string]*ghwpci.Device),
				}
				grp.AddDevice(device)
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

func GetAllDevices(kernelmodules ...bool) []string {
	iommu := NewIOMMU()
	var lspci_devs []string

	// Iterate through the IOMMU groups and get the device info
	for id := 0; id < len(iommu.Groups); id++ {
		// Iterate each device
		for _, device := range iommu.Groups[id].Devices {
			// Generate the device list with the data we want
			line := generateDevList(id, device, len(kernelmodules))
			lspci_devs = append(lspci_devs, line)
		}
	}

	return lspci_devs
}

func MatchSubclass(searchval string, kernelmodules ...bool) []string {
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
				line := generateDevList(id, device, len(kernelmodules))
				devs = append(devs, line)
			}
		}
	}

	return devs
}

// Function to print everything inside a specific IOMMU group
func GetDevicesFromGroups(groups []int, kernelmodules ...bool) []string {
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
					line := generateDevList(group, device, len(kernelmodules))

					// Append line to output
					output = append(output, line)
				}
			}
		}
	}
	return output
}

// Old deprecated functions marked for removal/rework below this comment

func MatchDEVs(kernelmodules bool, regex string) []string {
	var devs []string

	output := GetAllDevices(kernelmodules)
	gpuReg, err := regexp.Compile(regex)
	ErrorCheck(err)

	for _, line := range output {
		if gpuReg.MatchString(line) {
			devs = append(devs, line)
		}
	}

	return devs
}

package iommu

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/ls-iommu/pkg/params"
	"github.com/jaypipes/ghw"
	ghwpci "github.com/jaypipes/ghw/pkg/pci"
)

type IOMMU struct {
	Groups map[int]*Group
}

// Adds a Group struct to the IOMMU struct
func (i *IOMMU) AddGroup(group *Group) {
	i.Groups[group.ID] = group
}

type Group struct {
	ID      int
	Devices map[string]*ghwpci.Device
}

// Adds a new pci device to the Group struct
func (g *Group) AddDevice(device *ghwpci.Device) {
	g.Devices[(*device).Address] = device
}

// Creates a new Group struct
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
	errorcheck.ErrorCheck(err)
	pci, err := ghw.PCI(ghw.WithDisableWarnings())
	errorcheck.ErrorCheck(err)

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

		// Regex to check for a valid PCI domain to avoid a
		// ghw bug with Intel VMA which uses the PCI domain 10000, which ghw cannot handle
		r := regexp.MustCompile(`^([a-z0-9]{1,4}):`)

		// Only match valid PCI domains (start with 4 hexadecimal characters followed by a :)
		if r.MatchString(device_id) {
			device := pci.GetDevice(device_id)
			// If the group doesn't exist in the struct, add it
			_, exists := i.Groups[group_id]

			// If the group does not exist in our struct
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
				// Add the device to the existing group ID
				i.Groups[group_id].AddDevice(device)
			}
		}
	}

	// If we have 0 groups so far, IOMMU is probably disabled
	if len(i.Groups) == 0 {
		log.Fatal("IOMMU Disabled in UEFI/BIOS and/or not enabled in boot arguments!")
	}
}

// Creates an IOMMU struct-
func NewIOMMU() *IOMMU {
	// Make an empty IOMMU struct
	iommu := &IOMMU{}

	// Get all the IOMMU data
	iommu.Read()

	// Return the struct with the data
	return iommu
}

func GetAllDevices(pArg *params.Params) []string {
	// Get all the IOMMU data and put it into a variable
	iommu := NewIOMMU()

	// Prepare a string slice for storing our output
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

	// Make a regex to find devices that will need special exceptions for relative search
	specialRelativeExceptions := regexp.MustCompile(`^(SATA|USB) controller`)

	// Iterate through the groups
	for id := 0; id < len(alldevs.Groups); id++ {
		// For each device
		for _, device := range alldevs.Groups[id].Devices {
			// If the device has a subclass matching what we are looking for
			if strings.Contains(device.Subclass.Name, searchval) {
				if len(pArg.IntList["iommu_group"]) == 0 && !pArg.Flag["rom"] {
					// Generate the device list with the data we want
					line := generateDevList(id, device, pArg)
					devs = append(devs, line)

					// If we want to search for related devices
					if pArg.FlagCounter["related"] > 0 && !specialRelativeExceptions.MatchString(searchval) {
						// Find relatives and add them to the list
						related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
						devs = append(devs, related_list...)

					} else if pArg.FlagCounter["related"] > 0 && specialRelativeExceptions.MatchString(searchval) {
						// Prevent an infinite loop by passing 0 instead of related
						other := GetDevicesFromGroups([]int{id}, 0, pArg)
						devs = append(devs, other...)
					}

				} else if pArg.Flag["rom"] && pArg.Flag["gpu"] {
					// If we are asked to get the path to the gpu vbios
					if len(pArg.IntList["iommu_group"]) > 0 {
						// If we are asked to only get a specific IOMMU group
						for _, group := range pArg.IntList["iommu_group"] {
							// If the iommu group matches the one we are currently processing
							if id == group {
								// Print the GPU rom path
								devs = append(devs, GetRomPath(device, pArg)...)
							}
						}

					} else {
						// Else get the vbios path for any gpu
						devs = append(devs, GetRomPath(device, pArg)...)
					}

				} else {
					for _, group := range pArg.IntList["iommu_group"] {
						if id == group {
							// If we want the Device ID or PCI Address
							if pArg.Flag["id"] && !pArg.Flag["pciaddr"] {
								// If --id is supplied as an argument we display the VendorID:DeviceID
								devs = append(devs, fmt.Sprintf("%s:%s\n", device.Vendor.ID, device.Product.ID))

							} else if !pArg.Flag["id"] && pArg.Flag["pciaddr"] {
								// If --pciaddr is supplied as an argument we display the PCI Address
								devs = append(devs, fmt.Sprintf("%s\n", device.Address))

							} else {
								// Generate the device list with the data we want
								line := generateDevList(id, device, pArg)
								devs = append(devs, line)
							}

							// If we want to search for related devices
							if pArg.FlagCounter["related"] > 0 && !specialRelativeExceptions.MatchString(searchval) {
								// Find relatives and add them to the list
								related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
								devs = append(devs, related_list...)

							} else if pArg.FlagCounter["related"] > 0 && specialRelativeExceptions.MatchString(searchval) {
								// Prevent an infinite loop by passing 0 instead of related
								other := GetDevicesFromGroups([]int{id}, 0, pArg)
								devs = append(devs, other...)
							}
						}
					}
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
				errorcheck.ErrorCheck(fmt.Errorf("IOMMU Group %v does not exist", group))

			} else {
				// For each device in specified IOMMU group
				for _, device := range alldevs.Groups[group].Devices {
					// If we do not want the Device IDs or PCI Address
					if !pArg.Flag["id"] && !pArg.Flag["pciaddr"] {
						// Generate the device list with the data we want
						line := generateDevList(group, device, pArg)

						// Append line to output
						output = append(output, line)

						if related > 0 {
							// Find relatives and add them to the list
							related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
							output = append(output, related_list...)
						}

					} else if !strings.Contains(device.Subclass.Name, "bridge") {
						if pArg.Flag["id"] && !pArg.Flag["pciaddr"] {
							// If --id is supplied as an argument we display the VendorID:DeviceID
							output = append(output, fmt.Sprintf("%s:%s\n", device.Vendor.ID, device.Product.ID))

							if related > 0 {
								// Find relatives and add them to the list
								related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
								output = append(output, related_list...)
							}

						} else if !pArg.Flag["id"] && pArg.Flag["pciaddr"] {
							// If --pciaddr is supplied as an argument we display the PCI Address
							output = append(output, fmt.Sprintf("%s\n", device.Address))

							if related > 0 {
								// Find relatives and add them to the list
								related_list := findRelatedDevices(device.Vendor.ID, pArg.FlagCounter["related"], pArg)
								output = append(output, related_list...)
							}
						}
					}
				}
			}
		}
	}
	return output
}

// Find related devices based on VendorID, and do a deeper search in the same IOMMU group if specified
func findRelatedDevices(vendorid string, related int, pArg *params.Params) []string {
	// Make a string slice for our output
	var devs []string

	// Get all IOMMU devices
	alldevs := NewIOMMU()

	// Iterate through the groups
	for id := 0; id < len(alldevs.Groups); id++ {
		// For each device
		for _, device := range alldevs.Groups[id].Devices {
			// If the device has a vendor ID matching what we are looking for
			if strings.Contains(device.Vendor.ID, vendorid) {
				// Make a variable to decide if device should be ignored
				ignoreDevice := false

				// Check if we should ignore the device
				for _, ignore := range pArg.StringList["ignore"] {
					// If the VendorID is to be ignored
					if ignore == device.Vendor.ID {
						// Mark device to be ignored
						ignoreDevice = true
					}
				}

				if !ignoreDevice {
					// If we do not want the Device IDs or PCI Address
					if !pArg.Flag["id"] && !pArg.Flag["pciaddr"] {
						// Generate the device list with the data we want
						line := generateDevList(id, device, pArg)
						devs = append(devs, line)

						if related > 1 {
							// Prevent an infinite loop by passing 0 instead of related variable
							other := GetDevicesFromGroups([]int{id}, 0, pArg)
							devs = append(devs, other...)
						}

					} else if !strings.Contains(device.Subclass.Name, "bridge") {
						if pArg.Flag["id"] && !pArg.Flag["pciaddr"] {
							// If --id is supplied as an argument we display the VendorID:DeviceID
							devs = append(devs, fmt.Sprintf("%s:%s\n", device.Vendor.ID, device.Product.ID))
						} else if !pArg.Flag["id"] && pArg.Flag["pciaddr"] {
							// If --pciaddr is supplied as an argument we display the PCI Address
							devs = append(devs, fmt.Sprintf("%s\n", device.Address))
						}
					}
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
	errorcheck.ErrorCheck(err)

	for _, line := range output {
		if gpuReg.MatchString(line) {
			devs = append(devs, line)
		}
	}

	return devs
}

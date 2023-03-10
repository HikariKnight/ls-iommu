package iommu

import (
	"fmt"
	"log"
	"sort"

	"github.com/HikariKnight/ls-iommu/lib/params"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/pci"
)

// Generates a line with the Device info and formats it properly to be similar to the bash version of ls-iommu
func GenDeviceLine(group int, device *pci.Device, pArg *params.Params) string {
	var line string

	// If we want legacy output (to be output compatible with the bash version)
	var iommu_group string
	if pArg.Flag["legacyoutput"] {
		// Do not pad the group number
		iommu_group = fmt.Sprintf("%d", group)
	} else {
		// Else we pad the group number to make it sortable
		iommu_group = fmt.Sprintf("% 3d", group)
	}

	// If the device has no revision, ommit the (rev ID), in both cases we generate the line with device info
	if device.Revision != "0x00" {
		line = fmt.Sprintf("IOMMU Group %s: %s %s [%s%s]: %s %s [%s:%s] (rev %s)\n",
			iommu_group,
			device.Address,
			device.Subclass.Name,
			device.Class.ID,
			device.Subclass.ID,
			device.Vendor.Name,
			device.Product.Name,
			device.Vendor.ID,
			device.Product.ID,
			device.Revision[len(device.Revision)-2:],
		)
	} else {
		line = fmt.Sprintf("IOMMU Group %s: %s %s [%s%s]: %s %s [%s:%s]\n",
			iommu_group,
			device.Address,
			device.Subclass.Name,
			device.Class.ID,
			device.Subclass.ID,
			device.Vendor.Name,
			device.Product.Name,
			device.Vendor.ID,
			device.Product.ID,
		)
	}

	return line
}

// Generates the kernel driver info for a device
func GenKernelInfo(group int, device *pci.Device) string {
	var line string
	var subsystem_name string
	var subvendor_name string

	// We need to probe some extra info here so we need to use ghw
	pci, err := ghw.PCI()
	if err != nil {
		fmt.Printf("Error getting PCI info: %v", err)
	}

	// Get the subvendor
	subvendor := pci.Vendors[device.Subsystem.VendorID]

	// If subvendor does exist
	if subvendor != nil {
		// Get the subvendor name
		subvendor_name = pci.Vendors[device.Subsystem.VendorID].Name
	} else {
		// Else slap the vendor name on
		subvendor_name = device.Vendor.Name
	}

	// If the subsystem name is unknown then use the product name instead
	if device.Subsystem.Name == "unknown" {
		subsystem_name = device.Product.Name
	} else {
		subsystem_name = device.Subsystem.Name
	}

	// Add the subSystemID to a string so we can check if its valid
	subSystemID := fmt.Sprintf("%s:%s", device.Subsystem.VendorID, device.Subsystem.ID)

	// If we have a valid (not just 0s) ID
	if subSystemID != "0000:0000" {
		// Add the Subsystem data
		line = fmt.Sprintf(
			"\tSubsystem: %s %s [%s:%s]\n",
			subvendor_name,
			subsystem_name,
			device.Subsystem.VendorID,
			device.Subsystem.ID,
		)
	}

	// If we do not have an empty driver string
	if device.Driver != "" {
		// Add the driver data
		line = fmt.Sprintf("%s\tKernel driver in use: %s\n",
			line,
			device.Driver,
		)
	}

	return line
}

// Generates a line for our device list
func generateDevList(id int, device *pci.Device, pArg *params.Params) string {
	var line string

	// If user requested kernel modules
	if pArg.Flag["kernelmodules"] {
		// Generate the line with kernel modules
		line = fmt.Sprintf(
			"%s%s",
			GenDeviceLine(id, device, pArg),
			GenKernelInfo(id, device),
		)
	} else {
		// Generate the line without the kernel modules
		line = GenDeviceLine(id, device, pArg)
	}

	return line
}

// Function to just print out a string array to STDOUT
func PrintOutput(out []string, pArg *params.Params) {
	if len(out) == 0 && (!pArg.Flag["id"] && !pArg.Flag["pciaddr"]) {
		log.Fatal("IOMMU disabled in UEFI/BIOS and/or you have not configured your\n\t\t    bootloader to enable IOMMU with the kernel boot arguments!")
	} else if len(out) == 0 && (pArg.Flag["id"] || pArg.Flag["pciaddr"]) {
		log.Fatal("--id or --pciaddr detected however no Device IDs or PCI Addresses were printed\nThis could be because there were only bridge devices found.\nIf other flags also results in a blank output then you do not have IOMMU enabled in UEFI/BIOS or you are missing the kernel flags in your bootloader!")
	}

	// Remove duplicate lines
	output := removeDuplicateLines(out)
	// Sort cleaned output
	sort.Strings(output)

	// Print output line by line
	for _, line := range output {
		fmt.Print(line)
	}
}

// Removes duplicate lines from a string slice, useful for cleaning up the output if doing multiple scans
func removeDuplicateLines(lines []string) []string {
	// Make a map to keep track of which strings have been processed
	keys := make(map[string]bool)

	// Make a new string slice
	var list []string

	// For each line
	for _, entry := range lines {
		// If the line has not been processed before
		if _, value := keys[entry]; !value {
			// Mark it as processed in our map
			keys[entry] = true

			// Add line to our list
			list = append(list, entry)
		}
	}
	return list
}

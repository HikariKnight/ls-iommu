package main

import (
	"fmt"
	"log"
	"os"
	"sort"

	iommu "github.com/HikariKnight/ls-iommu/lib/iommu"
	params "github.com/HikariKnight/ls-iommu/lib/params"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	// Work with the output depending on arguments given
	if pArg.Flag["gpu"] {
		// Get all GPUs (3d controllers are ignored)
		output := iommu.MatchSubclass(`VGA`, pArg)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(pArg.IntList["iommu_group"], pArg.FlagCounter["related"], pArg)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if pArg.Flag["usb"] {
		// Get all USB controllers
		output := iommu.MatchSubclass(`USB controller`, pArg)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(pArg.IntList["iommu_group"], pArg.FlagCounter["related"], pArg)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if pArg.Flag["nic"] {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`Ethernet controller`, pArg)

		// Get all Wi-Fi controllers
		wifi := iommu.MatchSubclass(`Network controller`, pArg)
		output = append(output, wifi...)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(pArg.IntList["iommu_group"], pArg.FlagCounter["related"], pArg)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if len(pArg.IntList["iommu_group"]) > 0 {
		// Get all devices in specified IOMMU groups and append it to the output
		output := iommu.GetDevicesFromGroups(pArg.IntList["iommu_group"], pArg.FlagCounter["related"], pArg)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(pArg)
		printoutput(out)
	}
}

// Function to just print out a string array to STDOUT
func printoutput(out []string) {
	if len(out) == 0 {
		log.Fatal("IOMMU disabled in UEFI/BIOS and/or you have not configured your\n\t\t    bootloader to enable IOMMU with the kernel boot arguments!")
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

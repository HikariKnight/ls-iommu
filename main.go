package main

import (
	"fmt"
	"os"
	"sort"

	iommu "github.com/HikariKnight/ls-iommu/lib"
	"github.com/akamensky/argparse"
)

func main() {
	// Setup the parser for arguments
	parser := argparse.NewParser("ls-iommu", "A Tool to print out all devices and their IOMMU groups")

	// Configure arguments
	gpu := parser.Flag("g", "gpu", &argparse.Options{
		Required: false,
		Help:     "List all GPUs.",
	})

	usb := parser.Flag("u", "usb", &argparse.Options{
		Required: false,
		Help:     "List all USB controllers.",
	})

	nic := parser.Flag("n", "network", &argparse.Options{
		Required: false,
		Help:     "List all Network controllers.",
	})

	related := parser.FlagCounter("r", "related", &argparse.Options{
		Required: false,
		Help:     "Attempt to list related devices that share Vendor ID or\n\t\t IOMMU Groups (used with -g -u and -n), pass -rr if you want to search using both when used with -g or -n\n\t\t Note: -rr can be inaccurate or too broad when many devices share Vendor ID",
	})

	iommu_group := parser.IntList("i", "group", &argparse.Options{
		Required: false,
		Help:     "List everything in the IOMMU groups given. Supply argument multiple times to list additional groups.",
	})

	kernelmodules := parser.Flag("k", "kernel", &argparse.Options{
		Required: false,
		Help:     "Lists subsystems and kernel drivers using the devices.",
		Default:  false,
	})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(4)
	}

	// Work with the output depending on arguments given
	if *gpu {
		// Get all GPUs (3d controllers are ignored)
		output := iommu.MatchSubclass(`VGA`, *related, *kernelmodules)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(*iommu_group, *kernelmodules)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if *usb {
		// Get all USB controllers
		output := iommu.MatchSubclass(`USB controller`, *related, *kernelmodules)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(*iommu_group, *kernelmodules)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if *nic {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`Ethernet controller`, *related, *kernelmodules)

		// Get all Wi-Fi controllers
		wifi := iommu.MatchSubclass(`Network controller`, *related, *kernelmodules)
		output = append(output, wifi...)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(*iommu_group, *kernelmodules)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if len(*iommu_group) > 0 {
		// Get all devices in specified IOMMU groups and append it to the output
		output := iommu.GetDevicesFromGroups(*iommu_group, *kernelmodules)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(*kernelmodules)
		printoutput(out)
	}
}

// Function to just print out a string array to STDOUT
func printoutput(out []string) {
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

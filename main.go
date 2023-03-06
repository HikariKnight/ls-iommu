package main

import (
	"fmt"
	"os"
	"strings"

	iommu "github.com/HikariKnight/ls-iommu/lib"
	"github.com/akamensky/argparse"
)
func main() {
	// Setup the parser for arguments
	parser := argparse.NewParser("ls-iommu", "A Tool to print out all devices and their IOMMU groups")

	// Configure arguments
	gpu := parser.Flag("g", "gpu", &argparse.Options{
		Required: false,
		Help: "List all GPUs and devices related to them.",
	})

	usb := parser.Flag("u", "usb", &argparse.Options{
		Required: false,
		Help: "List all USB controllers.",
	})

	iommu_group := parser.IntList("i","group", &argparse.Options{
		Required: false,
		Help: "List everything in the IOMMU groups given. Supply argument multiple times to list additional groups.",
	})

	kernelmodules := parser.Flag("k","kernel", &argparse.Options{
		Required: false,
		Help: "Lists kernel modules using the devices and subsystems. (ignored if other options are present)",
	})

	test := parser.Flag("t", "test", &argparse.Options{
		Required: false,
		Help: "function im actively testing, does not do anything you care about (might be broken)",
	})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	// Work with the output depending on arguments given
	if *gpu {
		// List all GPUs (3d controllers are ignored)
		gpus := iommu.MatchSubclass(`VGA`)
		printoutput(gpus)
		printIOMMUgroup(*iommu_group)
	} else if *usb {
		// List all USB controllers
		usbs := iommu.MatchSubclass(`USB controller`)
		printoutput(usbs)
		printIOMMUgroup(*iommu_group)
	} else if *test {
		newTest(false, `USB controller`)
	} else if len(*iommu_group) > 0 {
		printIOMMUgroup(*iommu_group)
	}  else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(*kernelmodules)
		printoutput(out)
	}
}

// Function to just print out a string array to STDOUT
func printoutput(out []string)  {
	for _, line := range out {
		fmt.Print(line)
	}
}

// Function to print everything inside a specific IOMMU group
func printIOMMUgroup(groups []int) {
	// As long as we are asked to get devices from any specific IOMMU groups
	if len(groups) > 0 {
		// Get all IOMMU devices
		alldevs := iommu.NewIOMMU()
		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			// For each device in specified IOMMU group
			for _, device := range alldevs.Groups[group].Devices {
				// Generate output line
				line := iommu.GenDeviceLine(group, device)

				// Print the device info
				fmt.Print(line)
			}
		}
	}
	os.Exit(0)
}

func newTest(kernelmodules bool, searchval string) []string{
	var devs []string

	// Get all IOMMU devices
	alldevs := iommu.NewIOMMU()

	// Iterate through the groups
	for _, group := range alldevs.Groups {
		// For each device
		for _, device := range group.Devices {
			// If the device has a subclass matching what we are looking for
			if strings.Contains(device.Subclass.Name,searchval) {
				// Generate the device line
				line := iommu.GenDeviceLine(group.ID, device)
				// Append device line
				devs = append(devs, line)
			}
		}
	}

	return devs
}
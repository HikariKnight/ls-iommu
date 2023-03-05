package main

import (
	"fmt"
	"os"

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

	iommu_group := parser.StringList("i","group", &argparse.Options{
		Required: false,
		Help: "List everything in the IOMMU groups given. Supply argument multiple times to list additional groups.",
	})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	// Get all the IOMMU groups and their devices as a string array
	groups := iommu.GetIOMMU_Groups()

	// Work with the output depending on arguments given
	if *gpu {
		// List all GPUs (3d controllers are ignored)
		gpus := iommu.MatchDEVs(groups, `VGA`)
		printoutput(gpus)
		printIOMMUgroup(*iommu_group)
	} else if *usb {
		// List all USB controllers
		usbs := iommu.MatchDEVs(groups,`USB controller`)
		printoutput(usbs)
		printIOMMUgroup(*iommu_group)
	} else if len(*iommu_group) > 0 {
		printIOMMUgroup(*iommu_group)
	}  else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(groups)
		printoutput(out)
	}
}

// Print all devices in IOMMU group
func printIOMMUgroup(groups []string) {
	if len(groups) > 0 {
		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			devs := iommu.MatchDEVs(groups, `Group ` + group)
			printoutput(devs)
		}
	}
	os.Exit(0)
}

// Function to just print out a string array to STDOUT
func printoutput(out []string)  {
	for _, line := range out {
		fmt.Print(line)
	}
}
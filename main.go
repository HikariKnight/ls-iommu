package main

import (
	"fmt"
	"os"
	"strconv"

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
		gpus := iommu.MatchDEVs(false, `VGA`)
		printoutput(gpus)
		printIOMMUgroup(*iommu_group)
	} else if *usb {
		// List all USB controllers
		usbs := iommu.MatchDEVs(false,`USB controller`)
		printoutput(usbs)
		printIOMMUgroup(*iommu_group)
	} else if *test {
		newTest(*iommu_group)
	} else if len(*iommu_group) > 0 {
		printIOMMUgroup(*iommu_group)
	}  else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(*kernelmodules)
		printoutput(out)
	}
}

// Print all devices in IOMMU group
func printIOMMUgroup(groups []int) {
	if len(groups) > 0 {
		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			devs := iommu.MatchDEVs(false, `Group ` + strconv.Itoa(group))
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

func newTest(groups []int) {
	if len(groups) > 0 {
		devs := iommu.NewIOMMU()
		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			//devs := iommu.MatchDEVs(false, `Group ` + strconv.Itoa(group))
			fmt.Println(devs.Groups[group])
			//printoutput(devs)
		}
	}
	os.Exit(0)
}
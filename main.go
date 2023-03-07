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
		// Get all GPUs (3d controllers are ignored)
		output := iommu.MatchSubclass(`VGA`)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(*iommu_group)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if *usb {
		// Get all USB controllers
		output := iommu.MatchSubclass(`USB controller`)

		// Get all devices in specified IOMMU groups and append it to the output
		other := iommu.GetDevicesFromGroups(*iommu_group)
		output = append(output, other...)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	} else if *test {
		newTest(false, `USB controller`)
	} else if len(*iommu_group) > 0 {
		// Get all devices in specified IOMMU groups and append it to the output
		output := iommu.GetDevicesFromGroups(*iommu_group)

		// Print the output and exit
		printoutput(output)
		os.Exit(0)
	}  else {
		// Default behaviour mimicks the bash variant that this is based on
		out := iommu.GetAllDevices(*kernelmodules)
		printoutput(out)
	}
}

// Function to just print out a string array to STDOUT
func printoutput(out []string)  {
	output := removeDuplicateLines(out)
	for _, line := range output {
		fmt.Print(line)
	}
}

// Removes duplicate lines from a string slice, useful for cleaning up the output if doing multiple scans
func removeDuplicateLines(intSlice []string) []string {
    keys := make(map[string]bool)
    list := []string{}	
    for _, entry := range intSlice {
        if _, value := keys[entry]; !value {
            keys[entry] = true
            list = append(list, entry)
        }
    }    
    return list
}

func newTest(kernelmodules bool, searchval string) []string{
	devs := []string{"test", "test", "not test", "tast", "3", "3"}

	return devs
}
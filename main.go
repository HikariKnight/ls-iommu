package main

import (
	"fmt"
	"os"
	"regexp"

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
		newTest(false, `VGA`)
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
		devs := iommu.NewIOMMU()
		// For each IOMMU group given we will print the devices in each group
		for _, group := range groups {
			// For each device in specified IOMMU group
			for _, device := range devs.Groups[group].Devices {
				var line string

				// If the device has no revision, ommit the (rev ID), in both cases we generate the line with device info
				if device.Revision != "0x00" {
					line = fmt.Sprintf("IOMMU Group %v: %s %s: %s %s [%s:%s] (rev %s)\n",
					group,
					device.Address,
					device.Subclass.Name,
					device.Vendor.Name,
					device.Product.Name,
					device.Vendor.ID,
					device.Product.ID,
					device.Revision[len(device.Revision)-2:],
					)
				} else {
					line = fmt.Sprintf("IOMMU Group %v: %s %s: %s %s [%s:%s]\n",
					group,
					device.Address,
					device.Subclass.Name,
					device.Vendor.Name,
					device.Product.Name,
					device.Vendor.ID,
					device.Product.ID,
					)
				}

				// Print the device info
				fmt.Print(line)
			}
		}
	}
	os.Exit(0)
}

func newTest(kernelmodules bool, regex string) []string{
	var devs []string

	output := iommu.GetAllDevices(kernelmodules)
	gpuReg, err := regexp.Compile(regex)
	iommu.ErrorCheck(err)

	for _, line := range output {
		if gpuReg.MatchString(line) {
			devs = append(devs, line)
		}
	}

	return devs
}
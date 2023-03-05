package main

import (
	"fmt"
	"os"

	iommu "github.com/HikariKnight/ls-iommu/lib"
	"github.com/akamensky/argparse"
)
func main() {
	parser := argparse.NewParser("ls-iommu", "A Tool to print out all devices and their IOMMU groups")

	//var pf = fmt.Printf

	// Placeholder arg atm
	gpu := parser.Flag("g", "gpu", &argparse.Options{
		Required: false,
		Help: "List all GPUs and devices related to them",
	})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	groups := iommu.GetIOMMU_Groups()

	if *gpu {
		iommu.GetGPUs(groups)
	} else {
		out := iommu.GetAllDevices(groups)
		for _, line := range out {
			fmt.Print(line)
		}

	}
}
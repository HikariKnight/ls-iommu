package main

import (
	"fmt"
	"os"

	iommu "github.com/HikariKnight/ls-iommu/lib"
	"github.com/akamensky/argparse"
)
func main() {
	parser := argparse.NewParser("ls-iommu", "A Tool to print out all devices and their IOMMU groups")
	gpu := parser.Flag("g", "gpu", &argparse.Options{Required: false, Help: "List all GPUs and devices related to them"})

	var pf = fmt.Printf
	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
	}

	if *gpu {
		pf("GPU arg: %v\n", *gpu)
	}

	iommu.GetIOMMU_Groups()

}
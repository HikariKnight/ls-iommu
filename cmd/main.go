package main

import (
	"os"

	iommu "github.com/HikariKnight/ls-iommu/pkg/iommu"
	params "github.com/HikariKnight/ls-iommu/pkg/params"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	// Work with the output depending on arguments given
	if pArg.Flag["gpu"] {
		// Get all GPUs (3d controllers are ignored)
		output := iommu.MatchSubclass(`VGA`, pArg)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else if pArg.Flag["usb"] {
		// Get all USB controllers
		output := iommu.MatchSubclass(`USB controller`, pArg)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else if pArg.Flag["nic"] {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`Ethernet controller`, pArg)

		// Get all Wi-Fi controllers
		wifi := iommu.MatchSubclass(`Network controller`, pArg)
		output = append(output, wifi...)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else if pArg.Flag["sata"] {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`SATA controller`, pArg)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else if len(pArg.IntList["iommu_group"]) > 0 {
		// Get all devices in specified IOMMU groups and append it to the output
		output := iommu.GetDevicesFromGroups(pArg.IntList["iommu_group"], pArg.FlagCounter["related"], pArg)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else {
		// Default behaviour mimicks the bash variant that this is based on
		output := iommu.GetAllDevices(pArg)
		iommu.PrintOutput(output, pArg)
	}
}

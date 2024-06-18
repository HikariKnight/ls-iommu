package main

import (
	"fmt"
	"os"

	"github.com/HikariKnight/ls-iommu/internal/version"
	iommu "github.com/HikariKnight/ls-iommu/pkg/iommu"
	params "github.com/HikariKnight/ls-iommu/pkg/params"
)

func main() {
	// Get all our arguments in 1 neat struct
	pArg := params.NewParams()

	// Display version and exit if the version flag is present
	if pArg.Flag["version"] {
		fmt.Printf("ls-iommu version %s built in Go\n", version.Version)
		os.Exit(0)
	}

	// Work with the output depending on arguments given
	if pArg.Flag["gpu"] {
		// Get all GPUs (3d controllers are ignored)
		output := iommu.MatchSubclass(`VGA`, pArg)

		// Get all 3D controllers
		controller3d := iommu.MatchSubclass(`3D`, pArg)
		output = append(output, controller3d...)

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

	} else if pArg.Flag["nvme"] {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`Non-Volatile memory controller`, pArg)

		// Print the output and exit
		iommu.PrintOutput(output, pArg)
		os.Exit(0)

	} else if pArg.Flag["audio"] {
		// Get all Ethernet controllers
		output := iommu.MatchSubclass(`Audio device`, pArg)

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

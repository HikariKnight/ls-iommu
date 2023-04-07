package iommu

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/HikariKnight/ls-iommu/pkg/errorcheck"
	"github.com/HikariKnight/ls-iommu/pkg/params"
	"github.com/jaypipes/ghw/pkg/pci"
)

// Function to get the vbios path for a device
func GetRomPath(device *pci.Device, pArg *params.Params) []string {
	// Make a string slice to contain our paths
	var roms []string

	// Walk through /sys/devices/ and add all paths that has a rom file and matches the device.Address to the roms variable
	err := filepath.Walk("/sys/devices/", func(path string, info fs.FileInfo, err error) error {
		errorcheck.ErrorCheck(err, "Unable to walk /sys/devices/")

		// If the file name is "rom" and the path contains the PCI address for our device
		if info.Name() == "rom" && strings.Contains(path, device.Address) {
			// Add the filepath to our roms variable
			roms = append(roms, fmt.Sprintf("%s\n", path))
		}
		return nil
	})
	errorcheck.ErrorCheck(err, "Failed to walk /sys/devices/ to find the device rom")

	// Return all found rom files
	return roms
}

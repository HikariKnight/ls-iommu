package iommu

import (
	"fmt"
	"regexp"

	"github.com/jaypipes/ghw/pkg/pci"
)

func GenDeviceLine(group int, device *pci.Device, legacyoutput ...bool) string {
	var line string
	pciaddrclean := regexp.MustCompile(`^\d+:`)

	// If we want legacy output (to be output compatible with the bash version)
	var iommu_group string
	if len(legacyoutput) > 0 {
		// Do not pad the group number
		iommu_group = fmt.Sprintf("%d", group)
	} else {
		// Else we pad the group number to make it sortable
		iommu_group = fmt.Sprintf("% 3d", group)
	}

	// If the device has no revision, ommit the (rev ID), in both cases we generate the line with device info
	if device.Revision != "0x00" {
		line = fmt.Sprintf("IOMMU Group %s: %s %s [%s%s]: %s %s [%s:%s] (rev %s)\n",
			iommu_group,
			pciaddrclean.ReplaceAllString(device.Address, ""),
			device.Subclass.Name,
			device.Class.ID,
			device.Subclass.ID,
			device.Vendor.Name,
			device.Product.Name,
			device.Vendor.ID,
			device.Product.ID,
			device.Revision[len(device.Revision)-2:],
		)
	} else {
		line = fmt.Sprintf("IOMMU Group %s: %s %s [%s%s]: %s %s [%s:%s]\n",
			iommu_group,
			pciaddrclean.ReplaceAllString(device.Address, ""),
			device.Subclass.Name,
			device.Class.ID,
			device.Subclass.ID,
			device.Vendor.Name,
			device.Product.Name,
			device.Vendor.ID,
			device.Product.ID,
		)
	}

	return line
}

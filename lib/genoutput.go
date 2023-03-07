package iommu

import (
	"fmt"
	"regexp"

	"github.com/jaypipes/ghw/pkg/pci"
)

func GenDeviceLine(group int, device *pci.Device) string {
	var line string
	pciaddrclean := regexp.MustCompile(`^\d+:`)
	cleanoutput := regexp.MustCompile(`^IOMMU Group\s{1}\d+:`)

	// If the device has no revision, ommit the (rev ID), in both cases we generate the line with device info
	if device.Revision != "0x00" {
		line = fmt.Sprintf("IOMMU Group % 3d: %s %s [%s%s]: %s %s [%s:%s] (rev %s)\n",
		group,
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
		line = fmt.Sprintf("IOMMU Group % 3d: %s %s [%s%s]: %s %s [%s:%s]\n",
		group,
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
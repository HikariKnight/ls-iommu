package iommu

import (
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/jaypipes/ghw"
	ghwpci "github.com/jaypipes/ghw/pkg/pci"
)

type IOMMU struct {
	Groups map[int]*Group
}

func (i *IOMMU) AddGroup(group *Group) {
	i.Groups[group.ID] = group
}

type Group struct {
	ID      int
	Devices map[string]*ghwpci.Device
}

func (g *Group) AddDevice(device *ghwpci.Device) {
	g.Devices[(*device).Address] = device
}

func NewGroup(id int, devices map[string]*ghwpci.Device) *Group {
	return &Group{
		ID:      id,
		Devices: devices,
	}
}

func (i *IOMMU) Read() {
	i.Groups = make(map[int]*Group)
	// Get all groups and associated devices
	iommu_devices, err := filepath.Glob("/sys/kernel/iommu_groups/*/devices/*")
	if err != nil {
		log.Fatalln(err)
	}
	pci, err := ghw.PCI(ghw.WithDisableWarnings())
	if err != nil {
		log.Fatalln(err)
	}

	iommu_regex := regexp.MustCompile(`/sys/kernel/iommu_groups/(.*)/devices/(.*)`)
	for _, iommu_device := range iommu_devices {
		matches := iommu_regex.FindStringSubmatch(iommu_device)
		group_id, err := strconv.Atoi(matches[1])
		if err != nil {
			// Failed to properly parse groupid into integer, invalid for a group id, skip it
			continue
		}
		device_id := matches[2]

		if strings.HasPrefix(device_id, "0000:") {
			device := pci.GetDevice(device_id)

			// If the group doesn't exist in the struct, add it
			_, exists := i.Groups[group_id]
			if !exists {
				grp := &Group{
					ID:      group_id,
					Devices: make(map[string]*ghwpci.Device),
				}
				grp.AddDevice(device)
				i.AddGroup(grp)
			} else {
				i.Groups[group_id].AddDevice(device)
			}
		}
	}
}

func NewIOMMU() *IOMMU {
	iommu := &IOMMU{}
	iommu.Read()
	return iommu
}

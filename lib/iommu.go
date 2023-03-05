package iommu

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func ListDirs(path string) []string{
	// Make a string array to hold all dir names
	var dirs []string
	files, err := os.ReadDir(path)
	ErrorCheck(err)

	for _, file := range files {
		dirs = append(dirs, file.Name())
	}

	return dirs
}

// Function to get all IOMMU Groups on the machine
func GetIOMMU_Groups() []string{
	groups := ListDirs("/sys/kernel/iommu_groups")

	return groups
}

func GetIOMMU_Devices(groups []string) []string{
	var devices []string

	for _, group := range groups {
		devices = ListDirs("/sys/kernel/iommu_groups/" + group + "/devices/")
	}

	return devices
}

func GetAllDevices(groups []string) []string {
	var lspci_devs []string

	for _, group := range groups {
		devices := GetIOMMU_Devices([]string{group})
		for _, device := range devices {
			cmd := exec.Command("lspci", "-nns", device)
			
			var out bytes.Buffer
			cmd.Stdout = &out

			err := cmd.Run()
			ErrorCheck(err)

			lspci_devs = append(lspci_devs, fmt.Sprintf("IOMMU Group %s: %s", group, out.String()))
		}
	}

	return lspci_devs
}

func GetGPUs(groups []string) {
	output := GetAllDevices(groups)

	for _, line := range output {
		fmt.Print(line)
	}
}
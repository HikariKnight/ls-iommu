package iommu

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func listDirs(path string) []string{
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
	groups := listDirs("/sys/kernel/iommu_groups")

	return groups
}

func getIOMMU_Devices(groups []string) []string{
	var devices []string

	for _, group := range groups {
		devices = listDirs(fmt.Sprintf("/sys/kernel/iommu_groups/%s/devices/", group))
	}

	return devices
}

func GetAllDevices(groups []string) []string {
	var lspci_devs []string

	for _, group := range groups {
		devices := getIOMMU_Devices([]string{group})
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

func MatchDEVs(groups []string, regex string) []string{
	var devs []string

	output := GetAllDevices(groups)
	gpuReg, err := regexp.Compile(regex)
	ErrorCheck(err)

	for _, line := range output {
		if gpuReg.MatchString(line) {
			devs = append(devs, line)
		}
	}

	return devs
}
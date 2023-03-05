package iommu

import (
	"fmt"
	"os"
)

// Function to get all IOMMU Groups on the machine
func GetIOMMU_Groups() {
	// Make an array of strings for us to store the IOMMU groups
	var groups []string

	files, err := os.ReadDir("/sys/kernel/iommu_groups")
	ErrorCheck(err)

	for _, group := range files {
		if group.IsDir() {
			groups = append(groups, group.Name())
		}
	}

	fmt.Println(groups)
} 
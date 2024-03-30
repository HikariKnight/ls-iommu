# ls-iommu
A tool to list devices in iommu groups, useful for setting up VFIO

This tool is a go implementation and extension to the small bash script provided by Wendell from Level1Techs.
It's purpose was to just list every device and their associated IOMMU group.

This project does the same thing, but extends the functionality by implementing arguments to list the output in different ways or provide just the details needed without having to hop through grep, sed, awk or perl pipe hoops.

Currently the program supports the same behavior as Wendell's ls-iommu script but can also be told to display only devices in selected IOMMU groups, only GPUs or only USB controllers and also locate related devices.<br>
More functionality can be added if it is deemed useful, just open an issue with the request.

Note: arm builds are generated but not tested as I lack the relevant hardware.

![helptext](https://github.com/HikariKnight/ls-iommu/assets/2557889/4c04e171-5b21-4858-8810-76daa7d15303)
![gpu output](https://github.com/HikariKnight/ls-iommu/assets/2557889/a6ef282b-dc25-493b-8a7c-3f5cd8fcff3f)
![getting info from specific devices](https://github.com/HikariKnight/ls-iommu/assets/2557889/19bf09d6-76e6-47c7-8875-3aa73c327e15)
![display specific info](https://github.com/HikariKnight/ls-iommu/assets/2557889/c0453d50-db09-4d41-8701-59477a567654)


## Features
* Has a flag to generate default ouptut text compatible with the bash version of ls-iommu (`--legacy`)
* Can locate and display only GPUs, USB Controllers and Network cards
* List devices of individual or multiple IOMMU groups
* Can show you kernel driver info for the devices
* Can attempt to get related devices (devices that share vendorID) or devices located in the same IOMMU group as the device(s), the best method is used for -r depending on what you look for.
* Flag to list only VendorID:DeviceID for all devices found
* Flag to list only PCI Address for all devices found
* A flag you can use to ignore specific VendorIDs when doing a related devices search
* Can show only the info you are interested in on the device line using -F and then providing a comma separated list of objects you want to show on the device line (this does not affect the extra lines provided by `-k`)
* Has a flag to get the vbios path for gpus (`-g --rom` or `-g -i X --rom`)


## Build instructions
Prerequisites: 
* Go 1.20+
* git

This will build the latest `ls-iommu` and set the version to the latest commit hash.
```bash
git clone https://github.com/HikariKnight/ls-iommu.git
cd ls-iommu
CGO_ENABLED=0 go build -ldflags="-X github.com/HikariKnight/ls-iommu/internal/version.Version=$(git rev-parse --short HEAD)" -o ls-iommu cmd/main.go
```

The binary `ls-iommu` will now be located the root of the project directory.

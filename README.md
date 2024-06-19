# ls-iommu
A tool to list devices in iommu groups, useful for setting up VFIO and is a dependency for [quickpassthrough](https://github.com/HikariKnight/quickpassthrough)

This tool is a go implementation and extension to the small bash script provided by Wendell from Level1Techs.
It's purpose was to just list every device and their associated IOMMU group.

This project does the same thing, but extends the functionality by implementing arguments to list the output in different ways or provide just the details needed without having to hop through grep, sed, awk or perl pipe hoops.

You can [download the latest release](https://github.com/HikariKnight/ls-iommu/releases/) from the releases page

## Features
Currently the program supports the same behavior as Wendell's ls-iommu script but can also be told to display:
- Output sorted properly by IOMMU group
- Only devices in selected IOMMU groups
- Only specific devices (ex: GPUs, USB controllers, SATA controllers, etc)
- Tailor the output to show only what you care about
- Locate relative devices sharing the same IOMMU group
- Display currently used kernel driver for listed devices
- Display only device IDs for queried devices
- Display only PCI addresses queried devices
- Display rom path for GPUs (or the selected GPU using `-i` to only show devices in a specific IOMMU group)

More functionality can be added if it is deemed useful, just open an issue with the request.
<br>
**Note:** arm64 builds are generated but not tested as I lack the relevant hardware.

![helptext](https://github.com/HikariKnight/ls-iommu/assets/2557889/4c04e171-5b21-4858-8810-76daa7d15303)
![gpu output](https://github.com/HikariKnight/ls-iommu/assets/2557889/a6ef282b-dc25-493b-8a7c-3f5cd8fcff3f)
![getting info from specific devices](https://github.com/HikariKnight/ls-iommu/assets/2557889/19bf09d6-76e6-47c7-8875-3aa73c327e15)
![display specific info](https://github.com/HikariKnight/ls-iommu/assets/2557889/c0453d50-db09-4d41-8701-59477a567654)


## Build instructions
Prerequisites: 
* Go 1.20+
* git

This will build the latest `ls-iommu` and set the version to the latest commit hash.
```bash
git clone https://github.com/HikariKnight/ls-iommu.git
cd ls-iommu
go mod download
CGO_ENABLED=0 go build -ldflags="-X github.com/HikariKnight/ls-iommu/internal/version.Version=$(git rev-parse --short HEAD)" -o ls-iommu cmd/main.go
```
NOTE: you can build with newer dependencies (can break things) by running `go get -u ./cmd` after `go mod download`

The binary `ls-iommu` will now be located the root of the project directory.

# ls-iommu
A tool to list devices in iommu groups, useful for setting up VFIO

This tool is a go implementation and extension to the small bash script provided by Wendell from Level1Techs.
It's purpose was to just list every device and their associated IOMMU group.

This project does the same thing, but extends the functionality by implementing arguments to list the output in different ways or provide just the details needed without having to hop through grep, sed, awk or perl pipe hoops.

Currently the program supports the same behavior as Wendell's ls-iommu script but can also be told to display only devices in selected IOMMU groups, only GPUs or only USB controllers.<br>
More extended functionality is planned.
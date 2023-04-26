package params

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

/*
	The whole purpose of this module is to make a struct
	to just carry all our parsed arguments around between functions

	Create a Params struct with all the argparse arguments
	pArg := params.NewParams()
*/

type Params struct {
	Flag        map[string]bool
	FlagCounter map[string]int
	IntList     map[string][]int
	StringList  map[string][]string
	String      map[string]string
}

func (p *Params) addFlag(name string, flag bool) {
	p.Flag[name] = flag
}

func (p *Params) addFlagCounter(name string, flag int) {
	p.FlagCounter[name] = flag
}

func (p *Params) addIntList(name string, flag []int) {
	p.IntList[name] = flag
}

func (p *Params) addStringList(name string, flag []string) {
	p.StringList[name] = flag
}

func (p *Params) addString(name string, flag string) {
	p.String[name] = flag
}

func NewParams() *Params {
	// Setup the parser for arguments
	parser := argparse.NewParser("ls-iommu", "A Tool to print out all devices and their IOMMU groups")

	// Configure arguments
	gpu := parser.Flag("g", "gpu", &argparse.Options{
		Required: false,
		Help:     "List all GPUs. (use -i # to only display results from specified IOMMU group)",
	})

	usb := parser.Flag("u", "usb", &argparse.Options{
		Required: false,
		Help:     "List all USB controllers. (use -i # to only display results from specified IOMMU group)",
	})

	nic := parser.Flag("n", "network", &argparse.Options{
		Required: false,
		Help:     "List all Network controllers. (use -i # to only display results from specified IOMMU group)",
	})

	sata := parser.Flag("s", "sata", &argparse.Options{
		Required: false,
		Help:     "List all SATA controllers. (use -i # to only display results from specified IOMMU group)",
	})

	iommu_group := parser.IntList("i", "group", &argparse.Options{
		Required: false,
		Help:     "List everything in the IOMMU groups given. Supply argument multiple times to list additional groups.",
	})

	related := parser.FlagCounter("r", "related", &argparse.Options{
		Required: false,
		Help:     "Attempt to list related devices that share IOMMU Groups or\n\t\t Vendor ID (used with -g -u -i -s and -n), pass -rr if you want to search using both when used with -g -i -s or -n\n\t\t Note: -rr can be inaccurate or too broad when many devices share Vendor ID",
	})

	ignore := parser.StringList("R", "ignore", &argparse.Options{
		Required: false,
		Help:     "Ignores passed VendorID (Left part of : in [VendorID:DeviceID]) outside of the selected IOMMU group when doing a --related search, you can use this to ignore unreliable Vendor IDs when doing related searches. (works with -g -i -u -s and -n)",
	})

	kernelmodules := parser.Flag("k", "kernel", &argparse.Options{
		Required: false,
		Help:     "Lists subsystems and kernel drivers using the devices.",
		Default:  false,
	})

	legacyoutput := parser.Flag("", "legacy", &argparse.Options{
		Required: false,
		Help:     "Generate the output unsorted and be the same output as the old bash script",
		Default:  false,
	})

	id := parser.Flag("", "id", &argparse.Options{
		Required: false,
		Help:     "Print out only VendorID:DeviceID for non bridge devices (Only works with -i)",
	})

	pciaddr := parser.Flag("", "pciaddr", &argparse.Options{
		Required: false,
		Help:     "Print out only the PCI Address for non bridge devices (Only works with -i)",
	})

	rom := parser.Flag("", "rom", &argparse.Options{
		Required: false,
		Help:     "Print out the rom path GPUs. (must be used with -g or --gpu)",
	})

	format := parser.String("F", "format", &argparse.Options{
		Required: false,
		Help:     "Formats the device line output the way you want it (omit what you do not want)\n\t\t Supported objects: pciaddr, subclass_name, subclass_name:, subclass_id, subclass_id:, name, name:, device_id, device_id:, vendor, vendor:, oem, oem:, prod_name, prod_name:, revision, optional_revision",
		Default:  "pciaddr,subclass_name,subclass_id,name,device_id,optional_revision",
	})

	// Parse arguments
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(4)
	}

	// Make our struct
	pArg := &Params{
		Flag:        make(map[string]bool),
		FlagCounter: make(map[string]int),
		IntList:     make(map[string][]int),
		StringList:  make(map[string][]string),
		String:      make(map[string]string),
	}

	// Add all parsed arguments to a struct for portability since we will use them all over the program
	pArg.addFlag("gpu", *gpu)
	pArg.addFlag("usb", *usb)
	pArg.addFlag("nic", *nic)
	pArg.addFlag("sata", *sata)
	pArg.addFlagCounter("related", *related)
	pArg.addStringList("ignore", *ignore)
	pArg.addIntList("iommu_group", *iommu_group)
	pArg.addFlag("kernelmodules", *kernelmodules)
	pArg.addFlag("legacyoutput", *legacyoutput)
	pArg.addFlag("id", *id)
	pArg.addFlag("pciaddr", *pciaddr)
	pArg.addFlag("rom", *rom)
	pArg.addString("format", *format)

	return pArg
}

package params

type Params struct {
	Flag        map[string]bool
	FlagCounter map[string]int
	IntList     map[string][]int
}

func (p *Params) AddFlag(name string, flag bool) {
	p.Flag[name] = flag
}

func (p *Params) AddFlagCounter(name string, flag int) {
	p.FlagCounter[name] = flag
}

func (p *Params) AddIntList(name string, flag []int) {
	p.IntList[name] = flag
}

/*
Create a Params struct with
	pArg := params.NewParams()
*/

func NewParams() *Params {
	p := &Params{
		Flag:        make(map[string]bool),
		FlagCounter: make(map[string]int),
		IntList:     make(map[string][]int),
	}
	return p
}

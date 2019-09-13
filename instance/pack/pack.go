package pack

type Options map[string]*Item

type Pack struct {
	Required []*Path  `json:"required"`
	Options  *Options `json:"optional"`
}

type Item struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Enabled      bool     `json:"enabled"`
	Dependencies []string `json:"dependencies"`
	Files        []*Path  `json:"files"`
}

type Path struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

func (p *Pack) Sync(o *Options) *Options {
	if o != nil {
		for k := range *p.Options {
			if v, ok := (*o)[k]; ok {
				(*p.Options)[k] = v
			}
		}
	}
	return p.Options
}

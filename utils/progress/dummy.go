package progress

var dum = &dummy{}

type dummy struct {
}

func Dummy() Listener {
	return dum
}

func (d *dummy) Stat(string, float64)   {}
func (d *dummy) GlobalStatus(string)    {}
func (d *dummy) GlobalProgress(float64) {}
func (d *dummy) TaskStatus(string)      {}
func (d *dummy) TaskProgress(float64)   {}
func (d *dummy) Wait()                  {}

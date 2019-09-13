package version

import (
	"bytes"
	"strconv"
	"strings"
)

const (
	NoUpgrade    = 0
	PatchUpgrade = 1
	MinorUpgrade = 2
	MajorUpgrade = 3
)

type Version []int

func Parse(s string) *Version {
	var version Version
	parseVersion(s, &version)
	return &version
}

func parseVersion(s string, v *Version) {
	parts := strings.Split(s, ".")
	for _, p := range parts {
		i, e := strconv.Atoi(p)
		if e == nil {
			*v = append(*v, i)
		}
	}
}

func (v *Version) Major() int {
	return v.pointAt(0)
}

func (v *Version) Minor() int {
	return v.pointAt(1)
}

func (v *Version) Patch() int {
	return v.pointAt(2)
}

func (v *Version) Latest(o *Version) *Version {
	i := v.Compare(o)
	if i < 0 {
		return o
	}
	return v
}

func (v *Version) UpgradeIndex(o *Version) int {
	if v.Major() != o.Major() {
		return MajorUpgrade
	}
	if v.Minor() != o.Minor() {
		return MinorUpgrade
	}
	if v.Patch() != o.Patch() {
		return PatchUpgrade
	}
	return NoUpgrade
}

func (v *Version) Matches(o *Version) bool {
	return v.Compare(o) == 0
}

func (v *Version) Compare(o *Version) int {
	max := max(len(*v), len(*o))
	for i := 0; i < max; i++ {
		iv := v.pointAt(i)
		io := o.pointAt(i)
		if iv != io {
			return iv - io
		}
	}
	return 0
}

func (v *Version) ToString(points int) string {
	wr := bytes.Buffer{}
	for i := 0; i < points; i++ {
		if i > 0 {
			wr.WriteString(".")
		}
		p := v.pointAt(i)
		wr.WriteString(strconv.Itoa(p))
	}
	return wr.String()
}

func (v *Version) String() string {
	return v.ToString(len(*v))
}

func (v *Version) pointAt(i int) int {
	if i < len(*v) {
		return (*v)[i]
	}
	return 0
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

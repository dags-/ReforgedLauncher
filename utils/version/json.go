package version

func (v *Version) MarshalJSON() ([]byte, error) {
	s := v.String()
	return []byte(`"` + s + `"`), nil
}

func (v *Version) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]
	parseVersion(s, v)
	return nil
}

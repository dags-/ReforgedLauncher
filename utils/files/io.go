package files

import "io/ioutil"

func ReadText(path string) (string, error) {
	d, e := ioutil.ReadFile(path)
	if e != nil {
		return "", e
	}
	return string(d), nil
}

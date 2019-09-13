package errs

import (
	"fmt"
)

func Panic(name string, e error) {
	if e != nil {
		fmt.Println("Crashed at " + name + ":")
		panic(e)
	}
}

func Log(name string, e error) {
	if e != nil {
		fmt.Println(name, e.Error())
	}
}

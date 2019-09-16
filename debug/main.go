package main

import (
	"fmt"
	"github.com/dags-/systray"
	"io/ioutil"
	"log"
)

func main() {
	systray.Run(func() {
		b, e := ioutil.ReadFile("assets/public/assets/image/tray.png")
		fmt.Println(e)

		systray.SetIcon(b)
		setup()
	}, func() {
		log.Println("close")
	})
}

func setup() {
	log.Println("setup")
	i := systray.AddMenuItem("PLS", "")
	<-i.ClickedCh
	fmt.Println("done")
	systray.Quit()
}

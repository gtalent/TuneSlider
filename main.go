package main

import (
	"fmt"
	"io/ioutil"
	"json"
)

type show struct {
	Slides []slide
}

type slide struct {
	Contents string
}

func main() {
	var show show
	show.Slides = make([]slide, 10)
	slides := show.Slides
	no := '0'
	for i := 0; i < 10; i++ {
		slides[i].Contents = "Slide " + string(no)
		no++
	}
	output, err := json.Marshal(&show)
	if err != nil {
		fmt.Println(err.String())
	}
	ioutil.WriteFile("song.sld", output, 0644)
}


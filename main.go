package main

import (
	"fmt"
	"io/ioutil"
	"json"
	gfx "wombat/core/graphics"
)

type show struct {
	Slides []slide
}

type slide struct {
	TextColor       gfx.Color
	BackgroundColor gfx.Color
	Contents        string
}

func main() {
	var show show
	show.Slides = make([]slide, 10)
	slides := show.Slides
	no := '0'
	for i := 0; i < 10; i++ {
		slides[i].Contents = "Slide " + string(no)
		slides[i].TextColor = gfx.Color{255, 255, 255}
		slides[i].BackgroundColor = gfx.Color{0, 0, 0}
		no++
	}
	output, err := json.Marshal(&show)
	if err != nil {
		fmt.Println(err.String())
	}
	ioutil.WriteFile("song.sld", output, 0644)
}

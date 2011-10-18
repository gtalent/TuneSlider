package main

import (
	"fmt"
	"io/ioutil"
	"json"
	"time"
	gfx "wombat/core/graphics"
	"wombat/core/input"
)

type show struct {
	Slides []slide
}

type slide struct {
	TextColor       gfx.Color
	BackgroundColor gfx.Color
	Duration        int64 //seconds
	Contents        string
}

type Drawer struct {
	screen *gfx.Display
	slide
	text  gfx.Text
	font  *gfx.Font
}

func (me *Drawer) init() {
	me.font.SetColor(me.TextColor)
	me.font.Write(me.Contents, &me.text)
}

func (me *Drawer) Draw(c *gfx.Canvas) {
	c.SetColor(me.BackgroundColor)
	c.FillRect(0, 0, me.screen.GetWidth(), me.screen.GetHeight())
	c.SetColor(me.TextColor)
	c.DrawText(&me.text, me.screen.GetWidth() / 2 - me.text.Width() / 2, me.screen.GetHeight() / 2 - me.text.Height() / 2)
}

func buildTestInput() {
	var show show
	show.Slides = make([]slide, 10)
	slides := show.Slides
	no := '0'
	for i := 0; i < 10; i++ {
		slides[i].Contents = "Slide " + string(no)
		slides[i].TextColor = gfx.Color{255, 255, 255}
		slides[i].BackgroundColor = gfx.Color{0, 0, 0}
		slides[i].Duration = 2
		no++
	}
	output, err := json.Marshal(&show)
	if err != nil {
		fmt.Println(err.String())
	}
	ioutil.WriteFile("song.sld", output, 0644)
}

func main() {
	buildTestInput()
	file, err := ioutil.ReadFile("song.sld")
	var show show
	json.Unmarshal(file, &show)
	if err != nil {
		fmt.Println(err.String())
	}
	running := true
	screen := gfx.NewDisplay()
	screen.Open(800, 600)
	input.Init()
	input.AddQuit(func() {
		running = false
		screen.Close()
	})
	slides := show.Slides
	var d Drawer
	d.screen = screen
	d.font = gfx.LoadFont("font.otf", 32)
	screen.AddDrawer(&d)
	for i := 0; running && i < len(slides); i++ {
		slide := slides[i]
		d.slide = slide
		d.init()
		time.Sleep(slide.Duration * 1000000000)
	}
}

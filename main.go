/*
   Copyright 2011 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"json"
	"time"
	"sdl/mixer"
	gfx "wombat/core/graphics"
	"wombat/core/input"
)

var kill chan interface{} = make(chan interface{})

type show struct {
	Slides    []slide
	AudioPath string
}

type slide struct {
	TextColor       gfx.Color
	BackgroundColor gfx.Color
	Duration        int64 //seconds
	Contents        string
}

type Drawer struct {
	slide
	initialized bool
	screen      *gfx.Display
	text        gfx.Text
	font        *gfx.Font
}

func (me *Drawer) init() {
	me.font.SetColor(me.TextColor)
	me.font.Write(me.Contents, &me.text)
	me.initialized = true
}

func (me *Drawer) Draw(c *gfx.Canvas) {
	c.SetColor(me.BackgroundColor)
	c.FillRect(0, 0, me.screen.GetWidth(), me.screen.GetHeight())
	c.SetColor(me.TextColor)
	if me.initialized {
		c.DrawText(&me.text, me.screen.GetWidth()/2-me.text.Width()/2, me.screen.GetHeight()/2-me.text.Height()/2)
	}
}

func buildTestInput() {
	var show show
	show.AudioPath = "song.ogg"
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

func playShow(show show, d *Drawer) {
	slides := show.Slides
	song := mixer.LoadMUS(show.AudioPath)
	song.PlayMusic(1)
	for i := 0; true; {
		slide := slides[i]
		d.slide = slide
		d.init()
		i++
		if i < len(slides) {
			time.Sleep(slide.Duration * 1000000000)
		} else {
			break
		}
	}
	mixer.FadeOutMusic(3000)
	song.Free()
}

func main() {
	buildTestInput()
	file, err := ioutil.ReadFile("song.sld")
	var show show
	json.Unmarshal(file, &show)
	if err != nil {
		fmt.Println(err.String())
	}
	screen := gfx.NewDisplay()
	screen.Open(800, 600)
	mixer.OpenAudio(mixer.DEFAULT_FREQUENCY, mixer.DEFAULT_FORMAT, mixer.DEFAULT_CHANNELS, 4096)
	input.Init()
	input.AddQuit(func() {
		kill <- nil
		screen.Close()
	})

	var d Drawer
	d.screen = screen
	d.font = gfx.LoadFont("font.otf", 32)
	screen.AddDrawer(&d)
	go playShow(show, &d)
	<-kill
}

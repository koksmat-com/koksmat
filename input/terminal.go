package input

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

func GetString(defaultValue string, caption string) string {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	fmt.Println(caption + ":")
	value := readInputWithDefault(defaultValue)
	return value

}

func readInputWithDefault(defaultValue string) string {
	input := []rune{}
	cursorPos := 0

	displayInput(input, cursorPos, defaultValue)

	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return defaultValue
			case termbox.KeyEnter:
				if len(input) == 0 {
					return defaultValue
				}
				return string(input)
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
					cursorPos--
				}
			default:
				if ev.Ch != 0 {
					input = append(input, ev.Ch)
					cursorPos++
				}
			}
			displayInput(input, cursorPos, defaultValue)
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}

func displayInput(input []rune, cursorPos int, defaultValue string) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	x := 0

	if len(input) == 0 {
		printWithDim(defaultValue, x)
	} else {
		printWithNormal(string(input), x)
	}
	termbox.SetCursor(cursorPos, 0)
	termbox.Flush()
}

func printWithDim(s string, x int) {
	for _, c := range s {
		termbox.SetCell(x, 0, c, termbox.ColorBlack, termbox.ColorWhite)
		x++
	}
}

func printWithNormal(s string, x int) {
	for _, c := range s {
		termbox.SetCell(x, 0, c, termbox.ColorDefault, termbox.ColorDefault)
		x++
	}
}

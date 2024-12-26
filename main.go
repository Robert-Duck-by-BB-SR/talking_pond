package main

import (
	"fmt"
	"os"

	dd "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/duck_dom"
	"golang.org/x/term"
)

func main() {
	old_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), old_state)

	dd.ClearScreen()

	screen := dd.Screen{State: &dd.Normal, EventLoopIsRunning: true}
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	screen.Width = width
	screen.Height = height

	sidebar := dd.CreateWindow(dd.Styles{
		Width:      50,
		Height:     screen.Height - 1,
		Background: dd.MakeRGBBackground(69, 150, 100),
		Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	})

	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts123123", dd.Styles{
			Width:      10,
			Height:     5,
			Background: dd.MakeRGBBackground(250, 0, 0),
			TextColor:  dd.MakeRGBTextColor(0, 0, 0),
			// Paddding:   1,
			Border: dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		}),
	)
	// FIX: I think that new block component after inline component should start from the beninging of a parent
	sidebar.AddComponent(
		dd.CreateComponent("Deez nuts", dd.Styles{
			Width:      10,
			Height:     10,
			Background: dd.MakeRGBBackground(250, 0, 0),
			TextColor:  dd.MakeRGBTextColor(0, 0, 0),
			Paddding:   1,
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
		),
	)
	screen.AddWindow(sidebar)

	content := dd.CreateWindow(dd.Styles{
		Width:      screen.Width - sidebar.Styles.Width - 1,
		Height:     int(float32(screen.Height)*0.7) + 1,
		Background: dd.MakeRGBBackground(69, 150, 100),
		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	})

	content.Position.StartingCol = sidebar.StartingCol + sidebar.Width
	content.Position.StartingRow = sidebar.StartingRow
	content.AddComponent(
		dd.CreateComponent(
			"|SIMD|",
			dd.Styles{
				Width:      10,
				Height:     10,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))
	content.AddComponent(
		dd.CreateComponent(
			"LIGMA???",
			dd.Styles{
				Width:      10,
				Height:     10,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		))

	screen.AddWindow(content)

	input_bar := &dd.Window{
		Position: dd.Position{StartingRow: content.Height + 1, StartingCol: sidebar.Width + 1},
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     int(float32(screen.Height)*0.3) - 1,
			Background: dd.MakeRGBBackground(150, 150, 40),
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	screen.AddWindow(input_bar)

	screen.StatusBar = dd.Window{
		Position: dd.Position{StartingRow: screen.Height, StartingCol: 1},
		Styles: dd.Styles{
			Width:      screen.Width,
			Height:     1,
			Background: dd.MakeRGBBackground(80, 40, 100),
		},
	}
	screen.StatusBar.Components = []*dd.Component{
		{
			Position: dd.Position{StartingRow: screen.Height, StartingCol: 2},
			Buffer:   dd.NORMAL,
			Styles: dd.Styles{
				Width:  len(dd.COMMAND),
				Height: 1,
			},
		},
	}

	screen.Render()
	screen.Activate()

	// TODO: Check if its possible to accept more than one byte
	stdin_buffer := make([]byte, 1)
	for screen.EventLoopIsRunning {
		for len(screen.RenderQueue) > 0 {
			item_to_render := screen.RenderQueue[0]
			fmt.Print(item_to_render)
			screen.RenderQueue = screen.RenderQueue[1:]
		}

		fmt.Printf(dd.MOVE_CURSOR_TO_POSITION, screen.CursorPosition.StartingRow, screen.CursorPosition.StartingCol)

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		screen.State.HandleKeypress(&screen, stdin_buffer)
	}
	// restart to default settings
	fmt.Print(dd.SHOW_CURSOR)
	// TODO: any assert should have show cursor
}

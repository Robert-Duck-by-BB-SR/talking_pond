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
		dd.CreateComponent("|Deez nuts|", dd.Styles{
				Width: len("|Deez nuts|") + 10,
				Height: 1,
				Background: dd.MakeRGBBackground(80, 40, 100),
				Border: dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
			},
		),
	)

	// sidebar.AddComponent(
	// 	&dd.Component{
	// 		// Position: dd.Position{StartingRow: 5, StartingCol: uint(sidebar.StartingCol) + 2},
	// 		Buffer: "|got em|",
	// 		Styles: dd.Styles{
	// 			Width: len("|got em|"),
	// 			Direction: dd.Inline,
	// 			Height:     1,
	// 			Background: dd.MakeRGBBackground(80, 40, 100),
	// 			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 		},
	// 	})

	screen.AddWindow(sidebar)

	// content := dd.Window{
	// 	Position: dd.Position{StartingRow: 1, StartingCol: uint(sidebar.Styles.Width) + 1},
	// 	Styles: dd.Styles{
	// 		Width:      screen.Width - sidebar.Styles.Width - 1,
	// 		Height:     int(float32(screen.Height)*0.7) + 1,
	// 		Background: dd.MakeRGBBackground(69, 150, 100),
	// 		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 	},
	// }
	// content.AddComponent(
	// 	&dd.Component{
	// 		Position: dd.Position{StartingRow: 2, StartingCol: uint(content.StartingCol) + 2},
	// 		Buffer:   "|SIMD|",
	// 		Styles: dd.Styles{
	// 			Width:      screen.Width - sidebar.Styles.Width - 1,
	// 			Height:     screen.Height,
	// 			Background: dd.MakeRGBBackground(80, 40, 100),
	// 			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 		},
	// 	})
	//
	// content.AddComponent(
	// 	&dd.Component{
	// 		Position: dd.Position{StartingRow: 4, StartingCol: uint(content.StartingCol) + 2},
	// 		Buffer:   "|Ligma?|",
	// 		Styles: dd.Styles{
	// 			Width:      screen.Width - sidebar.Styles.Width - 1,
	// 			Height:     screen.Height,
	// 			Background: dd.MakeRGBBackground(80, 40, 100),
	// 			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 		},
	// 	})

	// screen.AddWindow(&content)
	//
	// input_bar := dd.Window{
	// 	Position: dd.Position{StartingRow: uint(content.Height) + 1, StartingCol: uint(sidebar.Width) + 1},
	// 	Styles: dd.Styles{
	// 		Width:      screen.Width - sidebar.Styles.Width - 1,
	// 		Height:     int(float32(screen.Height)*0.3) - 1,
	// 		Background: dd.MakeRGBBackground(150, 150, 40),
	// 		Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
	// 	},
	// }

	// screen.AddWindow(&input_bar)
	//
	// screen.StatusBar = dd.Window{
	// 	Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 1},
	// 	Styles: dd.Styles{
	// 		Width:      screen.Width,
	// 		Height:     1,
	// 		Background: dd.MakeRGBBackground(80, 40, 100),
	// 	},
	// }
	// screen.StatusBar.AddComponent(
	// 	&dd.Component{
	// 		Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 2},
	// 		Buffer:   dd.NORMAL,
	// 		Styles: dd.Styles{
	// 			Width:  screen.Width,
	// 			Height: 1,
	// 		},
	// 	})

	screen.Activate()
	screen.Render()

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
}

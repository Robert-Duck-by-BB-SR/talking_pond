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

	sidebar := dd.Window{
		Position: dd.Position{StartingRow: 1, StartingCol: 1},
		Styles: dd.Styles{
			Width:      50,
			Height:     screen.Height - 1,
			Background: dd.MakeRGBBackground(69, 150, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	content := dd.Window{
		Position: dd.Position{StartingRow: 1, StartingCol: uint(sidebar.Styles.Width) + 2},
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height - 1,
			Background: dd.MakeRGBBackground(69, 150, 100),
			// Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	status_bar_component := dd.Component{
		Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 2},
		Buffer:   "NORMAL",
		Styles: dd.Styles{
			Width:  screen.Width,
			Height: 30,
		},
	}

	status_bar := dd.Window{
		Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 1},
		Styles: dd.Styles{
			Width:      screen.Width,
			Height:     1,
			Background: dd.MakeRGBBackground(80, 40, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
		Components: []dd.Component{status_bar_component},
	}

	item := dd.Component{
		Position: dd.Position{StartingRow: 3, StartingCol: uint(sidebar.StartingCol) + 2},
		Buffer:   "|Deez nuts|",
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height,
			Background: dd.MakeRGBBackground(69, 150, 100),
			Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	item_two := dd.Component{
		Position: dd.Position{StartingRow: 5, StartingCol: uint(sidebar.StartingCol) + 2},
		Buffer:   "|got em|",
	}

	sidebar.Components = []dd.Component{item, item_two}

	item_three := dd.Component{
		Position: dd.Position{StartingRow: 2, StartingCol: uint(content.StartingCol) + 2},
		Buffer:   "|SIMD|",
	}

	item_four := dd.Component{
		Position: dd.Position{StartingRow: 4, StartingCol: uint(content.StartingCol) + 2},
		Buffer:   "|Ligma?|",
	}

	content.Components = []dd.Component{item_three, item_four}

	screen.RenderQueue = append(screen.RenderQueue, sidebar.Render())
	screen.RenderQueue = append(screen.RenderQueue, content.Render())
	screen.RenderQueue = append(screen.RenderQueue, status_bar.Render())

	// TODO: render them together with parents
	for _, comp := range sidebar.Components {
		comp.Render()
		screen.RenderQueue = append(screen.RenderQueue, comp.Content)
	}
	for _, comp := range content.Components {
		comp.Render()
		screen.RenderQueue = append(screen.RenderQueue, comp.Content)
	}

	for _, comp := range status_bar.Components {
		comp.Render()
		screen.RenderQueue = append(screen.RenderQueue, comp.Content)
	}

	stdin_buffer := make([]byte, 1)
	for screen.EventLoopIsRunning {
		for len(screen.RenderQueue) > 0 {
			item_to_render := screen.RenderQueue[0]
			fmt.Print(item_to_render)
			dd.FileDebugMeDaddy(item_to_render)
			screen.RenderQueue = screen.RenderQueue[1:]
		}

		fmt.Printf(dd.MOVE_CURSOR_TO_POSITION, screen.CursorPosition.StartingRow, screen.CursorPosition.StartingCol)

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		screen.State.HandleKeypress(&screen, stdin_buffer)

		// case 'j':
		// 	move_cursor(&screen, screen.Active(), 1)
		// case 'k':
		// 	move_cursor(&screen, screen.Active(), -1)
		// case 'h':
		// 	move_cursor(&screen, &screen, -1)
		// 	screen.CursorPos = screen.Active().Active().GetPos()
		// case 'l':
		// 	move_cursor(&screen, &screen, 1)
		// 	screen.CursorPos = screen.Active().Active().GetPos()
	}
	// restart to default settings
	fmt.Print(dd.SHOW_CURSOR)
}

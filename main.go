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
			Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	content := dd.Window{
		Position: dd.Position{StartingRow: 1, StartingCol: uint(sidebar.Styles.Width) + 1},
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     int(float32(screen.Height)*0.7) + 1,
			Background: dd.MakeRGBBackground(69, 150, 100),
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	input_bar := dd.Window{
		Position: dd.Position{StartingRow: uint(content.Height) + 1, StartingCol: uint(sidebar.Width) + 1},
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     int(float32(screen.Height)*0.3) - 1,
			Background: dd.MakeRGBBackground(150, 150, 40),
			Border:     dd.Border{Style: dd.RoundedBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}
	status_bar_component := dd.Component{
		Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 2},
		Buffer:   dd.NORMAL,
		Styles: dd.Styles{
			Width:  screen.Width,
			Height: 1,
		},
	}

	status_bar := dd.Window{
		Position: dd.Position{StartingRow: uint(screen.Height), StartingCol: 1},
		Styles: dd.Styles{
			Width:      screen.Width,
			Height:     1,
			Background: dd.MakeRGBBackground(80, 40, 100),
		},
		Components: []dd.Component{status_bar_component},
	}

	item := dd.Component{
		Position: dd.Position{StartingRow: 3, StartingCol: uint(sidebar.StartingCol) + 2},
		Buffer:   "|Deez nuts|",
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height,
			Background: dd.MakeRGBBackground(80, 40, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	item_two := dd.Component{
		Position: dd.Position{StartingRow: 5, StartingCol: uint(sidebar.StartingCol) + 2},
		Buffer:   "|got em|",
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height,
			Background: dd.MakeRGBBackground(80, 40, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	sidebar.Components = []dd.Component{item, item_two}

	item_three := dd.Component{
		Position: dd.Position{StartingRow: 2, StartingCol: uint(content.StartingCol) + 2},
		Buffer:   "|SIMD|",
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height,
			Background: dd.MakeRGBBackground(80, 40, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	item_four := dd.Component{
		Position: dd.Position{StartingRow: 4, StartingCol: uint(content.StartingCol) + 2},
		Buffer:   "|Ligma?|",
		Styles: dd.Styles{
			Width:      screen.Width - sidebar.Styles.Width - 1,
			Height:     screen.Height,
			Background: dd.MakeRGBBackground(80, 40, 100),
			// Border:     dd.Border{Style: dd.BoldBorder, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	content.Components = []dd.Component{item_three, item_four}

	screen.Windows = append(screen.Windows, sidebar, content, input_bar)

	for _, window := range screen.Windows {
		for _, component := range window.Components {
			screen.RenderQueue = append(screen.RenderQueue, component.Render())
		}
		screen.RenderQueue = append(screen.RenderQueue, window.Render())
	}

	screen.Activate()
	screen.RenderQueue = append(screen.RenderQueue, status_bar.Render(), status_bar_component.Render())

	stdin_buffer := make([]byte, 1)
	for screen.EventLoopIsRunning {
		dd.DebugMeDaddy(&screen, fmt.Sprint(len(screen.Windows)))
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

package main

import (
	"fmt"
	"os"

	dd "github.com/Robert-Duck-by-BB-SR/talking_pond/internal/duck_dom"
	"golang.org/x/term"
)

// import (
// 	"bufio"
// 	"fmt"
// 	"image"
// 	_ "image/jpeg"
// 	_ "image/png"
// 	"log"
// 	"net"
// 	"os"
// 	"os/signal"
// )
//
// var frame_chars = []byte{' ', '`', '.', ',', '~', '+', '*', '&', '#', '@'}
//
// func main() {
//
// 	// Define the WebSocket server URL (replace with your server's address)
// 	serverAddr := "localhost:6969"
//
// 	// Dial the WebSocket server
// 	log.Printf("Connecting to %s...", serverAddr)
// 	conn, err := net.Dial("tcp", serverAddr)
// 	if err != nil {
// 		log.Fatalf("Failed to connect: %v", err)
// 	}
// 	defer conn.Close()
// 	log.Println("Connected to server")
//
// 	// Set up interrupt handling for graceful shutdown
// 	interrupt := make(chan os.Signal, 1)
// 	signal.Notify(interrupt, os.Interrupt)
//
// 	done := make(chan struct{})
// 	go func() {
// 		defer close(done)
// 		for {
// 			message, _, err := bufio.NewReader(conn).ReadLine()
// 			if err != nil {
// 				log.Println("Read error:", err)
// 				return
// 			}
// 			fmt.Println(string(message))
// 		}
// 	}()
//
// 	go func() {
// 		scanner := bufio.NewScanner(os.Stdin)
// 		for scanner.Scan() {
// 			data := []byte{127}
// 			data = append(data, scanner.Bytes()...)
// 			data = append(data, '\n')
// 			writer := bufio.NewWriter(conn)
// 			_, err := writer.Write(data)
// 			if err != nil {
// 				log.Fatalf("Write error: %v", err)
// 			}
// 			err = writer.Flush()
// 			if err != nil {
// 				log.Fatalf("Flush error: %v", err)
// 			}
// 		}
// 	}()
//
// 	// Wait for interrupt signal to close the connection
// 	for {
// 		select {
// 		case <-done:
// 			return
// 		case <-interrupt:
// 			log.Println("Interrupt received, closing connection")
// 			// Send a close message to the server
// 			data := []byte{127}
// 			data = append(data, "close please"...)
// 			_, err = bufio.NewWriter(conn).Write(data)
// 			if err != nil {
// 				log.Println("Close error:", err)
// 			}
// 			return
// 		}
// 	}
//
// }
//
// type CharMeDaddy struct {
// 	char, count, r, g, b byte
// }
//
// func encode_frame(img image.Image) []byte {
// 	orig_bounds := img.Bounds().Max
//
// 	scale_x := orig_bounds.X / 80
// 	scale_y := orig_bounds.Y / 40
// 	new_img_x := orig_bounds.X / scale_x
// 	new_img_y := orig_bounds.Y / scale_y
//
// 	encoded_data := []byte{}
// 	all_rle := []CharMeDaddy{}
//
// 	for y := range new_img_y {
// 		for x := range new_img_x {
// 			r, g, b, _ := img.At(x*scale_x, y*scale_y).RGBA()
// 			lum := (19595*r + 38470*g + 7471*b + 1<<15) >> 24
// 			indx := lum * uint32(len(frame_chars)) / 256
// 			// sliding window -> 5 bytes
// 			// 0 - char
// 			// 1 - repeat
// 			// 2 - r
// 			// 3 - g
// 			// 4 - b
// 			// 5 - new line
// 			if x == 0 {
// 				all_rle = append(all_rle, CharMeDaddy{frame_chars[indx], 1, uint8(r), uint8(g), uint8(b)})
// 			} else {
// 				curr_rle := &all_rle[len(all_rle)-1]
// 				if frame_chars[indx] == curr_rle.char &&
// 					uint8(r) == curr_rle.r &&
// 					uint8(g) == curr_rle.g &&
// 					uint8(b) == curr_rle.b {
// 					curr_rle.count += 1
// 				} else {
// 					all_rle = append(all_rle, CharMeDaddy{frame_chars[indx], 1, uint8(r), uint8(g), uint8(b)})
// 				}
// 			}
// 		}
// 		for _, el := range all_rle {
// 			encoded_data = append(encoded_data, el.char, el.count, el.r, el.g, el.b)
// 		}
// 		all_rle = []CharMeDaddy{}
// 		encoded_data = append(encoded_data, '\n')
// 	}
// 	return encoded_data
// }
//
// func decode_frame(enc_data []byte) {
// 	// sliding window -> 5 bytes
// 	// 0 - char
// 	// 1 - repeat
// 	// 2 - r
// 	// 3 - g
// 	// 4 - b
// 	// 5 - new line
// 	fmt.Print("\033[2J\033[H")
// 	for i := 0; i < len(enc_data); i += 5 {
// 		if enc_data[i] == '\n' {
// 			// or i -= 4
// 			i += 1
// 			fmt.Println()
// 			if i >= len(enc_data) {
// 				break
// 			}
// 		}
//
// 		for reps := 0; reps < int(enc_data[i+1]); reps += 1 {
// 			r := enc_data[i+2]
// 			g := enc_data[i+3]
// 			b := enc_data[i+4]
//
// 			var cell string = fmt.Sprintf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, enc_data[i])
// 			fmt.Print(cell)
// 		}
//
// 	}
// }

// func move_cursor(screen *dd.Screen, item dd.Renderable, direction int) {
// 	new_index := item.ActiveIndex() + direction
// 	if new_index >= 0 && new_index < len(screen.Windows) {
// 		active_item := item.Active()
// 		active_item.SetBackground("")
//
// 		item.SetActive(new_index)
//
// 		next_active_item := item.Active()
// 		next_active_item.SetBackground(dd.INVERT_STYLES)
//
// 		screen.RenderQueue = append(screen.RenderQueue, active_item, next_active_item)
// 		screen.CursorPos = next_active_item.GetPos()
// 	}
// }

func main() {
	old_state, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), old_state)

	dd.ClearScreen()

	screen := dd.Screen{}
	width, height, _ := term.GetSize(int(os.Stdin.Fd()))
	screen.MaxCols = width
	screen.MaxRows = height

	// item := dd.Button{
	// 	Pos:     dd.Position{Row: 3, Col: 1},
	// 	Content: "|Deez nuts|",
	// 	Styles:  dd.INVERT_STYLES,
	// }
	// screen.CursorPos = item.Pos

	// item_two := dd.Button{
	// 	Pos:     dd.Position{Row: 5, Col: 1},
	// 	Content: "|got em|",
	// }
	//
	// item_three := dd.Button{
	// 	Pos:     dd.Position{Row: 1, Col: 21},
	// 	Content: "|SIMD|",
	// 	Styles:  dd.INVERT_STYLES,
	// }
	//
	// item_four := dd.Button{
	// 	// NOTE: should we make item position relative or absolute?
	// 	Pos:     dd.Position{Row: 3, Col: 21},
	// 	Content: "|Ligma?|",
	// }
	//
	// sidebar := dd.Window{
	// 	Pos: dd.Position{Row: 0, Col: 0},
	// }
	//
	// sidebar.Children = append(sidebar.Children, &item, &item_two)
	// main_win.Children = append(main_win.Children, &item_three, &item_four)
	// screen.Windows = append(screen.Windows, &sidebar, &main_win)
	//
	// screen.RenderQueue = append(screen.RenderQueue, screen.Windows...)
	// screen.RenderQueue = append(screen.RenderQueue, sidebar.Children...)
	// screen.RenderQueue = append(screen.RenderQueue, main_win.Children...)

	sidebar := dd.Window{
		// NOTE: should we make item position relative or absolute?
		Position: dd.Position{StartingRow: 1, StartingCol: 1},
		Styles: dd.Styles{
			Width:      50,
			Height:     screen.MaxRows,
			Background: dd.DEBUG_STYLES,
			Border:     dd.Border{Width: 1, Style: dd.Solid, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	content := dd.Window{
		// NOTE: should we make item position relative or absolute?
		Position: dd.Position{StartingRow: 1, StartingCol: uint(sidebar.Styles.Width) + 2},
		Styles: dd.Styles{
			Width:      screen.MaxCols - sidebar.Styles.Width - 1,
			Height:     screen.MaxRows,
			Background: dd.MakeRGBBackground(69, 150, 100),
			Border:     dd.Border{Width: 1, Style: dd.Solid, Color: dd.MakeRGBTextColor(100, 100, 100)},
		},
	}

	screen.RenderQueue = append(screen.RenderQueue, sidebar.Render())
	screen.RenderQueue = append(screen.RenderQueue, content.Render())

	stdin_buffer := make([]byte, 1)
	running_on_my_nuts := true
	fmt.Print(dd.HIDE_CURSOR)
	for running_on_my_nuts {
		for len(screen.RenderQueue) > 0 {
			item_to_render := screen.RenderQueue[0]
			fmt.Print(item_to_render)
			screen.RenderQueue = screen.RenderQueue[1:]
		}

		fmt.Printf(dd.MOVE_CURSOR_TO_POSITION, screen.CursorPos.StartingRow, screen.CursorPos.StartingCol)

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		switch stdin_buffer[0] {
		case 'q':
			running_on_my_nuts = false
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
	}
}

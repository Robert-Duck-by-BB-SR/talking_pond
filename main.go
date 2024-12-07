package main

import (
	"fmt"
	dd "github.com/nodaridev/talking_pond/internal/duck_dom"
	"golang.org/x/term"
	"os"
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

func clear_screen() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[H")
}

func move_cursor(screen *dd.Screen, direction int, render_queue *[]dd.Item) {
	fmt.Print(direction, screen.Active_child_indx)
	if int(screen.Active_child_indx)+direction >= 0 &&
		int(screen.Active_child_indx)+direction < len(screen.Children) {
		active_item := &screen.Children[screen.Active_child_indx]
		active_item.Styles = "\033[0m"
		screen.Active_child_indx = uint(int(screen.Active_child_indx) + direction)

		next_active_item := &screen.Children[active_item.Active_child_indx]
		next_active_item.Styles = "\033[7m"

		*render_queue = append(*render_queue, *active_item, *next_active_item)
	}

	// fmt.Printf("\033[%d;%dH", row, col)
}

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println("Error enabling raw mode:", err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	clear_screen()

	stdin_buffer := make([]byte, 1)

	screen := dd.Screen{}
	item := dd.Item{
		Row:     10,
		Col:     5,
		Content: "|BLYA OTO DVIZ|",
	}

	item_two := dd.Item{
		Row:     12,
		Col:     5,
		Content: "|BLYA OTO DVIZ X2|",
	}
	screen.Children = append(screen.Children, item, item_two)

	render_queue := []dd.Item{}
	// now we render all

	render_queue = append(render_queue, screen.Children...)

	buffer := ""
	running_on_my_nuts := true
	for running_on_my_nuts {
		for len(render_queue) > 0 {
			item_to_render := render_queue[0]
			buffer += item_to_render.Render()
			render_queue = render_queue[1:]
		}

		if len(buffer) > 0 {
			fmt.Print(buffer)
			buffer = ""
		}

		_, err := os.Stdin.Read(stdin_buffer)
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}

		switch stdin_buffer[0] {
		case 'q':
			running_on_my_nuts = false
		case 'j':
			move_cursor(&screen, -1, &render_queue)
		case 'k':
			move_cursor(&screen, +1, &render_queue)
		}
	}
}

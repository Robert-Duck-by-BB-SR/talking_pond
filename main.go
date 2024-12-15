package main

import (
	"bufio"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

var frame_chars = []byte{' ', '`', '.', ',', '~', '+', '*', '&', '#', '@'}

func main() {

	var config [2]string
	server_port := ":6969"

	file, err := os.Open(".secrets")
	defer file.Close()
	log.Printf("%v, %v", file, err)
	if err != nil {
		i := 0
		for i < len(config) && config[i] == "" {
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("-> ")
			for scanner.Scan() {
				config[i] = scanner.Text()
				i += 1
				break
			}
		}
		f, err := os.Create(".secrets")
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(config[0] + "\n" + config[1] + "\n"))
		if err != nil {
			panic(err)
		}
	} else {
		i := 0
		scanner := bufio.NewScanner(file)
		for i < len(config) && config[i] == "" {
			for scanner.Scan() {
				fmt.Print(scanner.Text())
				config[i] = scanner.Text()
				i += 1
				break
			}
		}
	}
	server_addr := config[0] + server_port

	// Dial the WebSocket server
	log.Printf("Connecting to %s...", config[0])
	conn, err := net.Dial("tcp", server_addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server")

	// NOTE: testing requests
	data := []byte{1}
	data = append(data, "get:users"...)
	data = append(data, '\n')

	writer := bufio.NewWriter(conn)
	i, err := writer.Write(data)
	log.Println(i)
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Flush error: %v", err)
	}

	message, _, err := bufio.NewReader(conn).ReadLine()
	fmt.Println(string(message))
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	// Set up interrupt handling for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// NOTE: simulating conversation creation
	data = []byte{1}
	data = append(data, "create:conversation;key:"...)
	data = append(data, []byte(config[1])...)
	data = append(data, ";users:deeznuts"...)
	data = append(data, '\n')
	i, err = writer.Write(data)
	log.Println(i)
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Flush error: %v", err)
	}

	message, _, err = bufio.NewReader(conn).ReadLine()
	fmt.Println(string(message))
	if err != nil {
		log.Println("Read error:", err)
		return
	}

	convo := ""

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			message, _, err := bufio.NewReader(conn).ReadLine()
			fmt.Println(message)
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			if message[1] == 69 {
				command := strings.Split(string(message[2:]), ":")
				if command[0] == "convo" {
					convo = command[1]
				}
			} else {
				fmt.Println(string(message))
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			data := []byte{0, byte(len(config[1]))}
			data = append(data, []byte(config[1])...)
			data = append(data, []byte(convo)...)
			data = append(data, 0)
			data = append(data, scanner.Bytes()...)
			data = append(data, '\n')
			writer := bufio.NewWriter(conn)
			_, err := writer.Write(data)
			if err != nil {
				log.Fatalf("Write error: %v", err)
			}
			err = writer.Flush()
			if err != nil {
				log.Fatalf("Flush error: %v", err)
			}
		}
	}()

	// Wait for interrupt signal to close the connection
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt received, closing connection")
			// Send a close message to the server
			data := []byte{127}
			data = append(data, "close please"...)
			_, err = bufio.NewWriter(conn).Write(data)
			if err != nil {
				log.Println("Close error:", err)
			}
			return
		}
	}

}

// type CharMeDaddy struct {
// 	char, count, r, g, b byte
// }

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

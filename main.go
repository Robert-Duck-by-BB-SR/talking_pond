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

// var frame_chars = []byte{' ', '`', '.', ',', '~', '+', '*', '&', '#', '@'}

func create_convesation(key string, conn net.Conn) []byte {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(key)
	data.WriteByte(255)
	data.WriteString("create")
	data.WriteByte(255)
	data.WriteString("tredstart") // TODO: should be a param
	data.WriteByte('\n')
	send(conn, []byte(data.String()))

	message, _, err := bufio.NewReader(conn).ReadLine()
	fmt.Println(string(message))
	if err != nil {
		log.Println("Read error:", err)
		return nil
	}
	return message
}

func request_messages(key string, conn net.Conn, convo []byte) {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(key)
	data.WriteByte(255)
	data.WriteString("get")
	data.WriteByte(255)
	data.WriteString("message")
	data.WriteByte(255)
	data.WriteString(string(convo))
	data.WriteByte('\n')
	send(conn, []byte(data.String()))

	receive(conn)
}

func request_to_connect(key string, conn net.Conn) {
	var data strings.Builder
	data.WriteByte(2)
	data.WriteString(key)
	data.WriteByte('\n')
	send(conn, []byte(data.String()))
}

func send(conn net.Conn, data []byte) {
	writer := bufio.NewWriter(conn)
	_, err := writer.Write(data)
	if err != nil {
		log.Fatalf("Write error: %v", err)
	}
	err = writer.Flush()
	if err != nil {
		log.Fatalf("Flush error: %v", err)
	}
	log.Println("data sent")
}

func send_message(conn net.Conn, key, convo string, scanner *bufio.Scanner) {
	var data strings.Builder
	data.WriteByte(0)
	data.WriteString(key)
	data.WriteByte(255)
	data.WriteString(convo)
	data.WriteByte(255)
	data.WriteByte(0) // message type
	data.WriteByte(255)
	data.WriteString(scanner.Text())
	data.WriteByte('\n')
	send(conn, []byte(data.String()))
}

func receive(conn net.Conn) error {
	message, _, err := bufio.NewReader(conn).ReadLine()
	log.Println(message)
	if err != nil {
		log.Println("Read error:", err)
		return err
	}
	parts := strings.Split(string(message), string([]byte{254}))
	log.Println(parts)
	for _, part := range parts {
		m := strings.Split(part, string([]byte{255}))
		log.Println(len(m))
		if len(m) == 5 {
			fmt.Println(m[1], ":", m[4], "->", m[3])
		}
	}
	return nil
}

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

	log.Printf("Connecting to %s...", config[0])
	conn, err := net.Dial("tcp", server_addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server")

	// Set up interrupt handling for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	request_to_connect(config[1], conn)

	// convo := create_convesation(config[1], conn)
	// request_messages(config[1], conn, convo)
	convo := "a5c2fe80-22b7-495e-b2a6-79bf4eacf173"
	request_messages(config[1], conn, []byte(convo))

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			log.Println("starting receive")
			if err := receive(conn); err != nil {
				break
			}
		}
	}()

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			send_message(conn, config[1], string(convo), scanner)
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
			// FIXME: this is bs in the current setup
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

package tpsclient

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
)

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

type Client struct {
	server_port string
	server_addr string
	config      [2]string
}

func (client *Client) LoadClient() bool {
	client.server_port = ":6969"
	file, err := os.Open(".secrets")
	defer file.Close()
	log.Printf("%v, %v", file, err)

	// read config from a file into a struct
	if err == nil {
		i := 0
		scanner := bufio.NewScanner(file)
		for i < len(client.config) && client.config[i] == "" {
			for scanner.Scan() {
				client.config[i] = scanner.Text()
				i += 1
				break
			}
		}
		return false

		// i := 0
		// for i < len(client.config) && client.config[i] == "" {
		// 	scanner := bufio.NewScanner(os.Stdin)
		// 	for scanner.Scan() {
		// 		client.config[i] = scanner.Text()
		// 		i += 1
		// 		break
		// 	}
		// }
		// f, err := os.Create(".secrets")
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = f.Write([]byte(client.config[0] + "\n" + client.config[1] + "\n"))
		// if err != nil {
		// 	panic(err)
		// }
		// client.server_addr = client.config[0] + client.server_port
		// return true
	}
	return false
}

func (client *Client) placeholder() {

	log.Printf("Connecting to %s...", client.config[0])
	conn, err := net.Dial("tcp", client.server_addr)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to server")

	// Set up interrupt handling for graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	request_to_connect(client.config[1], conn)

	// convo := create_convesation(config[1], conn)
	// request_messages(config[1], conn, convo)
	convo := "a5c2fe80-22b7-495e-b2a6-79bf4eacf173"
	request_messages(client.config[1], conn, []byte(convo))

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
			send_message(conn, client.config[1], string(convo), scanner)
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

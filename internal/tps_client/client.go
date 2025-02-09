package tpsclient

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"
)

type Client struct {
	ServerPort   string
	ServerAddr   string
	Conversation string
	// {0: host, 1: key}
	Config [2]string
	Conn   net.Conn
}

func CreateConversation(client Client, users string) {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(client.Config[1])
	data.WriteByte(255)
	data.WriteString("create")
	data.WriteByte(255)
	data.WriteString(users)
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
}

func RequestMessages(client *Client) {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(client.Config[1])
	data.WriteByte(255)
	data.WriteString("get")
	data.WriteByte(255)
	data.WriteString("message")
	data.WriteByte(255)
	data.WriteString(client.Conversation)
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
}

func RequestConversations(client *Client) {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(client.Config[1])
	data.WriteByte(255)
	data.WriteString("get")
	data.WriteByte(255)
	data.WriteString("conversation")
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
}

func RequestUsers(client Client) {
	var data strings.Builder
	data.WriteByte(1)
	data.WriteString(client.Config[1])
	data.WriteByte(255)
	data.WriteString("get")
	data.WriteByte(255)
	data.WriteString("users")
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
}

func RequestToConnect(client *Client) {
	var data strings.Builder
	data.WriteByte(2)
	data.WriteString(client.Config[1])
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
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
}

func SendMessage(client *Client, message string) {
	var data strings.Builder
	data.WriteByte(0)
	data.WriteString(client.Config[1])
	data.WriteByte(255)
	data.WriteString(client.Conversation)
	data.WriteByte(255)
	data.WriteByte(0) // message type
	data.WriteByte(255)
	data.WriteString(message)
	data.WriteByte('\n')
	send(client.Conn, []byte(data.String()))
}

// returns [strings] split by 254 (item) separator
func Receive(conn net.Conn) string {
	scanner := bufio.NewScanner(conn)
	var message string
	if scanner.Scan() {
		message = scanner.Text()
	}
	return strings.Clone(message)
}

func (client *Client) LoadClient() bool {
	client.ServerPort = ":6969"
	file, err := os.Open(".secrets")
	defer file.Close()

	// read config from a file into a struct
	if err == nil {
		i := 0
		scanner := bufio.NewScanner(file)
		for i < len(client.Config) && client.Config[i] == "" {
			for scanner.Scan() {
				client.Config[i] = scanner.Text()
				i += 1
			}
		}
		client.ServerAddr = client.Config[0] + client.ServerPort
		return true
	}
	return false
}

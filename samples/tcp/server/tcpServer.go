package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	serverIP   = ""
	serverPort = "3206"
	serverType = "tcp4"
	bufferSize = 2048
)

func main() {
	listenAddress, err := net.ResolveTCPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}

	socket, err := net.ListenTCP(serverType, listenAddress)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("TCP Server Socket Program Example in Go\n")
	fmt.Printf("Press Ctrl+C or Cmd+C to stop the program\n")
	fmt.Printf("[%s] Listening on: %s\n", serverType, socket.Addr())

	defer socket.Close()

	for {
		connection, err := socket.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}

		go connectionHandler(connection)
	}
}

func connectionHandler(connection *net.TCPConn) {
	fmt.Printf("[%s] Receive connection from %s\n", serverType, connection.RemoteAddr())
	fmt.Printf("[%s] [Client: %s] Creating receive buffer for connection of size %d\n", serverType, connection.RemoteAddr(), bufferSize)
	receiveBuffer := make([]byte, bufferSize)

	defer connection.Close()

	receiveLength, err := connection.Read(receiveBuffer)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] [Client: %s] Received %d bytes of message\n", serverType, connection.RemoteAddr(), receiveLength)
	message := string(receiveBuffer[:receiveLength])

	fmt.Printf("[%s] [Client: %s] Message: %s\n", serverType, connection.RemoteAddr(), message)

	response, err := logic(message)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] [Client: %s] Sending Response: %s\n", serverType, connection.RemoteAddr(), response)
	_, err = connection.Write([]byte(response))
	if err != nil {
		log.Fatalln(err)
	}

}

func logic(input string) (string, error) {
	return strings.ToUpper(input), nil
}

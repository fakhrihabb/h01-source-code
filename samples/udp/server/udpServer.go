package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	serverIP   = ""
	serverPort = "54321"
	serverType = "udp4"
	bufferSize = 2048
)

func main() {
	listenAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}

	socket, err := net.ListenUDP(serverType, listenAddress)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("UDP Server Socket Program Example in Go\n")
	fmt.Printf("Press Ctrl+C or Cmd+C to stop the program\n")
	fmt.Printf("[%s] Listening on: %s\n", serverType, socket.LocalAddr())

	defer socket.Close()

	for {
		fmt.Printf("[%s] Creating receive buffer for next communication of size %d\n", serverType, bufferSize)
		receiveBuffer := make([]byte, bufferSize)

		receiveLength, address, err := socket.ReadFromUDP(receiveBuffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go connectionHandler(socket, address, receiveBuffer, receiveLength)
	}
}

func connectionHandler(socket *net.UDPConn, address *net.UDPAddr, receiveBuffer []byte, receiveLength int) {
	fmt.Printf("[%s] [Client: %s] Received %d bytes of message\n", serverType, address, receiveLength)
	message := string(receiveBuffer[:receiveLength])

	fmt.Printf("[%s] [Client: %s] Message: %s\n", serverType, address, message)

	response, err := logic(message)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("[%s] [Client: %s] Sending Response: %s\n", serverType, address, response)
	_, err = socket.WriteToUDP([]byte(response), address)
	if err != nil {
		log.Fatalln(err)
	}

}

func logic(input string) (string, error) {
	return strings.ToUpper(input), nil
}

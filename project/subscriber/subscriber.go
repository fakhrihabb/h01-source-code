package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP          = ""
	serverPort        = "3206"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2006142424"
	sslKeyLogFileName = "ssl-key.log"
)

func Handler(packet utils.LRTPIDSPacket) string {
	if packet.IsTrainArriving {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan tiba di Peron 1.", packet.Destination)
	}
	if packet.IsTrainDeparting {
		return fmt.Sprintf("Mohon perhatian, kereta tujuan %s akan diberangkatkan dari Peron 1.", packet.Destination)
	}
	return ""
}

func main() {
	localUdpAddress, err := net.ResolveUDPAddr(serverType, net.JoinHostPort(serverIP, serverPort))
	if err != nil {
		log.Fatalln(err)
	}
	socket, err := net.ListenUDP(serverType, localUdpAddress)
	if err != nil {
		log.Fatalln(err)
	}

	defer socket.Close()

	fmt.Printf("PIDS Display Node (Subscriber) - QUIC Server\n")
	fmt.Printf("[%s] Preparing UDP listening socket on %s\n", serverType, socket.LocalAddr())

	tlsConfig := &tls.Config{
		Certificates: utils.GenerateTLSSelfSignedCertificates(),
		NextProtos:   []string{appLayerProto},
	}
	listener, err := quic.Listen(socket, tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	defer listener.Close()

	fmt.Printf("[quic] Listening QUIC connections on %s\n", listener.Addr())

	for {
		connection, err := listener.Accept(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		go handleConnection(connection)
	}
}

func handleConnection(connection quic.Connection) {
	fmt.Printf("[quic] Receiving connection from %s\n", connection.RemoteAddr())

	for {
		stream, err := connection.AcceptStream(context.Background())
		if err != nil {
			fmt.Printf("[quic] [Client: %s] Connection closed or error accepting stream: %v\n", connection.RemoteAddr(), err)
			return
		}
		go handleStream(connection.RemoteAddr(), stream)
	}
}

func handleStream(clientAddress net.Addr, stream quic.Stream) {
	fmt.Printf("[quic] [Client: %s] Receive stream open request with ID %d\n", clientAddress, stream.StreamID())

	// Read the raw packet data
	buffer := make([]byte, bufferSize)
	n, err := stream.Read(buffer)
	if err != nil && err != io.EOF {
		fmt.Printf("[quic] Error reading from stream: %v\n", err)
		return
	}

	if n > 0 {
		// Decode the packet
		packet := utils.Decoder(buffer[:n])
		fmt.Printf("[quic] Received packet: Transaction ID %d, Train Number %d, Destination %s\n", 
			packet.TransactionId, packet.TrainNumber, packet.Destination)

		// Handle the packet and get the message
		message := Handler(packet)
		if message != "" {
			fmt.Println(message)
		}

		// Create ACK packet
		ackPacket := packet
		ackPacket.IsAck = true

		// Encode and send ACK
		ackData := utils.Encoder(ackPacket)
		_, err = stream.Write(ackData)
		if err != nil {
			fmt.Printf("[quic] Error sending ACK: %v\n", err)
		} else {
			fmt.Printf("[quic] ACK sent for Transaction ID %d\n", packet.TransactionId)
		}

		// Give the client time to read the ACK before closing
		time.Sleep(100 * time.Millisecond)
	}

	stream.Close()
}

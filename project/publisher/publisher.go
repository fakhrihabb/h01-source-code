package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/quic-go/quic-go"
	"jarkom.cs.ui.ac.id/h01/project/utils"
)

const (
	serverIP          = "3.83.102.64"
	serverPort        = "3206"
	serverType        = "udp4"
	bufferSize        = 2048
	appLayerProto     = "lrt-jabodebek-2006142424"
	sslKeyLogFileName = "ssl-key.log"
)

func main() {
	sslKeyLogFile, err := os.Create(sslKeyLogFileName)
	if err != nil {
		log.Fatalln(err)
	}
	defer sslKeyLogFile.Close()

	fmt.Printf("Station Control Node (Publisher) - QUIC Client\n")

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{appLayerProto},
		KeyLogWriter:       sslKeyLogFile,
	}
	connection, err := quic.DialAddr(context.Background(), net.JoinHostPort(serverIP, serverPort), tlsConfig, &quic.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	defer connection.CloseWithError(0x0, "No Error")

	fmt.Printf("[quic] Dialling from %s to %s\n", connection.LocalAddr(), connection.RemoteAddr())

	// Create Packet A: Train with trip number 42 destination Harjamukti arrives
	destination := "Harjamukti"
	packetA := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     1,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   true,
			IsTrainDeparting:  false,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}

	// Create Packet B: Train with trip number 42 destination Harjamukti departs
	packetB := utils.LRTPIDSPacket{
		LRTPIDSPacketFixed: utils.LRTPIDSPacketFixed{
			TransactionId:     2,
			IsAck:             false,
			IsNewTrain:        false,
			IsUpdateTrain:     false,
			IsDeleteTrain:     false,
			IsTrainArriving:   false,
			IsTrainDeparting:  true,
			TrainNumber:       42,
			DestinationLength: uint8(len(destination)),
		},
		Destination: destination,
	}

	// Send Packet A
	fmt.Printf("\n=== Sending Packet A (Train Arriving) ===\n")
	sendPacketAndReceiveACK(connection, packetA)

	// Send Packet B
	fmt.Printf("\n=== Sending Packet B (Train Departing) ===\n")
	sendPacketAndReceiveACK(connection, packetB)

	fmt.Printf("\n[quic] All packets sent successfully\n")
}

func sendPacketAndReceiveACK(connection quic.Connection, packet utils.LRTPIDSPacket) {
	stream, err := connection.OpenStreamSync(context.Background())
	if err != nil {
		log.Printf("[quic] Error opening stream: %v\n", err)
		return
	}
	defer stream.Close()

	fmt.Printf("[quic] Opened bidirectional stream %d to %s\n", stream.StreamID(), connection.RemoteAddr())

	// Encode and send packet
	packetData := utils.Encoder(packet)
	fmt.Printf("[quic] [Stream ID: %d] Sending packet with Transaction ID %d, Train Number %d, Destination %s\n", 
		stream.StreamID(), packet.TransactionId, packet.TrainNumber, packet.Destination)
	
	_, err = stream.Write(packetData)
	if err != nil {
		log.Printf("[quic] [Stream ID: %d] Error sending packet: %v\n", stream.StreamID(), err)
		return
	}
	fmt.Printf("[quic] [Stream ID: %d] Packet sent\n", stream.StreamID())

	// Read ACK
	receiveBuffer := make([]byte, bufferSize)
	receiveLength, err := stream.Read(receiveBuffer)
	if err != nil {
		log.Printf("[quic] [Stream ID: %d] Error receiving ACK: %v\n", stream.StreamID(), err)
		return
	}
	fmt.Printf("[quic] [Stream ID: %d] Received %d bytes of ACK from server\n", stream.StreamID(), receiveLength)

	// Decode ACK
	ackPacket := utils.Decoder(receiveBuffer[:receiveLength])
	fmt.Printf("[quic] [Stream ID: %d] Received ACK: Transaction ID %d, IsAck: %t\n", 
		stream.StreamID(), ackPacket.TransactionId, ackPacket.IsAck)
}

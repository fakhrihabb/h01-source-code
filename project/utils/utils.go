package utils

import (
	"bytes"
	"encoding/binary"
)

// Tipe data komponen boleh diubah, namun variabelnya jangan diubah
type LRTPIDSPacketFixed struct {
	TransactionId     uint16  // 16 bits
	IsAck             bool    // 1 bit
	IsNewTrain        bool    // 1 bit
	IsUpdateTrain     bool    // 1 bit
	IsDeleteTrain     bool    // 1 bit
	IsTrainArriving   bool    // 1 bit
	IsTrainDeparting  bool    // 1 bit
	TrainNumber       uint16  // 16 bits
	DestinationLength uint8   // 8 bits
}

type LRTPIDSPacket struct {
	LRTPIDSPacketFixed
	Destination string
}

func Encoder(packet LRTPIDSPacket) []byte {
	buffer := new(bytes.Buffer)

	// Write Transaction ID (16 bits)
	binary.Write(buffer, binary.BigEndian, packet.TransactionId)

	// Pack the 6 boolean flags into 1 byte (only using 6 bits)
	var flags uint8
	if packet.IsAck {
		flags |= 0x20 // bit 5
	}
	if packet.IsNewTrain {
		flags |= 0x10 // bit 4
	}
	if packet.IsUpdateTrain {
		flags |= 0x08 // bit 3
	}
	if packet.IsDeleteTrain {
		flags |= 0x04 // bit 2
	}
	if packet.IsTrainArriving {
		flags |= 0x02 // bit 1
	}
	if packet.IsTrainDeparting {
		flags |= 0x01 // bit 0
	}
	buffer.WriteByte(flags)

	// Write Train Number (16 bits)
	binary.Write(buffer, binary.BigEndian, packet.TrainNumber)

	// Write Destination Length (8 bits)
	buffer.WriteByte(packet.DestinationLength)

	// Write Destination string
	buffer.WriteString(packet.Destination)

	return buffer.Bytes()
}

func Decoder(rawMessage []byte) LRTPIDSPacket {
	buffer := bytes.NewReader(rawMessage)
	var packet LRTPIDSPacket

	// Read Transaction ID (16 bits)
	binary.Read(buffer, binary.BigEndian, &packet.TransactionId)

	// Read flags byte and unpack the boolean values
	var flags uint8
	binary.Read(buffer, binary.BigEndian, &flags)
	
	packet.IsAck = (flags & 0x20) != 0           // bit 5
	packet.IsNewTrain = (flags & 0x10) != 0      // bit 4
	packet.IsUpdateTrain = (flags & 0x08) != 0   // bit 3
	packet.IsDeleteTrain = (flags & 0x04) != 0   // bit 2
	packet.IsTrainArriving = (flags & 0x02) != 0 // bit 1
	packet.IsTrainDeparting = (flags & 0x01) != 0 // bit 0

	// Read Train Number (16 bits)
	binary.Read(buffer, binary.BigEndian, &packet.TrainNumber)

	// Read Destination Length (8 bits)
	binary.Read(buffer, binary.BigEndian, &packet.DestinationLength)

	// Read Destination string
	destinationBytes := make([]byte, packet.DestinationLength)
	buffer.Read(destinationBytes)
	packet.Destination = string(destinationBytes)

	return packet
}

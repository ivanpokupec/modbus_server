package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
)

type Header struct {
	TransactionId uint16
	ProtocolId    uint16
	Length        uint16
	UnitId        uint8
	FunctionCode  uint8
}

type ModbusMessage struct {
	StartAddress uint16
	Quantity     uint16
}

const port = ":502"

var holdingRegisters = make([]uint16, 65536)

func main() {
	holdingRegisters[0] = 53
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server started at %s\n", ln.Addr())
	conn, err := ln.Accept()
	if err != nil {
		log.Fatal(err)
	}

	for {
		reader := bufio.NewReader(conn)
		response := processMessage(reader)
		conn.Write([]byte(response + "\n"))

	}
}

func processMessage(message *bufio.Reader) string {
	var header Header
	err := binary.Read(message, binary.BigEndian, &header)
	if err != nil {
		log.Fatal(err)
	}

	s := strconv.FormatUint(uint64(header.Length), 10)
	fmt.Print("header length received:", s)

	s = strconv.FormatUint(uint64(header.ProtocolId), 10)
	fmt.Print("Protocol id received:", s)

	s = strconv.FormatUint(uint64(header.FunctionCode), 10)
	fmt.Print("Function code received:", s)

	if header.FunctionCode == 1 {
		// Read Coil Status
		fmt.Print("Function code 1 received")
	} else if header.FunctionCode == 2 {
		// Read Input Status
		fmt.Print("Function code 2 received")
	} else if header.FunctionCode == 3 {
		// Read Holding Registers
		var modbusMessage ModbusMessage
		err = binary.Read(message, binary.BigEndian, &modbusMessage)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print("Function code 3 received")
		fmt.Print("Start Address:", modbusMessage.StartAddress)
		fmt.Print("Quantity:", modbusMessage.Quantity)

		// Read the holding registers
		registers := holdingRegisters[modbusMessage.StartAddress : modbusMessage.StartAddress+modbusMessage.Quantity]
		fmt.Print("Registers:", registers)

	} else if header.FunctionCode == 4 {
		// Read Input Registers
		fmt.Print("Function code 4 received")
	} else if header.FunctionCode == 5 {
		// Write Single Coil
		fmt.Print("Function code 5 received")
	} else if header.FunctionCode == 6 {
		// Write Single Register
		fmt.Print("Function code 6 received")
	} else if header.FunctionCode == 15 {
		// Write Multiple Coils
		fmt.Print("Function code 15 received")
	} else if header.FunctionCode == 16 {
		// Write Multiple Registers
		fmt.Print("Function code 16 received")
	} else {
		fmt.Print("Invalid Function code received")
	}
	return ""
}

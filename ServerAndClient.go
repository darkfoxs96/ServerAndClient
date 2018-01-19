package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func main() {
	port := "9922"
	host := ""
	if len(os.Args) >= 2 {
		host = os.Args[1]
	}
	serverHost := "" + host
	readChan := make(chan []byte)
	writeChan := make(chan []byte)
	StartServerIfClient(port, serverHost, readChan, writeChan)
}

func StartServerIfClient(port string, serverHost string, readChan chan []byte, writeChan chan []byte) {
	if serverHost != "" {
		RPCClient(serverHost, writeChan)
	} else {
		RPCServer(port, readChan)
	}
}

func RPCServer(port string, readChan chan<- []byte) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go readConn(conn, readChan)
	}
}

func readConn(conn net.Conn, readChan chan<- []byte) {
	fullArray := make([]byte, 1000)
	var writePointer int
	var sizeMessage int
	var newLen int
	var remainArray []byte
	for {
		readArray := fullArray[writePointer:]
		n, err := conn.Read(readArray)
		if err != nil {
			fmt.Println(err)
			break
		}
		fullArray = fullArray[:writePointer+n]
		writePointer = 0
		for {
			if len(fullArray) == 4 && sizeMessage == 0 {
				sizeMessage = countSize(fullArray)
				newLen = sizeMessage
				break
			} else {
				if len(fullArray) == sizeMessage {
					readChan <- fullArray[:]
					sizeMessage = 0
					break
				}
			}
			if len(fullArray) > 4 && sizeMessage == 0 {
				sizeMessage = countSize(fullArray)
				if len(fullArray)-4 == sizeMessage {
					readChan <- fullArray[4 : sizeMessage+4]
					sizeMessage = 0
					break
				}
				if len(fullArray)-4 < sizeMessage {
					newLen = sizeMessage
					remainArray = fullArray[4:]
					writePointer = len(remainArray)
					break
				}
				if len(fullArray)-4 > sizeMessage {
					readChan <- fullArray[4 : sizeMessage+4]
					fullArray = fullArray[4+sizeMessage:]
					sizeMessage = 0
				}
			}
			if len(fullArray) < sizeMessage && sizeMessage != 0 {
				remainArray = fullArray[:]
				writePointer = len(remainArray)
				newLen = sizeMessage
				break
			}
			if len(fullArray) > sizeMessage && sizeMessage != 0 {
				readChan <- fullArray[:sizeMessage]
				fullArray = fullArray[sizeMessage:]
				if len(fullArray) == 4 {
					sizeMessage = countSize(fullArray)
					break
				}
				if len(fullArray) < 4 {
					remainArray = fullArray[:]
					newLen = len(remainArray)
					writePointer = len(remainArray)
					break
				}
				if len(fullArray) > 4 {
					sizeMessage = countSize(fullArray)
					fullArray = fullArray[4:]
				}
			}
			if len(fullArray) < 4 && sizeMessage == 0 {
				remainArray = fullArray[:]
				newLen = len(remainArray)
				writePointer = len(remainArray)
				break
			}
		}
		fullArray = make([]byte, newLen+1000)
		newLen = 0
		copy(fullArray, remainArray)
		remainArray = nil
	}
}

func countSize(header []byte) (size int) {
	size = int(header[0]) + int(header[1])*256 + int(header[2])*256*256 + int(header[3])*256*256*256
	return
}

func calcSize(size int) (header []byte) {
	header = make([]byte, 4)
	header[0] = byte(size % 256)
	size /= 256
	header[1] = byte(size % 256)
	size /= 256
	header[2] = byte(size % 256)
	header[3] = byte(size / 256)
	return
}

func residueIfEnd(read []byte, end int) bool {
	var by byte
	col := bytes.IndexByte(read, by)
	if col != end {
		return true
	}
	return false
}

func realizationResidue(read []byte, endIn int) (previous []byte) {
	buf := read[endIn:]
	var empty byte
	end := bytes.IndexByte(buf, empty)
	previous = buf[:end]
	return
}

func RPCClient(serverHost string, writeChan <-chan []byte) {
	for {
		conn, err := net.Dial("tcp", serverHost)
		if err != nil {
			fmt.Println(err)
		} else {
			connWrite(conn, writeChan)
			conn.Close()
		}
	}
}

func connWrite(conn net.Conn, writeChan <-chan []byte) {
	for {
		message := <-writeChan
		_, err := conn.Write(calcSize(len(message)))
		if err != nil {
			fmt.Println(err)
			break
		}
		_, err = conn.Write(message)
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}

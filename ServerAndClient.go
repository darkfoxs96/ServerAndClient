package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"strconv"
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
	KnowTheSize := false
	var end int
	var previous []byte
	for {
		read := make([]byte, 1000)
		_, err := conn.Read(read)
		if err != nil {
			fmt.Println(err)
			break
		}
		if KnowTheSize {
			bu := new(bytes.Buffer)
			_, err = bu.Write(previous)
			if err != nil {
				fmt.Println(err)
			}
			_, err = bu.Write(read[:end])
			if err != nil {
				fmt.Println(err)
			}
			glued := bu.Bytes()
			readChan <- glued
			KnowTheSize = false
			if residueAndEnd(read, end) {
				KnowTheSize, previous, end = realizationResidue(readChan, read, end)
			}
		} else {
			end = sizeDetermination(read)
			if end != -1 {
				KnowTheSize = true
			}
		}
	}
}

func sizeDetermination(read []byte) int {
	size := ""
	for i := 0; i <= 3; i++ {
		Val := string(read[i])
		numeral, err := strconv.Atoi(Val)
		if err != nil {
			return -1
		}
		s := strconv.Itoa(numeral)
		size += s
	}
	numeral := -1
	numeral, _ = strconv.Atoi(size)
	///-1 Not the size
	return numeral
}

func residueAndEnd(read []byte, end int) bool {
	var by byte
	col := bytes.IndexByte(read, by)
	if col != end {
		return true
	}
	return false
}

func realizationResidue(readChan chan<- []byte, read []byte, endIn int) (bool, []byte, int) {
	by := read[endIn:]
	var end int
	var previous []byte
	for {
		end = sizeDetermination(by)
		if end == -1 {
			var b byte
			col := bytes.IndexByte(by, b)
			previous = by[:col]
			return false, previous, 0
		} else {
			if !residueAndEnd(by, 4) {
				return true, previous, end
			}
			readChan <- by[4 : 4+end]
			if residueAndEnd(by, end) {
				by = by[4+end:]
			} else {
				return false, previous, 0
			}
		}
	}
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
		n, err := conn.Write(<-writeChan)
		if err != nil {
			fmt.Println(n, err)
			break
		}
	}
}

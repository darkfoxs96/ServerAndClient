package main

import (
	"bytes"
	"testing"
)

const testResult = `Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
Hpt
Pp
0003
SOE0002po0001
`

func TestServerAndClient(t *testing.T) {
	port := "9922"
	host := "127.0.0.1:9922"
	readChan := make(chan []byte, 1)
	writeChan := make(chan []byte, 1)
	//Create server
	go StartServerIfClient(port, "", readChan, writeChan)
	//Create client
	go StartServerIfClient(port, host, readChan, writeChan)
	//Create writeChanFunc
	go writeChanFunc(writeChan)
	//readChan
	in := bytes.NewBuffer(<-readChan)
	_, _ = in.Write([]byte("\n"))
	for i := 0; i < 31; i++ {
		message := <-readChan
		_, err := in.Write(message[:])
		if err != nil {
			t.Errorf("Exeption write in Buffer")
		}
		_, _ = in.Write([]byte("\n"))
	}
	result := in.String()
	if result != testResult {
		t.Errorf("test for OK Failed - results not match\nGot:\n%v\nExpected:\n%v", result, testResult)
	}
}

func writeChanFunc(writeChan chan<- []byte) {
	for s := 0; s < 8; s++ {
		writeChan <- []byte("Hpt")
		writeChan <- []byte("Pp")
		writeChan <- []byte("0003")
		writeChan <- []byte("SOE0002po0001")
	}
}

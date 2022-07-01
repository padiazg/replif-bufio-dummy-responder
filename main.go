package main

import (
	"bufio"
	"bytes"
	"fmt"
	"time"
)

func main() {
	var (
		command  []byte = []byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A}
		response []byte = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	)

	b := new(bytes.Buffer)
	rw := bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b))

	go dummyCommandResponder(rw, &command, &response)

	res, e0 := SendCommand(rw, &[]byte{0xFF, 0x86, 0x00, 0x47, 0x00, 0xC7, 0x03, 0x0F, 0x5A})
	if e0 != nil {
		fmt.Printf("TestSendCommand | Sending command: %v", e0)
	}

	fmt.Printf("%s", ToHex(res))
}

func dummyCommandResponder(rw *bufio.ReadWriter, c *[]byte, r *[]byte) {
	b := make([]byte, 9) // buffer to receive response
	for {
		if _, e := rw.Read(b); e != nil { // read response from tty
      fmt.Printf("b: %s\n", ToHex(b))
			// command received?
			if b[0] == (*c)[0] && b[1] == (*c)[1] && b[8] == (*c)[8] {
				// write response
				_, e := rw.Write(*r)
				if e != nil { // write commant to tty
					fmt.Print(e)
				}
				rw.Writer.Flush()
				break
			} // if (*r)[0] == (*c)[0] &&  ...
		} // if _, e := r0.Read(r); e != nil ....
		time.Sleep(250 * time.Millisecond)
	} // for ...
} // DummyCommandResponder ...

func SendCommand(rw *bufio.ReadWriter, c *[]byte) ([]byte, error) {
	fmt.Printf("SendCommand | c => %v\n", ToHex(*c))
	if _, e := rw.Write(*c); e != nil { // write commant to tty
		return nil, e
	}

	rw.Writer.Flush()                   // flush write buffer
	time.Sleep(1500 * time.Millisecond) // wait for the response
	r := make([]byte, 9)                // buffer to receive response
	if _, e := rw.Read(r); e != nil {   // read response from tty
		return nil, e
	}
	fmt.Printf("SendCommand | r => %v\n", ToHex(r))
	return r, nil
} // SendCommand ...

func ToHex(data []byte) string {
	var result = ""

	for _, c := range data {
		result = fmt.Sprintf("%s%#02x ", result, c)
	}

	return result
} // ToHex ...

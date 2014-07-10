package main

import (
  "io/ioutil"
  "strings"
  "time"
  "math/rand"
  "fmt"
  "net"
  "os"
)

/* WANT:
CLI args for port
CLI args for file name
TESTS
Logging - logrus / built in?   (log)
UDP
*/

/* Notes
Ginkgo
gofmt
testing TCP??????
CSP - Communicating sequential processes
How to write go (idomatic go)
$CDPATH
*/

func main() {
  rand.Seed(time.Now().UnixNano())
  l, err := net.Listen("tcp", "localhost:3333")
  if err != nil {
    fmt.Println("Error listening:", err.Error())
      os.Exit(1)
  }
  defer l.Close()
  fmt.Println("Listening")
  for {
    conn, err := l.Accept()
      if err != nil {
        fmt.Println("Error accepting: ", err.Error())
          os.Exit(1)
      }
    go handleRequest(conn)
  }
}

func handleRequest(conn net.Conn) {
  quote := randomQuote("wisdom.txt")
  conn.Write([]byte(quote))
  conn.Write([]byte("\r\n"))
  conn.Close()
}

func randomQuote(fileName string) string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
  quotes := strings.Split(string(file), "\n%\n")
  randQuoteIndex := rand.Intn(len(quotes))
	return quotes[randQuoteIndex]
}

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
  _ "github.com/codegangsta/cli"
  "github.com/Sirupsen/logrus"
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

var log = logrus.New()

func init() {
	rand.Seed(time.Now().UnixNano())
  log.Formatter = new(logrus.TextFormatter)
}

func main() {
  port := "3333"
	l, err := net.Listen("tcp", "localhost:" + port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
  log.Info("QOTD Server Started on Port " + port)
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

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
  "strconv"

  _ "github.com/codegangsta/cli"
  "github.com/Sirupsen/logrus"
  "github.com/nu7hatch/gouuid"
)

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
  requestUUID, err := uuid.NewV4()
  if err != nil {
    fmt.Println("error:", err)
      return
  }

  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Request Received")

	quoteId, quote := randomQuote("wisdom.txt")
	conn.Write([]byte(quote))
	conn.Write([]byte("\r\n"))
  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Quote #" + strconv.Itoa(quoteId) + " Served")

	conn.Close()
  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Connection Closed")
}

func randomQuote(fileName string) (int,string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	quotes := strings.Split(string(file), "\n%\n")
	randQuoteIndex := rand.Intn(len(quotes))
	return randQuoteIndex, quotes[randQuoteIndex]
}

/* Notes
Ginkgo
gofmt
testing TCP??????
CSP - Communicating sequential processes
How to write go (idomatic go)
*/

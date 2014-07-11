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

  "github.com/codegangsta/cli"
  "github.com/Sirupsen/logrus"
  "github.com/nu7hatch/gouuid"
)

var log = logrus.New()

func init() {
  rand.Seed(time.Now().UnixNano())
  log.Formatter = new(logrus.TextFormatter)
}

func main() {
  app := cli.NewApp()
  app.Name = "QOTD"
  app.Usage = "Run a QOTD Server"
  app.Flags = []cli.Flag {
    cli.StringFlag{"port,p", "3333", "port to bind the server to"},
  }

  app.Action = func(c *cli.Context) {
    port := c.String("port")
    fileName := "wisdom.txt"
    l, err := net.Listen("tcp", "localhost:" + port)
    if err != nil {
      log.Fatal("Error listening: ", err.Error())
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
      go handleRequest(conn, fileName)
    }
  }

  app.Run(os.Args)
}

func handleRequest(conn net.Conn, filename string) {
  requestUUID, err := uuid.NewV4()
  if err != nil {
    fmt.Println("error:", err)
    return
  }

  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Request Received")

  quoteId, quote := randomQuote(filename)
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

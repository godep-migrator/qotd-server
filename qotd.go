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

const (
  RFC865MaxLength = 512
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
    cli.BoolFlag{"strict", "quotes served in RFC 865 strict mode"},
  }

  app.Action = func(c *cli.Context) {
    port := c.String("port")
    fileName := c.Args()[0]
    quotes := loadQuotes(fileName)
    strictMode := c.Bool("strict")

    tcp, err := net.Listen("tcp", "localhost:" + port)
    if err != nil {
      log.Fatal("Error listening: ", err.Error())
      os.Exit(1)
    }
    defer tcp.Close()
    log.Info("QOTD Server Started on Port " + port)
    for {
      conn, err := tcp.Accept()
      if err != nil {
        fmt.Println("Error accepting: ", err.Error())
        os.Exit(1)
      }
      go serveRandomQuote(conn, quotes, strictMode)
    }
  }

  app.Run(os.Args)
}

func serveRandomQuote(conn net.Conn, quotes []string, strictMode bool) {
  requestUUID, err := uuid.NewV4()
  if err != nil {
    fmt.Println("error:", err)
    return
  }

  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Request Received")

  quoteId := rand.Intn(len(quotes))
  var quote = quotes[quoteId]
  if (strictMode && len(quote) > RFC865MaxLength){
    // 3 bytes for ..., 2 bytes for closing \r\n
    littleQuote := []byte(quote)[0:RFC865MaxLength-6]
    conn.Write(littleQuote)
    conn.Write([]byte("..."))
  } else {
    conn.Write([]byte(quote))
    log.WithFields(logrus.Fields{
      "request": requestUUID.String(),
      "client": conn.RemoteAddr().String(),
    }).Info("Quote #" + strconv.Itoa(quoteId) + " Served")
  }

  conn.Write([]byte("\r\n"))
  conn.Close()
  log.WithFields(logrus.Fields{
    "request": requestUUID.String(),
    "client": conn.RemoteAddr().String(),
  }).Info("Connection Closed")
}

func loadQuotes(fileName string) []string{
  file, err := ioutil.ReadFile(fileName)
  if err != nil {
    log.Fatal(err.Error())
  }
  quotes := strings.Split(string(file), "\n%\n")
  return quotes
}

/* Notes
Ginkgo
gofmt
testing TCP??????
CSP - Communicating sequential processes
How to write go (idomatic go)
*/

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
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
	app.Flags = []cli.Flag{
		cli.StringFlag{"port,p", "3333", "port to bind the server to"},
		cli.BoolFlag{"strict", "quotes served in RFC 865 strict mode"},
		cli.BoolFlag{"no-tcp", "server does not listen on tcp"},
		cli.BoolFlag{"no-udp", "server does not listen on udp"},
	}

	app.Action = func(c *cli.Context) {
		port := c.String("port")
		fileName := c.Args()[0]
		quotes := loadQuotes(fileName)
		strictMode := c.Bool("strict")
		startUdp := !c.Bool("no-udp")
		startTcp := !c.Bool("no-tcp")

		if startUdp {
			go listenForUdp(port, quotes, strictMode)
		}

		if startTcp {
			go listenForTcp(port, quotes, strictMode)
		}

		if startTcp || startUdp {
			// Keep this busy
			for {
				time.Sleep(100 * time.Millisecond)
			}
		} else {
			log.Fatal("Server not started on TCP or UDP, don't pass both --no-tcp and --no-udp")
		}
	}

	app.Run(os.Args)
}

func listenForTcp(port string, quotes []string, strictMode bool) {
	tcp, err := net.Listen("tcp", "localhost:"+port)
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

func listenForUdp(port string, quotes []string, strictMode bool) {
	udpService := ":" + port
	updAddr, udpErr := net.ResolveUDPAddr("udp", udpService)
	if udpErr != nil {
		log.Fatal("Error Resolving UDP Address: ", udpErr.Error())
		os.Exit(1)
	}
	updSock, udpErr := net.ListenUDP("udp", updAddr)
	if udpErr != nil {
		log.Fatal("Error listening: ", udpErr.Error())
		os.Exit(1)
	}
	defer updSock.Close()
	for {
		serveUDPRandomQuote(updSock, quotes, strictMode)
	}
}

func serveUDPRandomQuote(conn *net.UDPConn, quotes []string, strictMode bool) {
	requestUUID, err := uuid.NewV4()
	buf := make([]byte, 512)

	_, addr, err := conn.ReadFromUDP(buf[0:])

	if err != nil {
		panic(err)
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  addr.String(),
	}).Info("UDP Request Received")

	quote, quoteId := randomQuoteFormattedForDelivery(quotes, strictMode)
	conn.WriteToUDP([]byte(quote), addr)

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  addr.String(),
	}).Info("UDP Quote #" + strconv.Itoa(quoteId) + " Served")
}

func serveRandomQuote(conn net.Conn, quotes []string, strictMode bool) {
	requestUUID, err := uuid.NewV4()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("TCP Request Received")

	quote, quoteId := randomQuoteFormattedForDelivery(quotes, strictMode)
	conn.Write([]byte(quote))
	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("TCP Quote #" + strconv.Itoa(quoteId) + " Served")

	conn.Close()
	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("Connection Closed")
}

func randomQuoteFormattedForDelivery(quotes []string, strictMode bool) (string, int) {
	quoteId := rand.Intn(len(quotes))
	var quote = quotes[quoteId]
	if strictMode && len(quote) > RFC865MaxLength {
		// 3 bytes for ..., 2 bytes for closing \r\n
		quote = string([]byte(quote)[0 : RFC865MaxLength-6])
		quote = quote + "..."
	}

	quote = quote + "\r\n"
	return quote, quoteId
}

func loadQuotes(fileName string) []string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	quotes := strings.Split(string(file), "\n%\n")
	return quotes
}

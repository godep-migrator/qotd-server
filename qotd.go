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

		go listenForUdp(port, quotes, strictMode)
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

	app.Run(os.Args)
}

func listenForUdp(port string, quotes []string, strictMode bool) {
	udpService := ":" + port
	println(udpService)
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
	buf := make([]byte, 512)

	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	println("Invoked")

	quote := "Here you are on UDP getting a quote"

	conn.WriteToUDP([]byte(quote), addr)
	conn.WriteToUDP([]byte("\r\n"), addr)
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
	}).Info("Request Received")

	quoteId := rand.Intn(len(quotes))
	var quote = quotes[quoteId]
	if strictMode && len(quote) > RFC865MaxLength {
		// 3 bytes for ..., 2 bytes for closing \r\n
		littleQuote := []byte(quote)[0 : RFC865MaxLength-6]
		conn.Write(littleQuote)
		conn.Write([]byte("..."))
	} else {
		conn.Write([]byte(quote))
		log.WithFields(logrus.Fields{
			"request": requestUUID.String(),
			"client":  conn.RemoteAddr().String(),
		}).Info("Quote #" + strconv.Itoa(quoteId) + " Served")
	}

	conn.Write([]byte("\r\n"))
	conn.Close()
	log.WithFields(logrus.Fields{
		"request": requestUUID.String(),
		"client":  conn.RemoteAddr().String(),
	}).Info("Connection Closed")
}

func loadQuotes(fileName string) []string {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err.Error())
	}
	quotes := strings.Split(string(file), "\n%\n")
	return quotes
}

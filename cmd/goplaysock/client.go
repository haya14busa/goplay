package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	goplay "github.com/haya14busa/goplay/socket"

	"golang.org/x/net/websocket"
	"golang.org/x/tools/playground/socket"
)

const origin = "http://127.0.0.1/"

func main() {
	// Serve websocket playground server
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	mu := http.NewServeMux()

	mu.Handle("/", socket.NewHandler(mustParseURL(fmt.Sprintf("http://%s", l.Addr()))))
	s := http.Server{Handler: mu}
	go s.Serve(l)

	url := fmt.Sprintf("ws://%s/", l.Addr())

	config, err := websocket.NewConfig(url, origin)
	if err != nil {
		log.Fatal(err)
	}
	ws, err := websocket.DialConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	cli := goplay.Client{Conn: ws}
	code, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	cli.Run(string(code))
}

func mustParseURL(u string) *url.URL {
	r, _ := url.Parse(u)
	return r
}

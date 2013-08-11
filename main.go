package main

import (
	"flag"
	"fmt"
)

func main() {
	p := NewProxy(parseFlags())
	w := NewWatcher()

	w.StartApplications(p)
}

func parseFlags() string {
	port := flag.String("port", "8080", "Port where the proxy listens")
	host := flag.String("host", "localhost", "Host where the proxy is binded")

	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *host, *port)

	return addr
}

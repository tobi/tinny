package main

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"exec"
	"http"
	"time"
	"flag"
	"log"
	"net"
	"strconv"
)

var port *int = flag.Int("port", 3333, "port")
var executable *string = flag.String("executable", "server", "executable to launch. Must support -pXXXX for port assignment")
var compiler *string = flag.String("compiler", "make", "compiler to execute. make by default.")

func error(client io.WriteCloser, s string) {
	log.Printf("Error: %s", s)
	client.Write([]byte("HTTP/1.0 500 Internal Error\n"))
	client.Write([]byte("Content-Type: text/plain\n"))
	client.Write([]byte("\n\n" + s + "\n"))
	client.Close()
}

func pipeRequestResponse(server, client net.Conn) os.Error {

	// Read request from wire
	req, err := http.ReadRequest(bufio.NewReader(client))
	if err != nil {
		return err
	}

	rawReq, err := http.DumpRequest(req, true)
	if err != nil {
		return err
	}

	// forward it on
	server.Write(rawReq)

	// Read response from wire
	resp, err := http.ReadResponse(bufio.NewReader(server), req)
	if err != nil {
		return err
	}

	rawResp, err := http.DumpResponse(resp, true)
	if err != nil {
		return err
	}

  log.Printf("%s %s [%s]", req.Method, req.RawURL, resp.Status)

	// forward it on
	client.Write(rawResp)

	return nil
}

func forward(client net.Conn) {
	defer client.Close()

	cmd := exec.Command(*compiler)
	output, err := cmd.CombinedOutput()

	if err != nil {
		error(client, string(output))
		return
	}

  client_port := *port+1024

	cmd = exec.Command(*executable, "-p", strconv.Itoa(client_port))
	err = cmd.Start()

	if err != nil {
		error(client, err.String())
		return
	}

	// 20 retries, ~= 10 secs of attempts
	for i := 0; i < 50; i++ {
		server, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", client_port))
		if err == nil {

			if pipeRequestResponse(server, client) != nil {
				error(client, err.String())
			}

			return
		}

		time.Sleep(2e8) // 200ms
	}

  error(client, fmt.Sprintf("Timeout waiting to connect to port %d. Perhaps the executable doesn't support the -p parameter?", client_port))
}

func main() {
	flag.Parse()

	server, err := net.Listen("tcp", ":"+strconv.Itoa(*port))
	if err != nil {
		log.Fatal("Listen: ", err.String())
    return
	}

	log.Printf("Listening on port %d\n", *port)
	log.Printf(" -> forwarding to executable %s -p%d\n", *executable, *port+1024)

	for {
		conn, err := server.Accept()

		if err != nil {
			log.Fatal("accept:", err.String()) 
		}

		go forward(conn)
	}

}

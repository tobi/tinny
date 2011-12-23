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
  client.Write([]byte("\n\n"+s+"\n"))
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

  cmd = exec.Command(*executable, "-p", strconv.Itoa(*port+1024))
  err = cmd.Start()
     
  if err != nil { 
    error(client, err.String())
    return
  }    

  // 20 retries, ~= 10 secs of attempts
  for i := 0; i < 20; i++ {
    server, err := net.Dial("tcp", "127.0.0.1:5050")    
    if err == nil {

      err = pipeRequestResponse(server, client)
      if err != nil {
        error(client, err.String())
        return
      }

      log.Print("Forwarded...")  

      return
    }

    time.Sleep(5e8) // 500ms
  }

  error(client, err.String())    
}

func main() {
  flag.Parse()
 	
	server, err := net.Listen("tcp", ":"+ strconv.Itoa(*port))  
	if err != nil {
		log.Fatal("Listen: ", err.String())
	}

  fmt.Printf("Listening on port %d\n", *port)
  fmt.Printf(" -> forwarding to executable %s -p%d\n", *executable, *port+1024)

  for {
     conn, err := server.Accept()
    
     if err != nil {
       log.Fatal("accept:", err.String()) // TODO(r): exit?
     }

     go forward(conn)    
   }

}

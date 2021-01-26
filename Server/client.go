// socket client for golang
package main

import "net"
import "fmt"
import "bufio"
import "os"
import "strings"

func main() {

  // connect to server
  conn,_  := net.Dial("tcp", "10.20.0.197:82")
  go listener(conn)
  for {
    // what to send?
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your command: ")
    text,_  := reader.ReadString('\n')
    // send to server
    text = strings.TrimSuffix(string(text),"\n")
    fmt.Fprintf(conn, text)
    
  }
}

func listener(connection net.Conn) {
	defer connection.Close()
	for {
		reply := make([]byte, 1024)
    	_, _ = connection.Read(reply)
    	println("\nReply=", string(reply))
    	fmt.Print("Enter your command: ")
	}
	return
}
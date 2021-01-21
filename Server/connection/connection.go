package connection

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
  	"bytes"
)

//UGHHHHH ...
var (
	CHOST = "10.20.0.197"
	CPORT = "82"
	CNET  = "tcp"
)

//MSG is used in transporting individual messages
type Message struct {
	username string
	msg      string
	isLast   bool

	conn *net.TCPConn
}

//Conns contains all connections
type AllConnections struct {
	connections map[string]*net.TCPConn
	sync.RWMutex
}

//User contains username
type User struct {
	Username string
}

//New ...
func (c *AllConnections) New(key string, conn *net.TCPConn) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.connections[key]; !ok {
		c.connections[key] = conn
		return true
	}
	return false
}

//Delete ...
func (c *AllConnections) Delete(key string) bool {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.connections[key]; ok {
		delete(c.connections, key)
		return true
	}
	return false
}

//Read ...
func (c *AllConnections) Read(key string) *net.TCPConn {
	c.RLock()
	defer c.RUnlock()

	if conn, ok := c.connections[key]; ok {
		return conn
	}

	return nil
}

func ListenForNewConnections() {
	//Create addr
	addr, err := net.ResolveTCPAddr(CNET, CHOST+":"+CPORT)
	handleError(err)

	//Create listener
	tcpListener, err := net.ListenTCP(CNET, addr)
	handleError(err)
	defer tcpListener.Close()

	//Create broadcaster
	var broadcaster = make(chan *Message, 1)
	defer close(broadcaster)
	c := &AllConnections{ connections: make(map[string]*net.TCPConn) }
	go c.broadcast(broadcaster)

	fmt.Println("Listening on " + CHOST + ":" + CPORT)
	for {
				// Listen for an incoming connection.
		conn, err := tcpListener.AcceptTCP()
		handleError(err)
		c.New(conn.RemoteAddr().String(), conn)

		// Handle connections in a new goroutine.
		go handleRequest(conn, broadcaster)
		fmt.Println("Handel connection!")

	}

	return
}

func (c *AllConnections) broadcast(messages chan *Message) {
	//Loop continually sending messages
	for {
		msg := <-messages
		//If user left, remove from conn pool
		if msg.isLast {
		  	c.Delete(msg.conn.RemoteAddr().String())
	  		//continue
		}

		for _, conn := range c.connections {
			//Write user message to group
	  		fmt.Println(msg.msg + "\n")
			_, err := conn.Write([]byte(msg.msg))
			if err != nil {
				c.Delete(msg.conn.RemoteAddr().String())
			}
		}
	}
}

func initUser(conn *net.TCPConn) *User {
	//Read data from connection
	buf := make([]byte, 2048)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return nil
	}

  buf = bytes.Trim(buf, "\x00")

	data := string(buf)
	if data == "" {
		return nil
	}

	info := strings.Split(data, ":")
	if len(info) < 2 || info[0] != "iam" {
		return nil
	}

	return &User{Username: info[1]}
}

func handleRequest(conn *net.TCPConn, messages chan *Message) {
	defer conn.Close()

	//First, get th euser name from the user. Client MUST first send an iam signal
	user := initUser(conn)
	if user == nil {
		return
	}
	fmt.Println(user.Username, " has joined\n")

	for {
		//Now it loops forver as it waits to reveive messages from that client

		buf := make([]byte, 2048)
		//Waits here until it recieves a message
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			//Remove connection from list
			//Client closed, alert other users
	 		msg := &Message{ username: user.Username, msg: "Disconnected!", conn: conn, isLast: true }
			messages <- msg
			return
		}

		buf = bytes.Trim(buf, "\x00")
		//Check if msg string has content and that it is a msg type
		content := string(buf)

		RouteRecievedMessage(user, content)
	}
	return
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

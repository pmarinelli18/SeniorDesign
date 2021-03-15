package connection

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
  	"bytes"
)

var (
	CHOST = "localhost"
	CPORT = "80"
	CNET  = "tcp"
)

var allConnections = &AllConnections{ connections: make(map[string]*net.TCPConn) }

var hardwareConnection *net.TCPConn

//Conns contains all connections
type AllConnections struct {
	connections map[string]*net.TCPConn
	sync.RWMutex
}


func createNewConnection(key string, conn *net.TCPConn) bool {
	allConnections.Lock()
	defer allConnections.Unlock()

	if _, ok := allConnections.connections[key]; !ok {
		allConnections.connections[key] = conn
		return true
	}
	return false
}

func SetHardwareConnection(conn *net.TCPConn){
	hardwareConnection = conn
}

//Delete ...
func DeleteConnection(withID string) bool {
	allConnections.Lock()
	defer allConnections.Unlock()

	if _, ok := allConnections.connections[withID]; ok {
		delete(allConnections.connections, withID)
		return true
	}
	return false
}

func SendMessageToAll(message string){
	for _, conn := range allConnections.connections {
		//Write user message to group
		_, err := conn.Write([]byte(message))
		if err != nil {
			DeleteConnection(conn.RemoteAddr().String())
		}
	}
}

func SendMessageToDevices(message []byte, devices []ConnectedDevices){
	for _, device := range devices{
		//Find the connection from allConnections
		for _, conn := range allConnections.connections {
			if (conn.RemoteAddr().String() == device.ipAddress){
				//Send out the message to that device
				conn.Write(message)
				break
			}
		}
	}
}

func SendMessageToHardware(message []byte){
	hardwareConnection.Write(message)
}

func GetConnection(device ConnectedDevices) *net.TCPConn{
		//Find the connection from allConnections
	for _, conn := range allConnections.connections {
		if (conn.RemoteAddr().String() == device.ipAddress){
			//Send out the message to that device
			return conn
		}
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

	fmt.Println("Listening on " + CHOST + ":" + CPORT)
	for {
		// Listen for an incoming connection.
		newConnection, err := tcpListener.AcceptTCP()
		handleError(err)
		
		//newConnection.RemoteAddr().String() is the IP address of the newly connected device. It will be used 
		//to identify the device
		createNewConnection(newConnection.RemoteAddr().String(), newConnection)
		
		// Handle connections in a new goroutine.
		go handleRequest(newConnection)
	}
	return
}

func handleRequest(connection *net.TCPConn) {
	defer connection.Close()
	fmt.Println(connection.RemoteAddr().String(), " has connected")

	for {
		//Now it loops forver as it waits to reveive messages from that client
		buf := make([]byte, 2048)
		//Waits here until it recieves a message
		_, err := connection.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			fmt.Println(connection.RemoteAddr().String(), " has disconnected")
			EndGame(connection)
			DeleteConnection(connection.RemoteAddr().String())
			return
		}

		buf = bytes.Trim(buf, "\x00")
		//Check if msg string has content and that it is a msg type
		content := string(buf)

		RouteRecievedMessage(connection, content)
	}
	return
}

func handleError(err error) {
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}



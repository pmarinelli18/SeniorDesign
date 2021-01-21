package main
import (
	Connection "./connection"
)

func main() {
	//Create addr
	Connection.ListenForNewConnections()

	return
}
package main
import (
	Connection "./connection"
)

func main() {
	//Create addr
	Connection.MakeDatabaseConnection()
	Connection.ListenForNewConnections()

	return
}
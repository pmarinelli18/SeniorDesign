package connection
import (
	"fmt"
	"strings"
	"net/url"
	"net"
	"encoding/json"
	"strconv"
)

/*
Messages are parsed like a URL, with a body and optional parameters
Examples of URLs are: 		msg?content=This is the message to print
					   		login?userName=MyName&password=12342
							triggerHardware/radarEffect	
*/

func RouteRecievedMessage(connection *net.TCPConn, messageContent string){

	urlRequest, err := url.Parse(messageContent)
	if err != nil {
        panic(err)
    }
    path := urlRequest.Path
    splitPath := strings.Split(path, "/")
	
	parameters, _ := url.ParseQuery(urlRequest.RawQuery)


	if len(splitPath) != 0 {
		if splitPath[0] == "msg" && len(parameters["string"]) > 0{
			fmt.Println("Display message")
			fmt.Println(connection.RemoteAddr().String() + ": " + parameters["string"][0])
			resondBack("displayMessage", connection, true)
		} else if splitPath[0] == "auth" {
			authRouter(splitPath, parameters, connection)
			
		} else if splitPath[0] == "get" {
			getRouter(splitPath, parameters, connection)
			
		}else if splitPath[0] == "game" {
			gameRouter(splitPath, parameters, connection)
			
		} else if splitPath[0] == "echo"  && len(parameters["string"]) > 0{
			fmt.Println("Echo message to all")
			SendMessageToAll(parameters["string"][0])
			resondBack("echo message", connection, true)
		} else{
			resondBack("bad api call", connection, false)
		}
	} else{
		resondBack("bad api call", connection, false)
	}
}

func getRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	if len(splitPath) > 1 && splitPath[1] == "getUsers" {
		fmt.Println("Get Users")
		getUsers()
	}
}
func authRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	result := false
	tag := "bad api call"

	if  len(splitPath) > 1 && splitPath[1] == "login" && len(parameters["userName"]) > 0 && len(parameters["password"]) > 0{
		fmt.Println("Login")
		tag = "login"
		result = CheckIfValidLogin(parameters["userName"][0], parameters["password"][0])
	} else if len(splitPath) > 1 && splitPath[1] == "createAccount" && len(parameters["userName"]) > 0 && len(parameters["password"]) > 0{
		fmt.Println("Create Account")
		tag = "createAccount"
		result = CreateNewAccount(parameters["userName"][0], parameters["password"][0])
	} else if len(splitPath) > 1 && splitPath[1] == "getUsers"{
		fmt.Println("Get Users")
		getUsers()
		return
	}

	if result {
		resondBack(tag, connection, true)
		InitNewUser(parameters["userName"][0], connection)
	} else{
		resondBack(tag, connection, false)
	}
}

func gameRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	result := false
	if  len(splitPath) > 1 && splitPath[1] == "recordResult" && len(parameters["miniGame"]) > 0 && len(parameters["score"]) > 0{
		fmt.Println("Saving game result")
		
		//result = CheckIfValidLogin(parameters["userName"][0], parameters["password"][0])
	} else if len(splitPath) > 1 && splitPath[1] == "repairShip" {
		fmt.Println("Repairing ship")
		// Table increases ship's health by 25
	} else if  len(splitPath) > 1 && splitPath[1] == "hitConfirm" && len(parameters["weaponType"]) > 0 {
		fmt.Println("Attack landed")
		// Update health of ship based on which weapoon type was used
	} else if  len(splitPath) > 1 && splitPath[1] == "navigation" && len(parameters["pos"]) > 0 {
                fmt.Println("Changing position")
                // Updates the navigation position of the ship
	} else if  len(splitPath) > 1 && splitPath[1] == "hackRadar" {
                fmt.Println("Radar Disabled")
                // Switches the radar of a ship to "0" (off)
	}
	if result {
		resondBack("TODO", connection, true)
	} else{
		resondBack("TODO", connection, false)
	}
}

func resondBack(tag string, connection *net.TCPConn, success bool){
	mapD := map[string]string{"id": tag, "result": strconv.FormatBool(success)}
	mapB, _ := json.Marshal(mapD)
	_, err := connection.Write([]byte(mapB))
	if err != nil {
		DeleteConnection(connection.RemoteAddr().String())
	}
}

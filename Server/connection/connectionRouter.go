package connection
import (
	"fmt"
	"strings"
	"net/url"
	"net"
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
			resondBack(connection, true)
		} else if splitPath[0] == "auth" {
			authRouter(splitPath, parameters, connection)
			
		} else if splitPath[0] == "get" {
			getRouter(splitPath, parameters, connection)
			
		}else if splitPath[0] == "game" {
			gameRouter(splitPath, parameters, connection)
			
		} else if splitPath[0] == "echo"  && len(parameters["string"]) > 0{
			fmt.Println("Echo message to all")
			SendMessageToAll(parameters["string"][0])
			resondBack(connection, true)
		} else{
			resondBack(connection, false)
		}
	} else{
		resondBack(connection, false)
	}
}

func getRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	if len(splitPath) > 1 && splitPath[1] == "getUsers" && len(parameters["userName"]) > 0 {
		fmt.Println("Get Users")
		getUsers(parameters["userName"][0], connection)
	}
}
func authRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	result := false
	if  len(splitPath) > 1 && splitPath[1] == "login" && len(parameters["userName"]) > 0 && len(parameters["password"]) > 0{
		fmt.Println("Login")
		result = CheckIfValidLogin(parameters["userName"][0], parameters["password"][0])
	} else if len(splitPath) > 1 && splitPath[1] == "createAccount" && len(parameters["userName"]) > 0 && len(parameters["password"]) > 0{
		fmt.Println("Create Account")
		result = CreateNewAccount(parameters["userName"][0], parameters["password"][0])
	} else if len(splitPath) > 1 && splitPath[1] == "getUsers" && len(parameters["userName"]) > 0 {
		fmt.Println("Get Users")
		getUsers(parameters["userName"][0], connection)
	}

	if result {
		resondBack(connection, true)
		InitNewUser(parameters["userName"][0], connection)
	} else{
		resondBack(connection, false)
	}
}

func gameRouter(splitPath []string, parameters url.Values, connection *net.TCPConn){
	result := false
	if  len(splitPath) > 1 && splitPath[1] == "recordResult" && len(parameters["miniGame"]) > 0 && len(parameters["score"]) > 0{
		fmt.Println("Saving game result")
		
		//result = CheckIfValidLogin(parameters["userName"][0], parameters["password"][0])
	}
	if result {
		resondBack(connection, true)
	} else{
		resondBack(connection, false)
	}
}

func resondBack(connection *net.TCPConn, success bool){
	message := "Failed"
	if success {
		message = "Success"
	}
	_, err := connection.Write([]byte(message))
	if err != nil {
		DeleteConnection(connection.RemoteAddr().String())
	}
}


func sendMessageToAllConnections(){

}
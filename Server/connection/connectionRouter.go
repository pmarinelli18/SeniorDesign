package connection
import (
	"fmt"
	"strings"
	"net/url"
)

/*
Messages are parsed like a URL, with a body and optional parameters
Examples of URLs are: 		msg?content=This is the message to print
					   		login?userName=MyName&password=12342
							triggerHardware/radarEffect	
*/

func RouteRecievedMessage(user *User, messageContent string){

	urlRequest, err := url.Parse(messageContent)
	if err != nil {
        panic(err)
    }
    path := urlRequest.Path
    splitPath := strings.Split(path, "/")
	
	parameters, _ := url.ParseQuery(urlRequest.RawQuery)


	if len(splitPath) != 0 {
		if splitPath[0] == "msg"{
			fmt.Println("Display message")
			fmt.Println(user.Username + ": " + parameters["string"][0])
		} else if splitPath[0] == "login" {
			fmt.Println("Login")
			CheckIfValidLogin("userName", "password")
		}
	}
}

func SendMessageToAllConnections(){

}
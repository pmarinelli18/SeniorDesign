package connection
import (
	"fmt"
	"strings"
)


func RouteRecievedMessage(user *User, messageContent string){

	info := strings.Split(messageContent, ":")

	if len(info) >= 2 {
		if info[0] == "msg"{
			fmt.Println("Display message")
			fmt.Println(user.Username + ": " + info[1])

		} else if info[0] == "login" && len(info) == 3 {
			fmt.Println("Login")
			fmt.Println("Username: " + info[1] + ", Password: " + info[2])
		}
	}
}

func SendMessageToAllConnections(){

}
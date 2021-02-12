package connection
import (
	"fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net"
)

var databaseConnection *sql.DB 

/*
type Tag struct {
    userName   string    `json:"userName"`
}
*/
type LogIn struct {
    success   int   `json:"success"`
}

type ConnectedDevices struct {
    ipAddress   string
    userName   string
}

func MakeDatabaseConnection(){
	var err error

    databaseConnection, err = sql.Open("mysql", "root:password@tcp(localhost)/BattleShip")

	// if there is an error opening the connection, handle it
    if err != nil {
        fmt.Println("ERROR: Cannot connect to the database")
        panic(err.Error())
    } else{
    	fmt.Println("Successful connected to database!")
    }
    err = databaseConnection.Ping()
    if err != nil{
        fmt.Println("ERROR: Cannot connect to the database on Ping")
    }
}

func CheckIfValidLogin(userName string, password string) bool{
	results, err := databaseConnection.Query("SELECT Count(*) success from Accounts where userName = \""+ userName + "\" AND password=\"" + password + "\";")
    if err != nil {
        return false
    }

    var logIn LogIn
    results.Next()
    err = results.Scan(&logIn.success)

    if logIn.success == 1{
        return true
    } else{
        return false
    }
}

func EndGame(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("DELETE FROM BoatState WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
}


func InitNewUser(userName string, ipAddress *net.TCPConn){
    //Init the user's game
        _, _ = databaseConnection.Query("INSERT INTO BoatState(IpAddress, UserName) Values (\""+ ipAddress.RemoteAddr().String() + "\", \""+ userName + "\");")
        dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true;")
        devices := make([]ConnectedDevices, 0)
        userNames := ""
        for dbConnections.Next() {
            var newDevice ConnectedDevices
            _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
            devices = append(devices, newDevice)
            userNames += newDevice.userName + " "
        }

        SendMessageToDevices(userNames, devices)
}
func CreateNewAccount(userName string, password string) bool{
    _, err := databaseConnection.Query("INSERT INTO Accounts(userName, password) VALUES (\""+ userName + "\", \"" + password + "\");")
    if err != nil{
        fmt.Println(err.Error())
        return false
    }
    return true
}



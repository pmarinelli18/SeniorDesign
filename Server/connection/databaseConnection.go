package connection
import (
	"fmt"
	//"strconv"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net"
    "encoding/json"
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
        err = databaseConnection.Ping()
        if err != nil{
            fmt.Println("ERROR: Cannot connect to the database on Ping")
        }else{
    	    fmt.Println("Successful connected to database!")
        }
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
        getUsers()
}

func FireWeapon(wep string, ipAddress *net.TCPConn) {
	if wep == "1" {
		fmt.Println("Cannon Incoming")
		_, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth - 25 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
		_, _ = databaseConnection.Query("UPDATE BoatState SET NumberOfCannons = NumberOfCannons - 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
	} else if wep == "2" {
		fmt.Println("Torpedo Incoming")
		_, _ = databaseConnection.Query("UPDATE BoatState SET TorpedoDamage = 40 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
		_, _ = databaseConnection.Query("UPDATE BoatState SET TorpedoState = \"Cooldown\" WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
	} else if wep == "3" {
                fmt.Println("Mounted MG Incoming")
                 _, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth - 10 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
	} else {
		fmt.Println("Invalid Weapon Type")
	}
}

func RepairShip(ipAddress *net.TCPConn) {
	_, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth + 25 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
}

func HackRadar(ipAddress *net.TCPConn) {
        _, _ = databaseConnection.Query("UPDATE BoatState SET RadarState = \"Hacked\" WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
}

func FixRadar(ipAddress *net.TCPConn) {
        _, _ = databaseConnection.Query("UPDATE BoatState SET RadarState = \"Enabled\" WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
}

func ChangePosition(pos string,ipAddress *net.TCPConn) {
	_, _ = databaseConnection.Query("UPDATE BoatState SET navigationPosition = " + pos + "  WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
}

func getUsers(){
    dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true;")
    devices := make([]ConnectedDevices, 0)
    userNames := make([]string, 0)

    for dbConnections.Next() {
        var newDevice ConnectedDevices
        _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
        devices = append(devices, newDevice)
        userNames = append(userNames, newDevice.userName)
    }

    mapD := map[string]interface{}{
        "id": "connectedDevices",
        "userNames": userNames,
    }
    mapB, _ := json.Marshal(mapD)
    

    SendMessageToDevices(mapB, devices)
}
func CreateNewAccount(userName string, password string) bool{
    _, err := databaseConnection.Query("INSERT INTO Accounts(userName, password) VALUES (\""+ userName + "\", \"" + password + "\");")
    if err != nil{
        fmt.Println(err.Error())
        return false
    }
    return true
}



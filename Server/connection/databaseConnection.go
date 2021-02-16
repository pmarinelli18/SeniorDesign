package connection
import (
	"fmt"
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

type BoatState struct {
    navigationPosition   string
    radarState   string
    shipHealth   string
    numberOfCannons   string
    torpedoState   string
    opponentHealth   string
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

func GetBoatState(ipAddress *net.TCPConn){
    boatStateResult := databaseConnection.QueryRow("SELECT navigationPosition, RadarState radarState, ShipHealth shipHealth, NumberOfCannons numberOfCannons, TorpedoState torpedoState, opponentBoat.opponentHealth from BoatState join (select ShipHealth opponentHealth from BoatState where IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1) as opponentBoat where IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    
    var boatState BoatState
    _ = boatStateResult.Scan(&boatState.navigationPosition, &boatState.radarState, &boatState.shipHealth, &boatState.numberOfCannons, &boatState.torpedoState, &boatState.opponentHealth)
    fmt.Println(boatState.shipHealth)
    cannonState := "Enabled"
    if boatState.numberOfCannons == "0"{
        cannonState = "Reloading"
    }
    mapD := map[string]interface{}{
        "id":"boatState",
        "boatHealth": map[string]interface{}{
            "yourHeath": boatState.shipHealth,
            "opponentHealth": boatState.opponentHealth,
        },
        "stateOfBoatFeatures": map[string]interface{}{
            "radar": boatState.radarState,
            "torpedo": boatState.torpedoState,
            "cannons": map[string]interface{}{
                "state": cannonState,
                "numberOfCannons": boatState.numberOfCannons,
            },
        },
        "boatPosition": boatState.navigationPosition,
    }
    mapB, _ := json.Marshal(mapD)
    _, _ = ipAddress.Write([]byte(mapB))
}




package connection
import (
	"fmt"
	"strconv"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
    "net"
    "encoding/json"
    "math/rand"
)

var databaseConnection *sql.DB 

var p1Address = ""
var p2Address = ""

var p1JustShotCannon = false
var p2JustShotCannon = false

/*
type Tag struct {
    userName   string    `json:"userName"`
}
*/
type LogIn struct {
    success   int   `json:"success"` 
}

type PlayersFinished struct {
    IpAddress   string
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
    opponentUserName   string
    userName     string
    ipAddress    string
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
            clearAccountsTable();
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

func RemovePlayerFromDB(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("DELETE FROM BoatState WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
} 

func EndGame() {
//Find all connections and get the conn value
    dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true ORDER BY ShipHealth asc;")
    var index = 0;
    for dbConnections.Next() {
        var newDevice ConnectedDevices
        _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
        conn := GetConnection(newDevice)
        if index == 0{
            mapD := map[string]interface{}{
                "id":"gameOver",
                "endResult":"lost",
            }
            mapB, _ := json.Marshal(mapD)
            _, _ = conn.Write([]byte(mapB))
            index += 1
        } else{
            mapD := map[string]interface{}{
                "id":"gameOver",
                "endResult":"won",
            }
            mapB, _ := json.Marshal(mapD)
            _, _ = conn.Write([]byte(mapB))
        }

        

    }
}


func InitNewUser(userName string, ipAddress *net.TCPConn){
    //Init the user's game
        _, _ = databaseConnection.Query("INSERT INTO BoatState(IpAddress, UserName) Values (\""+ ipAddress.RemoteAddr().String() + "\", \""+ userName + "\");")
        getUsers(ipAddress)
}

func clearAccountsTable() {
    _, _ = databaseConnection.Query("TRUNCATE TABLE BoatState;")
}

func FireWeapon(wep string, ipAddress *net.TCPConn) {
	if wep == "1" {
		fmt.Println("Cannon Incoming")
		_, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth - 25 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
		_, _ = databaseConnection.Query("UPDATE BoatState SET NumberOfCannons = NumberOfCannons - 1, FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
		dbConnections, _ := databaseConnection.Query("SELECT ShipHealth from BoatState WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
		var healths [2]string
		var index = 0
		for dbConnections.Next() {
			var boatState BoatState
			_ = dbConnections.Scan(&boatState.shipHealth)
			healths[index] = boatState.shipHealth
			index += 1
		}
		opponentHealth, _ := strconv.Atoi(healths[0])

        if (ipAddress.RemoteAddr().String() == p1Address){
            p1JustShotCannon = true;
        } else {
            p2JustShotCannon = true;
        }


		if opponentHealth <= 0 {
			EndGame() //End the game
            return
		}
	} else if wep == "2" {
		fmt.Println("Torpedo Incoming")
		_, _ = databaseConnection.Query("UPDATE BoatState SET TorpedoDamage = 40 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
		_, _ = databaseConnection.Query("UPDATE BoatState SET TorpedoState = \"Cooldown\", FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
        AddPendingTorpedoAttack(ipAddress.RemoteAddr().String())
	} else if wep == "3" {
                fmt.Println("Mounted MG Incoming")
                 _, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth - 10 WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
	} else {
		fmt.Println("Invalid Weapon Type")
	}
    checkIfBothPlayersAreFinished()
}

func RepairShip(ipAddress *net.TCPConn) {
	dbConnections, _ := databaseConnection.Query("SELECT ShipHealth from BoatState WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
                var healths [2]string
                var index = 0
                for dbConnections.Next() {
                        var boatState BoatState
                        _ = dbConnections.Scan(&boatState.shipHealth)
                        healths[index] = boatState.shipHealth
                        index += 1
                }
                n, _ := strconv.Atoi(healths[0])
                //fmt.Println(n)
                if n > 75 {
                         _, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = 100, FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
                } else {
			_, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth + 25, FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
		}
	checkIfBothPlayersAreFinished()
}

func HackRadar(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("UPDATE BoatState SET FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    _, _ = databaseConnection.Query("UPDATE BoatState SET RadarState = \"Hacked\" WHERE IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
    checkIfBothPlayersAreFinished()
}

func FixRadar(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("UPDATE BoatState SET FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    _, _ = databaseConnection.Query("UPDATE BoatState SET RadarState = \"Enabled\"WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1;")
    checkIfBothPlayersAreFinished()
}

func ChangePosition(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("UPDATE BoatState SET FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
	_, _ = databaseConnection.Query("UPDATE BoatState SET navigationPosition = " + strconv.Itoa(rand.Intn(100)) + " WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    checkIfBothPlayersAreFinished()
}

func PlayerLostMiniGame(ipAddress *net.TCPConn) {
    _, _ = databaseConnection.Query("UPDATE BoatState SET FinishedMiniGame = 1 WHERE IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    checkIfBothPlayersAreFinished()
}

func checkIfBothPlayersAreFinished(){
    addressesResult, _ := databaseConnection.Query("SELECT IpAddress from BoatState where FinishedMiniGame = 1;")

    var ipAddresses [2]string
    var index = 0;
    for addressesResult.Next(){
        var status PlayersFinished
        _ = addressesResult.Scan(&status.IpAddress)
        ipAddresses[index] = status.IpAddress

        index += 1
    }

    if index == 2{
        fmt.Println("Both players are finsied!")
	_, _ = databaseConnection.Query("UPDATE BoatState SET NumberOfCannons = NumberOfCannons + 1 WHERE IpAddress = \""+ipAddresses[0]+"\";")
        _, _ = databaseConnection.Query("UPDATE BoatState SET NumberOfCannons = NumberOfCannons + 1 WHERE IpAddress = \""+ipAddresses[1]+"\";")
	RoundEndedCheckForTorpedo()
        //Send the boatState to each player, reset FinishedMiniGame
        _, _ = databaseConnection.Query("UPDATE BoatState SET FinishedMiniGame = 0 WHERE IpAddress = \""+ipAddresses[0]+"\" OR IpAddress = \""+ipAddresses[1]+"\";")

        //Find all connections and get the conn value
        dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true;")
        for dbConnections.Next() {
            var newDevice ConnectedDevices
            _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
            conn := GetConnection(newDevice)

            mapD := map[string]interface{}{
                "id":"continueToInGameScreen",
            }
            mapB, _ := json.Marshal(mapD)
            _, _ = conn.Write([]byte(mapB))

        }

        SendAnimationDataToHardware()
    }

}

func StartGame(){
    dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true;")
    devices := make([]ConnectedDevices, 0)

    for dbConnections.Next() {
        var newDevice ConnectedDevices
        _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
        devices = append(devices, newDevice)
    }

    mapD := map[string]interface{}{
        "id": "startGame",
        "result": "true",
    }
    mapB, _ := json.Marshal(mapD)
    SendMessageToDevices(mapB, devices)

    SendAnimationDataToHardware()
}

func getUsers(sentAddress *net.TCPConn){
    dbConnections, _ := databaseConnection.Query("Select IpAddress ipAddress, UserName userName from BoatState where GameActive = true ORDER BY ipAddress = \"" + sentAddress.RemoteAddr().String() + "\" DESC;")
    devices := make([]ConnectedDevices, 0)
    userNames := make([]string, 0)

    for dbConnections.Next() {
        var newDevice ConnectedDevices
        _ = dbConnections.Scan(&newDevice.ipAddress, &newDevice.userName)
        devices = append(devices, newDevice)
        userNames = append(userNames, newDevice.userName)
    }

    var invert = false

    for _, device := range devices {
        conn := GetConnection(device)
        if invert && len(userNames) == 2{
            firstName := userNames[0]
            secondName := userNames[1]
            userNames[0] = secondName
            userNames[1] = firstName

        }
        invert = true

        mapD := map[string]interface{}{
            "id": "connectedDevices",
            "userNames": userNames,
        }
        mapB, _ := json.Marshal(mapD)
        _, _ = conn.Write([]byte(mapB))
    }

    SendPlayerNamesToHardWare()
    
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
    boatStateResult := databaseConnection.QueryRow("SELECT navigationPosition, RadarState radarState, ShipHealth shipHealth, NumberOfCannons numberOfCannons, TorpedoState torpedoState, opponentBoat.opponentHealth, opponentBoat.opponentUserName, UserName userName from BoatState join (select ShipHealth opponentHealth, UserName opponentUserName from BoatState where IpAddress != \""+ ipAddress.RemoteAddr().String() + "\" LIMIT 1) as opponentBoat where IpAddress = \""+ ipAddress.RemoteAddr().String() + "\";")
    
    var boatState BoatState
    _ = boatStateResult.Scan(&boatState.navigationPosition, &boatState.radarState, &boatState.shipHealth, &boatState.numberOfCannons, &boatState.torpedoState, &boatState.opponentHealth, &boatState.opponentUserName, &boatState.userName)
    cannonState := "Enabled"
    if boatState.numberOfCannons == "0"{
        cannonState = "Reloading"
    }
    mapD := map[string]interface{}{
        "id":"boatState",
        "userName": map[string]interface{}{
            "p1UserName": boatState.userName,
            "p2UserName": boatState.opponentUserName,
        },
        "boatHealth": map[string]interface{}{
            "p1Health": boatState.shipHealth,
            "p2Health": boatState.opponentHealth,
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

func SendAnimationDataToHardware(){
    //Find all connections and get the conn value
    dbConnections, _ := databaseConnection.Query("Select RadarState radarState, ShipHealth shipHealth from BoatState where GameActive = true ORDER BY id;")
    
    var radarStates [2]string
    var healths [2]string
    var index = 0

    for dbConnections.Next() {
        var boatState BoatState
        _ = dbConnections.Scan(&boatState.radarState, &boatState.shipHealth)
        radarStates[index] = boatState.radarState
        healths[index] = boatState.shipHealth
        index += 1
    }


    mapD := map[string]interface{}{
        "id": "FinishedMiniGame",
        "player1": map[string]interface{}{
            "health": healths[0],
            "radarState": radarStates[0],
            "shotCanon": p1JustShotCannon,
        },
        "player2": map[string]interface{}{
            "health": healths[1],
            "radarState": radarStates[1],
            "shotCanon": p2JustShotCannon,
        },
    }

    p1JustShotCannon = false;
    p2JustShotCannon = false;

    mapB, _ := json.Marshal(mapD)
    SendMessageToHardware(mapB)
}

func SendPlayerNamesToHardWare(){
    //Find all connections and get the conn value
    dbConnections, _ := databaseConnection.Query("Select UserName userName, IpAddress ipAddress from BoatState where GameActive = true ORDER BY id;")
    
    var userNames [2]string
    userNames[0] = ""
    userNames[1] = ""

    var addresses [2]string
    addresses[0] = ""
    addresses[1] = ""

    var index = 0
    for dbConnections.Next() {
        var boatState BoatState
        _ = dbConnections.Scan(&boatState.userName, &boatState.ipAddress)
        userNames[index] = boatState.userName
        addresses[index] = boatState.ipAddress
        index += 1
    }

    mapD := map[string]interface{}{
        "id": "displayUserNames",
        "player1": userNames[0],
        "player2": userNames[1],
    }

    p1Address = addresses[0]
    p2Address = addresses[1]


    mapB, _ := json.Marshal(mapD)
    SendMessageToHardware(mapB)
}

func SendIncomingTorpedoToHardware(addressOfAttacker string){
    var playerWhoIsGettingAttacked = "p1"
    if (addressOfAttacker == p1Address){
        playerWhoIsGettingAttacked = "p2"
    }



    mapD := map[string]interface{}{
        "id": "incomingTorpedo",
        "headingTowards": playerWhoIsGettingAttacked,
    }
    mapB, _ := json.Marshal(mapD)
    SendMessageToHardware(mapB)

}




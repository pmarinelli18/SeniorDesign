package connection

import (
	"fmt"
	"strconv"
)

//Number of rounds there have been in the game. Used to track progress over time.
var currentRoundNumber = 1
var pendingTorpedoAttacks []TorpedoHit

func RoundEndedCheckForTorpedo(){
	//Need to check for delayed damage last round
	//Need to check if player moved out of the way

	//Go through pending torpedos and see if there is one that needs to be sent 

	for _, attack := range pendingTorpedoAttacks {
        if attack.roundToMakeAttack == currentRoundNumber {
            //Check if the opponent has moved
            boatStateResult := databaseConnection.QueryRow("SELECT navigationPosition from BoatState where IpAddress != \""+ attack.playerWhoShotAddress + "\";")
    		var opponentBoatState BoatState
    		_ = boatStateResult.Scan(&opponentBoatState.navigationPosition)
    		
    		fmt.Println(attack.currentpositionOfOpponent)
    		fmt.Println(opponentBoatState.navigationPosition)

    		if (attack.currentpositionOfOpponent == opponentBoatState.navigationPosition){
    			fmt.Println("Actually do damage now!")
			_, _ = databaseConnection.Query("UPDATE BoatState SET shipHealth = shipHealth - 40 WHERE IpAddress != \""+ attack.playerWhoShotAddress + "\" LIMIT 1;")
			dbConnections, _ := databaseConnection.Query("SELECT ShipHealth from BoatState WHERE IpAddress != \""+ attack.playerWhoShotAddress + "\" LIMIT 1;")
        	var healths [2]string
        	var index = 0
        	for dbConnections.Next() {
                	var boatState BoatState
                	_ = dbConnections.Scan(&boatState.shipHealth)
                	healths[index] = boatState.shipHealth
                	index += 1
        	}
        	opponentHealth, _ := strconv.Atoi(healths[0])
			if opponentHealth <= 0 {
                EndGame() //Needs to be fixed
    			return
        	}
		} else{
    			fmt.Println("Opponent dodged attack!")
    		}
        }

        //It has been one turn since the torpedo hit or missed the other boat. Reenable it
        if attack.roundToMakeAttack + 1 == currentRoundNumber {
        	fmt.Println("Reenable the torpedo!")
       		_, _ = databaseConnection.Query("UPDATE BoatState SET TorpedoState = \"Standby\", FinishedMiniGame = 1 WHERE IpAddress = \""+ attack.playerWhoShotAddress + "\";")
		}
    }

	fmt.Print("Just finished round ")
	fmt.Println(currentRoundNumber)
	currentRoundNumber += 1
}

func AddPendingTorpedoAttack(addressOfAttacker string){
	//Change 1 to allow users more than 1 turn to navigate away

	//If the opponent has their radar enbaled, let them know they have an incoming torpedo!
	//get opponents position
	boatStateResult := databaseConnection.QueryRow("SELECT navigationPosition, RadarState radarState from BoatState where IpAddress != \""+ addressOfAttacker + "\";")
    var opponentBoatState BoatState
    _ = boatStateResult.Scan(&opponentBoatState.navigationPosition, &opponentBoatState.radarState)

	pendingTorpedoAttacks = append(pendingTorpedoAttacks, TorpedoHit{currentRoundNumber + 1, addressOfAttacker, opponentBoatState.navigationPosition})

	//Alert opponent that they are getting attacked.
	if (opponentBoatState.radarState == "Enabled"){
		SendIncomingTorpedoToHardware(addressOfAttacker)
	}
}


type TorpedoHit struct {
    roundToMakeAttack int
    playerWhoShotAddress  string
    currentpositionOfOpponent string
}

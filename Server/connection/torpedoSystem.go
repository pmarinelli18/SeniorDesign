package connection

import (
	"fmt"
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
	boatStateResult := databaseConnection.QueryRow("SELECT navigationPosition from BoatState where IpAddress != \""+ addressOfAttacker + "\";")
    var opponentBoatState BoatState
    _ = boatStateResult.Scan(&opponentBoatState.navigationPosition)

	pendingTorpedoAttacks = append(pendingTorpedoAttacks, TorpedoHit{currentRoundNumber + 1, addressOfAttacker, opponentBoatState.navigationPosition})
}


type TorpedoHit struct {
    roundToMakeAttack int
    playerWhoShotAddress  string
    currentpositionOfOpponent string
}
package connection
import (
	"fmt"
	"database/sql"
    _ "github.com/go-sql-driver/mysql"
)

var databaseConnection *sql.DB 


type Tag struct {
    userName   string    `json:"userName"`
}

func MakeDatabaseConnection(){
	var err error

	 databaseConnection, err = sql.Open("mysql", "root:password@tcp(localhost)/BattleShip")

	// if there is an error opening the connection, handle it
    if err != nil {
        panic(err.Error())
    } else{
    	fmt.Println("Successful connected to database!")
    }
}

func CheckIfValidLogin(userName string, password string){
	results, err := databaseConnection.Query("SELECT userName FROM Accounts")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    for results.Next() {
        var tag Tag
        // for each row, scan the result into our tag composite object
        err = results.Scan(&tag.userName)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }
                // and then print out the tag's Name attribute
        fmt.Println(tag.userName)
    }
}
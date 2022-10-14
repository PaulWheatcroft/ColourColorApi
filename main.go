package main

import (
	"database/sql"
	"net/http"
	"fmt"
	"os"
	"log"
	"sort"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/go-sql-driver/mysql"
)

type Colour struct {
	ID			int		`json:"id"`
	ColourName	string	`json:"colourName"`
	HexValue	string	`json:"hexValue"`
}

var db *sql.DB

func setDataBaseConnection() {
	dsn := mysql.Config{
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PASSWORD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DB_ADDRESS"),
			DBName:               os.Getenv("DB_NAME"),
			AllowNativePasswords: true,
	}

	// Get a database handle
	var err error
	db, err = sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func getColours(context *gin.Context) {
	var colours []Colour
	// Data Source Name Properties
	
	setDataBaseConnection()
	//Get all prime ministers
	queryResults, err := db.Query("select * from colours")
    if err != nil {
        fmt.Printf("The error is: %v", err)
		return
    }
	for queryResults.Next() {
        var colour Colour
        if err := queryResults.Scan(&colour.ID, &colour.ColourName, &colour.HexValue); err != nil {
            fmt.Printf("Error: %v", err)
			return
        }
        colours = append(colours, colour)
		
    }

	sort.SliceStable(colours, func(i, j int) bool {
		return colours[i].ID > colours[j].ID
	})

	//Put the results in to JSON and pass to the context
	fmt.Println(colours)
	context.IndentedJSON(http.StatusOK, colours)
	fmt.Println("The last line is", context.Request)

}

func addColour(context *gin.Context) {
	var newColour Colour
	setDataBaseConnection()
	fmt.Println("I ran addColour")
	if err := context.BindJSON(&newColour); err != nil {
		return
	}

	context.IndentedJSON(http.StatusCreated, newColour.ColourName)

	fmt.Println(newColour.ColourName, newColour.HexValue)

	db.Exec("INSERT INTO colours (colour_name, hex_value) VALUES (?, ?)", newColour.ColourName, newColour.HexValue)
	
}

func deleteColour(context *gin.Context) {
	var deleteColour Colour
	setDataBaseConnection()
	fmt.Println("I ran deleteColour")
	if err := context.BindJSON(&deleteColour); err != nil {
		return
	}

	context.IndentedJSON(http.StatusCreated, deleteColour.ColourName)

	fmt.Println(deleteColour.ColourName, deleteColour.HexValue, deleteColour.ID)

	db.Exec("delete from colours where (id)=(?)", deleteColour.ID)
	
}

func editColour(context *gin.Context) {
	var editColour Colour
	setDataBaseConnection()
	fmt.Println("I ran editColour")
	if err := context.BindJSON(&editColour); err != nil {
		return
	}

	context.IndentedJSON(http.StatusCreated, editColour.ColourName)

	fmt.Println(editColour.ColourName, editColour.HexValue, editColour.ID)

	db.Exec("UPDATE colours SET `colour_name` = ?, `hex_value` = ? WHERE `id` = ?", editColour.ColourName, editColour.HexValue, editColour.ID)
	
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
    }))
	fmt.Printf("I am running ColourColorApi")
	router.GET("/colours", getColours)
	router.POST("/add", addColour)
	router.DELETE("/delete", deleteColour)
	router.PATCH("/edit", editColour)
	router.Run("127.0.0.1:1212")
}
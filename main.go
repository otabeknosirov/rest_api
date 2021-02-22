package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	// "gopkg.in/gin-gonic/gin.v1"
)

type DBClient struct {
	db *sql.DB
}

// Model the record struct
type StationResource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	OpeningTime string `json:"opening_time"`
	ClosingTime string `json:"closing_time"`
}

// CreateStation handles the POST
func (DB *DBClient) CreateStation(c *gin.Context) {
	var station StationResource
	// Parse the body into our resrource
	if err := c.BindJSON(&station); err == nil {
		// Format Time to Go time format
		statement, err := DB.db.Prepare("INSERT INTO station(name, openingtime, closingtime) VALUES ($1, $2, $3)")
		if err != nil {
			fmt.Println("err = ", err)
			return
		}
		fmt.Println(station)
		result, _ := statement.Exec(station.Name,
			station.OpeningTime, station.ClosingTime)
		if err == nil {
			newID, _ := result.LastInsertId()
			station.ID = int(newID)
			c.JSON(http.StatusOK, gin.H{
				"result": station,
			})
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// RemoveStation handles the removing of resource
func (DB DBClient) RemoveStation(c *gin.Context) {
	id := c.Param("station-id")
	statement, _ := DB.db.Prepare("delete from station where id=?")
	_, err := statement.Exec(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "")
	}
}

// GetStation returns the station detail
func (DB DBClient) GetStation(c *gin.Context) {
	var station StationResource
	id := c.Param("station_id")
	err := DB.db.QueryRow("select ID, Name, CAST(OpeningTime as CHAR),CAST(ClosingTime as CHAR) from station where id=?", id).
		Scan(&station.ID, &station.Name, &station.OpeningTime, &station.ClosingTime)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"result": station,
		})
	}
}

func main() {
	// db, err := models.InitDB()
	// if err != nil {
	// 	log.Println(db)
	// }
	log.Println("------------")
	db, err := InitDB()
	log.Println("------------", db)
	if err != nil {
		log.Println(db)
	}
	dbclient := &DBClient{db: db}

	if err != nil {
		log.Println("------------", err)
		panic(err)
	}
	defer db.Close()

	//CREATE new router
	req := gin.Default()
	// Add routes to REST verbs
	req.GET("/v1/stations/:station_id", dbclient.GetStation)
	req.POST("/v1/stations", dbclient.CreateStation)
	req.DELETE("/v1/stations/:station_id", dbclient.RemoveStation)
	req.Run(":9090")
	//Attach an elegant path with handler
}
func InitDB() (*sql.DB, error) {
	var err error
	db, err := sql.Open("postgres",
		"postgres://userdb:cr7@localhost/userdb?sslmode=disable")
	if err != nil {
		return nil, err
	} else {
		stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS station(ID SERIAL PRIMARY KEY, Name VARCHAR(64) NULL, OpeniningTime TIME NULL, ClosingTime TIME NULL);`)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		result, err := stmt.Exec()
		log.Println(result)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return db, err
	}
}

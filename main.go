package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"encoding/json"
	"io/ioutil"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Id 	 int 	`json:"Id"`
	Name string `json:"Name"`
	Age  int 	`json:"Age"`
	Uri  string `json:Uri`
}

func main() {
	log.SetReportCaller(true)

    router := gin.Default()
    
	router.GET("/health", get_health)
	router.GET("/items", get_items)
    router.Run(":8080")
}

func get_health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Healthy")
}

func get_items(c *gin.Context) {
	jsonFile, err := os.Open(os.Getenv("ITEMS_PATH"))

	if err != nil {
		log.Panicf("Fatal error reading items: %v \n", err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var items []Item

	if err := json.Unmarshal(byteValue, &items); err != nil {
		log.Panicf("Fatar error unmarshaling data: %v \n", err)
	}

	c.IndentedJSON(http.StatusOK, items)
}
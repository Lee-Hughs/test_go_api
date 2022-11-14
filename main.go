package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Items struct {
	Items []Item `json:"items"`
}

type Item struct {
	Id   int    `json:"Id"`
	Name string `json:"Name"`
	Age  int    `json:"Age"`
	Uri  string `json:Uri`
}

const (
	host     = "localhost"
	port     = 5000
	user     = "postgres"
	password = "Guitar938"
	db_name  = "local_db"
)

func main() {
	log.SetReportCaller(true)

	viper.SetDefault("PORT", ":8080")

	viper.SetConfigName("")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetEnvPrefix("")
			viper.AutomaticEnv()
		} else {
			log.Panicf("Fatal error config file: %v \n", err)
		}
	}

	if viper.GetString("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.GET("/health", get_health)
	router.GET("/items", get_items)
	router.GET("/item/:id", get_item)
	router.Run(viper.GetString("PORT"))
}

func get_health(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Healthy")
}

func get_items(c *gin.Context) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db_name)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT item_id, name FROM items LIMIT $1", 10)
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			// handle this error
			panic(err)
		}
		items = append(items, fmt.Sprintf("%d: %s", id, name))
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	c.IndentedJSON(http.StatusOK, items)
}

func get_item(c *gin.Context) {
	item_id := c.Param("id")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, db_name)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	row := db.QueryRow("SELECT item_id, name FROM items WHERE item_id = $1 LIMIT 1", item_id)
	var id int
	var name string
	switch err := row.Scan(&id, &name); err {
	case sql.ErrNoRows:
		c.IndentedJSON(http.StatusNotFound, "Item not found")
	case nil:
		c.IndentedJSON(http.StatusOK, name)
	default:
		panic(err)
	}
}

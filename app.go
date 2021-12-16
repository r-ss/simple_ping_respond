package main

import (
	"time"
	"github.com/gin-gonic/gin"
	"encoding/json"
	"fmt"
	"io/ioutil"
    "os"
	"math/rand"
)

type AppConfig struct {
    Resource   string    `json:"resource"`
    Port   	   string    `json:"port"`
	Allowed    string    `json:"allowed"`
}


func randomStringSequence(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
    b := make([]rune, n)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}


func appShutdown() {
    fmt.Println("Shutting down...")
    os.Exit(0)
}

func main() {
	// gin.SetMode(gin.ReleaseMode)

	startTime := time.Now()

	jsonFile, err := os.Open("config.json")
	if err != nil {
		fmt.Println(err)
		appShutdown()
	}
    defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config AppConfig

	json.Unmarshal(byteValue, &config)

	router := gin.Default()
	router.SetTrustedProxies([]string{config.Allowed})
	

	router.GET("/probe", func(c *gin.Context) {

		// fmt.Printf("ClientIP: %s\n", c.ClientIP())
		
		uptime := time.Since(startTime)
		verbtime := int64(uptime) / int64(time.Second)

		c.JSON(200, gin.H{
			"resource": config.Resource,
			"random": randomStringSequence(8),
			"uptime": verbtime,
		})
	})

	var result map[string]interface{}
    json.Unmarshal([]byte(byteValue), &result)

	fmt.Println("config.Resource: " + config.Resource)
	fmt.Println("config.Port: " + config.Port)
	fmt.Println("config.Allowed: " + config.Allowed)

	// port := strconv.Itoa(config.Port)
	router.Run(":" + config.Port) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
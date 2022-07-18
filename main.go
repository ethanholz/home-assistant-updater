package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

// Posts update to Home Assistant Entity
func postUpdate() {
	endpoint := "http://localhost:8123/api/services/input_boolean/turn_on"
	// Set using an environment variables
	auth_token := os.Getenv("HASS_TOKEN")
	// Maybe setup to restart one day?
	object := map[string]string{"entity_id": "input_boolean.update"}
	jsonValue, err := json.Marshal(object)
	if err != nil {
		log.Fatal(err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatal(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth_token)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Posted request to HA")
	if resp.Status != "200 OK" {
		log.Fatal("Failed to make request to HA")
		return
	}

}

func backgroundPull() {
	imageName := "ghcr.io/home-assistant/home-assistant:stable"
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("Imported context")
	// Pull updated image
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	buf := new(strings.Builder)
	defer out.Close()
	_, err = io.Copy(buf, out)
	if err != nil {
		log.Fatal(err)
		return
	}

	match, _ := regexp.MatchString("Status: Downloaded newer image", buf.String())
	if match {
		log.Println("Pulled new image")
		go postUpdate()
	}
}

func pullContainer(c *gin.Context) {
	// Accepts and run
	c.JSON(http.StatusAccepted, nil)
	go backgroundPull()
}

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.POST("/update", pullContainer)
	router.Run("localhost:8081")
}

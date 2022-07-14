package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

// Posts update to Home Assistant Entity
func postUpdate() {
	endpoint := "http://localhost:8123/api/services/input_boolean/turn_on"
	auth_token := os.Getenv("HASS_TOKEN")
	object := map[string]string{"entity_id": "input_boolean.update"}
	jsonValue, err := json.Marshal(object)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+auth_token)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.Status != "200 OK" {
		panic("Failed to make request")
	}

}

func backgroundPull() {
	imageName := "ghcr.io/home-assistant/home-assistant:stable"
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	fmt.Println("Imported context")
	// Pull updated image
	out, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Pulled image")
	defer out.Close()
	io.Copy(os.Stdout, out)
	go postUpdate()
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

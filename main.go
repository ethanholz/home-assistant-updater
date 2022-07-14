package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

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

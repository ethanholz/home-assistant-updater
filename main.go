package main

import (
    "net/http"
    "context"
    "fmt"
    "os"
    "io"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/gin-gonic/gin"
)

func pullContainer(c *gin.Context) {
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
    c.JSON(http.StatusCreated, nil)
}

func main() {
    router := gin.Default()
    router.POST("/update", pullContainer)
    router.Run("localhost:8081")
}



package main

import (
    "net/http"
    "context"
    "fmt"
    "os"
    "io"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types/container"
    "github.com/gin-gonic/gin"
)

func pullContainer(c *gin.Context) {
    imageName := "ghcr.io/home-assistant/home-assistant:stable"
    containerName := "homeassistant"
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
    // Get a list of containers
    containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
    if err != nil {
        panic(err)
    }

    //Find HA container
    var haContainer types.Container
    for _, container := range containers {
        if container.Names[0][1:] == containerName {
            haContainer = container
            break
        }
    }
    if haContainer.ID == "" {
        fmt.Println("Does not exist")
    } else {
        // Stop the HA container
        if err := cli.ContainerStop(ctx, haContainer.ID, nil); err != nil {
            panic(err)
        }
        if err := cli.ContainerRemove(ctx, haContainer.ID, types.ContainerRemoveOptions{}); err != nil {
            panic(err)
        }
    }

    resp, err := cli.ContainerCreate(ctx, &container.Config{
        Image: imageName,
    }, nil, nil, nil, containerName)
    if err != nil {
        panic(err)
    }
    if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil{
        panic(err)
    }

    c.JSON(http.StatusCreated, nil)
}

func main() {
    router := gin.Default()
    router.POST("/update", pullContainer)
    router.Run("localhost:8080")
}



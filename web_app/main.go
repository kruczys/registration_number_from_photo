package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    
    router.Static("/", "./static")
    router.MaxMultipartMemory = 8 << 20

    fmt.Println("Server is running at port 8080")
    router.Run(":8080")
}

package main

import (
	"github.com/acouvreur/sablier/cmd"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	cmd.Execute()
}

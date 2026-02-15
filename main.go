
package main

import (
	"log"

	"github.com/JerryJeager/r3sonance-backend/cmd"
	"github.com/JerryJeager/r3sonance-backend/config"
)

func init() {
	config.LoadEnv()
	config.ConnectToDB()
}

func main() {
	log.Println("Starting Server")


	cmd.ExecuteApiRoutes()
}

	
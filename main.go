package main

import (
	"gmeroblog/model"
	"gmeroblog/routes"
)

func main() {
	model.InitDb()
	routes.InitRouter()
}

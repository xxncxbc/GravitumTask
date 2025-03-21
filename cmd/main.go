package cmd

import (
	"GravitumTask/configs"
	"fmt"
)

func main() {
	//loading configs
	conf := configs.LoadConfig()
	fmt.Println(conf)
	//setting up database

}

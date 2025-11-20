package main

import (
	"block_chain/app"
	"block_chain/config"
	"flag"
	"fmt"
)

var (
	configFlag = flag.String("environment", "./environment.toml", "environment toml file not found")
	difficulty = flag.Int("difficulty", 17, "difficulty err")
)

func main() {
	flag.Parse()

	c := config.NewConfig(*configFlag)

	app.NewApp(c, int64(*difficulty))

	fmt.Println(c.Info.Version)

	fmt.Println("test")
}

package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 这就像是 var a = new A();
	var config = &Config{}
	c, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	data,err := json.Marshal(c)
	fmt.Print(string(data))

	for _, s := range c.Repos {
		go sync(s, c.Config)
	}
}

func sync(config *syncConfig, globalConfig *globalConfig) {
}

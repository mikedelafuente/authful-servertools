package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

var configOnce sync.Once
var configInstance *AuthfulConfig

func GetConfig() *AuthfulConfig {
	configOnce.Do(func() {
		var err error

		configInstance, err = getConfigInstanceFromEnvironment()

		if err != nil {
			log.Println(err)
			panic(err)
		}
	})

	return configInstance
}

func getConfigInstanceFromEnvironment() (*AuthfulConfig, error) {
	log.Println("Loading config from environment")

	var myConfig *AuthfulConfig = &AuthfulConfig{}

	isDebug, err := strconv.ParseBool(os.Getenv("IS_DEBUG"))
	if err != nil {
		log.Println(err)
		myConfig.IsDebug = true
		return myConfig, err
	}
	myConfig.IsDebug = isDebug

	return myConfig, nil
}

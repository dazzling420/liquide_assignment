package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func GetConfig() *Config {
	//Get the environment to run flag
	envFlag := flag.String("env", "test", "Environment option default set to test")
	flag.Parse()
	*envFlag = strings.ToLower(*envFlag)

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("dir ::" + dir)

	configPath, err := filepath.Abs(dir + "/" + BaseConfigPath + *envFlag + ConfigFileSuffix)
	if err != nil {
		fmt.Println("Err ::", err.Error())
	}

	fmt.Println("config Path ::", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		return nil
	}

	fmt.Println("data -> ::", string(data))

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %s\n", err)
		return nil
	}

	xyz, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		fmt.Println("json error ", err.Error())
	}

	fmt.Println("Config xyz -> ::", string(xyz))
	fmt.Println("Config -> ::", config)
	return &config
}

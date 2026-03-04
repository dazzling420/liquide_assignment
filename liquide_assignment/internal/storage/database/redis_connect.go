package database

import (
	"fmt"
	"liquide_assignment/internal/config"
	"strings"

	"github.com/redis/rueidis"
)

func RedisConnect(conf config.Redis) (rueidis.Client, error) {
	connStringUrl := generateRedisConnectionString(conf.Host, conf.Port)

	c, err := rueidis.NewClient(rueidis.ClientOption{
		Username:     conf.Username,
		InitAddress:  connStringUrl,
		Password:     conf.Password,
		DisableCache: true,
	})

	if err != nil {
		panic(err)
	}

	return c, nil
}

func generateRedisConnectionString(url string, port int) []string {
	hosts := strings.Split(url, ",")
	var addresses []string
	for _, currentHost := range hosts {
		currentHost = strings.Trim(currentHost, " ")
		a := fmt.Sprintf("%s:%d", currentHost, port)
		addresses = append(addresses, a)
	}
	return addresses
}

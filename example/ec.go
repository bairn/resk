package main

import (
	"github.com/tietang/go-eureka-client/eureka"
	"time"
)

func main() {
	cfg := eureka.Config{
		DialTimeout:time.Second * 10,
	}

	client := eureka.NewClientByConfig([]string{
		"http://eureka.springcloud.cn/eureka",
	}, cfg)

	appName := "Go-Example"
	instance := eureka.NewInstanceInfo(
		"test.com", appName,
		"127.0.0.2",
		8080, 30,
		false)

	client.RegisterInstance(appName, instance)
	client.Start()
	c := make(chan int, 1)
	<-c
}

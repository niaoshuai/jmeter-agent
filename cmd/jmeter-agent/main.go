package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/niaoshuai/jmeter-agent/pkg/jmeter"
)

const HttpPort = "8080"

// 启动 Jmeter
func main() {
	ctx := context.Background()

	agent := jmeter.Agent{Version: "5.4.3", InstallPath: "/Users/coding/go/src/github.com/niaoshuai/jmeter-agent", HttpPort: HttpPort}
	agent.Register(ctx)

	r := gin.Default()
	r.GET("/start", func(c *gin.Context) {
		// 没有启动
		if _, err := agent.GetJmeterPID(); err != nil {
			agent.StartJmeterServer()
			agent.Status = jmeter.StatusStart
			c.JSON(200, gin.H{
				"message": "ok",
			})
			return
		}
		c.JSON(400, gin.H{
			"message": "已经启动了",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		if _, err := agent.GetJmeterPID(); err != nil {
			c.JSON(400, gin.H{
				"message": "没有正在进行的进程",
			})
			return
		}

		err := agent.StopJmeterServer()
		agent.Status = jmeter.StatusStop
		if err != nil {
			c.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": "ok",
		})

	})
	r.Run(fmt.Sprintf("%s:%s", agent.Ip, agent.HttpPort)) // listen and serve on 0.0.0.0:8080
}

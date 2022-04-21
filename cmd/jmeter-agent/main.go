package main

import (
	"github.com/gin-gonic/gin"
	"github.com/niaoshuai/jmeter-agent/pkg/jmeter"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const (
	JmeterAgentStatusStart jmeter.JmeterAgentStatus = "start"
	JmeterAgentStatusStop  jmeter.JmeterAgentStatus = "stop"
)

// 启动 Jmeter
func main() {
	agent := jmeter.JmeterAgent{Version: "5.1", InstallPath: "/usr/local"}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			// 获取进程号
			pid, err := agent.GetJmeterServerPid()
			if err != nil {
				log.Println(err)
				agent.Status = JmeterAgentStatusStop
				_ = os.Remove("jmeter.lock")
			} else {
				agent.Status = JmeterAgentStatusStart
				ioutil.WriteFile("jmeter.lock", pid, 0666)
			}
		}
	}()

	r := gin.Default()
	r.GET("/install", func(c *gin.Context) {
		agent.InstallJmeter()
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	r.GET("/start", func(c *gin.Context) {
		agent.StartJmeterServer()
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	r.GET("/stop", func(c *gin.Context) {
		agent.StopJmeterServer()
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

package task

import (
	"context"
	"fmt"
	"github.com/niaoshuai/jmeter-agent/pkg/redis"
	"log"
	"os/exec"
	"strings"
)

// 发起压测任务
func StartJmeterTask(jmeterFile string, ctx context.Context, installPath string, jmeterVersion string) {
	agentList, err := redis.GetAgentService(ctx)
	cmd := exec.Command("./bin/jmeter",
		"-n",
		"-t", jmeterFile,
		"-R", strings.Join(agentList, ","),
		"-l",
		"report.jtl",
		"-Dserver.rmi.ssl.disable=true",
	)
	cmd.Dir = installPath + fmt.Sprintf("/apache-jmeter-%s", jmeterVersion)
	msg, err := cmd.Output()
	if err != nil { // 运行命令
		log.Println(err)
	}
	fmt.Println(string(msg))
}

func cancelJmeterTask() {

}

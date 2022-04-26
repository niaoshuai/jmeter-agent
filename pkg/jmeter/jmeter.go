package jmeter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

import (
	"github.com/niaoshuai/jmeter-agent/pkg/ip"
	"github.com/niaoshuai/jmeter-agent/pkg/redis"
)

const ServerPort = "2099"

type AgentStatus string

const (
	StatusStart AgentStatus = "start"
	StatusStop  AgentStatus = "stop"
)

type Agent struct {
	Version     string
	InstallPath string
	Status      AgentStatus
	Ip          string
	HttpPort    string
	ServerPort  string
}

func (agent *Agent) DownloadJmeter() error {
	var (
		jmeterDownload = `https://mirrors.tuna.tsinghua.edu.cn/apache/jmeter/binaries/apache-jmeter-%s.tgz`
	)
	jmeterDownloadUrl := fmt.Sprintf(jmeterDownload, agent.Version)
	res, err := http.Get(jmeterDownloadUrl)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("文件不存在")
	}
	//创建下载存放exe
	f, err := os.Create(agent.InstallPath + fmt.Sprintf("/apache-jmeter-%s.tgz", agent.Version))
	if err != nil {
		return err
	}
	io.Copy(f, res.Body)
	defer f.Close()
	return nil
}

// InstallJmeter 安装Jmeter
func (agent *Agent) InstallJmeter() error {
	//执行该路径下的exe并安装
	fileName := fmt.Sprintf("apache-jmeter-%s.tgz", agent.Version)
	cmd := exec.Command("tar", "-zxvf", fileName)
	cmd.Dir = agent.InstallPath
	err := cmd.Run()
	if err != nil {
		return err
	}
	log.Printf("安装完成！安装路径为: %s \n", agent.InstallPath+"/"+fileName)
	return nil
}

// StartJmeterServer 启动JmeterServer
func (agent *Agent) StartJmeterServer() {
	//java -server -XX:+HeapDumpOnOutOfMemoryError -Xms1g -Xmx1g -XX:MaxMetaspaceSize=256m  -Djava.security.egd=file:/dev/urandom -Xdock:name=JMeter -Xdock:icon=docs/images/jmeter_square.png -Dapple.laf.useScreenMenuBar=true -Dapple.eawt.quitStrategy=CLOSE_ALL_WINDOWS -Djava.rmi.server.hostname=172.100.80.35 -Dserver.rmi.ssl.disable=true -jar bin/ApacheJMeter.jar -s -j demo.log
	//nohup./bin/jmeter-server -Djava.rmi.server.hostname=172.100.80.35 -Dserver_port=${SERVER_PORT:-1099} -Dserver.rmi.ssl.disable=true > runoob.log 2>&1 &

	cmd := exec.Command("./bin/jmeter",
		"-s",
		fmt.Sprintf("-Djava.rmi.server.hostname=%s", agent.Ip),
		fmt.Sprintf("-Dserver_port=%s", ServerPort),
		"-Dserver.rmi.ssl.disable=true")
	cmd.Dir = agent.InstallPath + fmt.Sprintf("/apache-jmeter-%s", agent.Version)
	if err := cmd.Start(); err != nil { // 运行命令
		log.Println(err)
	}
	agent.ServerPort = ServerPort
}

func (agent *Agent) GetJmeterPID() ([]byte, error) {
	cmd := exec.Command("pgrep",
		"-f",
		"ApacheJMeter")
	cmd.Dir = agent.InstallPath + fmt.Sprintf("/apache-jmeter-%s", agent.Version)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(output) < 1 {
		log.Println("null pid")
		return nil, err
	}
	return output, nil
}

func (agent *Agent) StopJmeterServer() error {
	pid, err := agent.GetJmeterPID()
	if err != nil {
		log.Println(err)
		return err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf("kill -9 %s", string(pid)))
	cmd.Dir = agent.InstallPath + fmt.Sprintf("/apache-jmeter-%s", agent.Version)
	msg, err := cmd.Output()
	if err != nil { // 运行命令
		log.Println(err)
		return err
	}
	log.Println(msg)
	return nil
}

// Register 注册Agent 信息
func (agent *Agent) Register(ctx context.Context) {
	ip, err := ip.GetOutBoundIP()
	if err != nil {
		fmt.Println(err)
	}
	agent.Ip = ip
	_, err = agent.GetJmeterPID()
	if err != nil {
		agent.StartJmeterServer()
	}
	// 存储信息到Redis
	redis.AddAgentService(ctx, agent.Ip, agent.ServerPort)
}

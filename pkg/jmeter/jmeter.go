package jmeter

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type JmeterAgentStatus string

type JmeterAgent struct {
	Version     string
	InstallPath string
	Status      JmeterAgentStatus
}

func (agent *JmeterAgent) downloadJmeter() {
	ch := make(chan bool)
	go func() {
		var (
			jmeterDownload = `https://dlcdn.apache.org//jmeter/binaries/apache-jmeter-%s.tgz`
		)
		res, err := http.Get(fmt.Sprintf(jmeterDownload, agent.Version))
		if err != nil {
			<-ch
		}
		//errOwn.Err(err)
		//创建下载存放exe
		f, err := os.Create(agent.InstallPath + fmt.Sprintf("apache-jmeter-%s.tgz", agent.Version))
		//errOwn.Err(err)
		io.Copy(f, res.Body)
		defer f.Close()
		<-ch
	}()
	ch <- true
}

func (agent *JmeterAgent) GetJmeterServerPid() ([]byte, error) {
	cmd := exec.Command("bash",
		"-c",
		"jps -v |grep ApacheJMeter |awk '{print $1}'")
	//获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer stdout.Close()
	err = cmd.Start()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(opBytes) < 1 {
		return nil, errors.New("PID NOT FOUND")
	}

	return opBytes, nil
}

func GetOutBoundIP() (ip string, err error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	ip = strings.Split(localAddr.String(), ":")[0]
	return
}

// InstallJmeter 安装Jmeter
func (agent *JmeterAgent) InstallJmeter() {
	agent.downloadJmeter()
	//执行该路径下的exe并安装
	cmd := exec.Command("tar", "-zxvf", fmt.Sprintf("apache-jmeter-%s.tgz", agent.Version))
	cmd.Start()
	fmt.Printf("安装完成！安装路径为: %s \n")
}

// StartJmeterServer 启动JmeterServer
func (agent *JmeterAgent) StartJmeterServer() {
	//java -server -XX:+HeapDumpOnOutOfMemoryError -Xms1g -Xmx1g -XX:MaxMetaspaceSize=256m  -Djava.security.egd=file:/dev/urandom -Xdock:name=JMeter -Xdock:icon=docs/images/jmeter_square.png -Dapple.laf.useScreenMenuBar=true -Dapple.eawt.quitStrategy=CLOSE_ALL_WINDOWS -Djava.rmi.server.hostname=172.100.80.35 -Dserver.rmi.ssl.disable=true -jar bin/ApacheJMeter.jar -s -j demo.log
	//nohup./bin/jmeter-server -Djava.rmi.server.hostname=172.100.80.35 -Dserver_port=${SERVER_PORT:-1099} -Dserver.rmi.ssl.disable=true > runoob.log 2>&1 &
	ip, err := GetOutBoundIP()
	if err != nil {
		fmt.Println(err)
	}
	ch := make(chan bool)
	go func() {
		cmd := exec.Command("./bin/jmeter-server",
			fmt.Sprintf("-Djava.rmi.server.hostname=%s", ip),
			"-Dserver_port=${SERVER_PORT:-1099}",
			"-Dserver.rmi.ssl.disable=true")
		cmd.Dir = agent.InstallPath + fmt.Sprintf("apache-jmeter-%s/", agent.Version)
		if err := cmd.Start(); err != nil { // 运行命令
			log.Println(err)
		}

		<-ch
	}()
	ch <- true

}

// StopJmeterServer StopJmeter 停止Jmeter
func (agent *JmeterAgent) StopJmeterServer() {
	ch := make(chan bool)
	go func() {
		jmeterLock, _ := ioutil.ReadFile("jmeter.lock")
		cmd := exec.Command("kill", "-9", string(jmeterLock))
		if err := cmd.Start(); err != nil { // 运行命令
			log.Println(err)
		}
		<-ch
	}()
	ch <- true
}

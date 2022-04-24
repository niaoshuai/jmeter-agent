package tests

import (
	"github.com/niaoshuai/jmeter-agent/pkg/jmeter"
	"testing"
)

var (
	agent = jmeter.Agent{
		Version:     "5.4.3",
		InstallPath: "/Users/coding/go/src/github.com/niaoshuai/jmeter-agent",
	}
)

// 下载Jmeter
func TestJmeterDownload(t *testing.T) {
	err := agent.DownloadJmeter()
	if err != nil {
		t.Error(err)
	}
}

func TestJmeterInstall(t *testing.T) {
	err := agent.InstallJmeter()
	if err != nil {
		t.Error(err)
	}
}

func TestJmeterStart(t *testing.T) {
	agent.StartJmeterServer()
}

func TestGetJmeterPid(t *testing.T) {
	pid, err := agent.GetJmeterPID()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pid)
}

func TestStopJmeter(t *testing.T) {
	err := agent.StopJmeterServer()
	if err != nil {
		t.Fatal(err)
	}
}

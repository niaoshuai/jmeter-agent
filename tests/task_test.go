package tests

import (
	"context"
	"github.com/niaoshuai/jmeter-agent/pkg/task"
	"testing"
)

func TestStartJmeterTask(t *testing.T) {
	tmpFileName := "/Users/coding/Downloads/HTTP请求.jmx"
	task.StartJmeterTask(tmpFileName, context.Background(), agent.InstallPath, agent.Version)
}

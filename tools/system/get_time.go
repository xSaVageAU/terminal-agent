package system

import (
	"time"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

type GetTimeArgs struct{}

type GetTimeResult struct {
	Time string `json:"time"`
}

func getTime(_ tool.Context, _ GetTimeArgs) (GetTimeResult, error) {
	return GetTimeResult{
		Time: time.Now().Format(time.RFC3339),
	}, nil
}

func GetTimeTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "get_time",
		Description: "Returns the current system date and time.",
	}, getTime)
}

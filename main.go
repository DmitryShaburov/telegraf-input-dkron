package main

import (
	"os"
	"fmt"
	"time"
	"strings"
	"gopkg.in/resty.v1"
)

type JobResp struct {
	SuccessCount int       `json:"success_count"`
	ErrorCount   int       `json:"error_count"`
	LastSuccess  time.Time `json:"last_success"`
	LastError    time.Time `json:"last_error"`
	Name         string    `json:"name"`
}

type ExecutionJobResp struct {
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Success    bool      `json:"success"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Print("Host must be specified")
		os.Exit(1)
	}

	host := os.Args[1]
	host = strings.TrimRight(host, "/")

	urlJobs := fmt.Sprintf("%s/v1/jobs", host)
	respJobs, errJobs := resty.R().SetResult(&[]JobResp{}).Get(urlJobs)
	if errJobs != nil {
		fmt.Print(errJobs.Error())
	}

	resultJobs := respJobs.Result().(*[]JobResp)
	for _, job := range (*resultJobs) {
		state := "1"
		if job.LastSuccess.Before(job.LastError) {
			state = "0"
		}

		urlExecutions := fmt.Sprintf("%s/v1/jobs/%s/executions", host, job.Name)
		respExecutions, errExecutions := resty.R().SetResult([]ExecutionJobResp{}).Get(urlExecutions)
		if errExecutions != nil {
			fmt.Print(errExecutions.Error())
		}
		executions := respExecutions.Result().(*[]ExecutionJobResp)
		var duration int64
		for i := len(*executions)-1; i >= 0; i-- {
			if (*executions)[i].Success {
				duration = (*executions)[i].FinishedAt.Sub((*executions)[i].StartedAt).Nanoseconds() / 1000000
				break
			}
		}

		fmt.Printf("dkron,job=%v state=%v,success_count=%v,error_count=%v,last_duration=%v\n",
			job.Name, state, job.SuccessCount, job.ErrorCount, duration)
	}
}
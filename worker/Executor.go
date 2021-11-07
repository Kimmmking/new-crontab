package worker

import (
	"github.com/Kimmmking/crontab/common"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

// 执行任务
func (executor *Executor) ExecuteJob(jobExecuteInfo *common.JobExecuteInfo) {
	go func() {
		var (
			cmd    *exec.Cmd
			err    error
			result *common.JobExecuteResult
			output []byte
		)

		result = &common.JobExecuteResult{
			ExecuteInfo: jobExecuteInfo,
			Output:      make([]byte, 0),
		}

		// 记录任务开始时间
		result.StartTime = time.Now()

		// 执行shell命令
		cmd = exec.CommandContext(jobExecuteInfo.CancelCtx, "/bin/bash", "-c", jobExecuteInfo.Job.Command)

		// 执行并捕获输出
		output, err = cmd.CombinedOutput()

		// 记录任务结束时间
		result.EndTime = time.Now()
		result.Output = output
		result.Err = err

		// TODO 任务执行完成后，把执行的结果返回给 scheduler，scheduler 会从 executingTable 中删除掉执行记录

	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}

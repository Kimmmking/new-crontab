package master

type Scheduler interface {
}

var G_scheduler = (*Scheduler)(nil)

func InitScheduler() (err error) {
	return
}

// 1. 负载均衡

// 2. 分派任务的时机

package schedule

type Scheduler interface {
	Schedule(job *Job)
	// TODO 增加 Info 接口，查看调度器内部详情
	//Info() string
}

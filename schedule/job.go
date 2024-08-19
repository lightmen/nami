package schedule

import "github.com/lightmen/nami/pkg/cast"

type Jobber func(*Job)

type Job struct {
	Jobber     Jobber
	Key        string
	ResultChan chan *Result
	Meta       any //定义传给Job的元数据
}

// 返回处理结果
type Result struct {
	Rsp any   //存储任务处理后的回包
	Err error //存储任务处理后的错误情况
}

func NewJob(key string, jobber Jobber, meta any) *Job {
	job := &Job{
		Key:        key,
		ResultChan: make(chan *Result, 1),
		Meta:       meta,
		Jobber:     jobber,
	}

	return job
}

func (j *Job) String() string {
	if j == nil {
		return ""
	}

	return j.Key + "|" + cast.ToJson(j.Meta)
}

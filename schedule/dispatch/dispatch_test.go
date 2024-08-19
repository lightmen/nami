package dispatch

import (
	"context"
	"testing"

	"github.com/lightmen/nami/pkg/cast"
	"github.com/lightmen/nami/schedule"
	"github.com/stretchr/testify/assert"
)

func TestDispatcher(t *testing.T) {
	ctx := context.Background()

	dispatch := New(ctx, WithMaxWorker(5))
	count := 0
	jobSize := 10
	jobs := make([]*schedule.Job, jobSize)
	for i := 0; i < jobSize; i++ {
		rsp := i
		jobs[i] = &schedule.Job{
			Key: "100000" + cast.ToString(i),
			Meta: map[string]any{
				"cmd": 1001,
				"key": 2001,
			},
			Jobber: func(job *schedule.Job) {
				count++
				job.ResultChan <- &schedule.Result{
					Rsp: rsp,
					Err: nil,
				}
			},
			ResultChan: make(chan *schedule.Result, 1),
		}

		dispatch.Schedule(jobs[i])
	}

	for i := 0; i < jobSize; i++ {
		result := <-jobs[i].ResultChan
		assert.Equal(t, i, result.Rsp,
			"the rsp not correctly")
	}

	assert.Equal(t, jobSize, count,
		"the count not correctly")
}

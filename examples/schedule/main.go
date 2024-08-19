package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"time"

	"github.com/lightmen/nami/schedule"
	"github.com/lightmen/nami/schedule/dqueue"
)

func main() {
	ctx := context.Background()
	q := dqueue.New(ctx, dqueue.WithMaxWorker(5000))
	//q := dispatch.New(ctx, dispatch.WithMaxWorker(5000))

	num := 10000
	for i := 0; i < num; i++ {
		go client(q, i)
	}

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}

func client(q schedule.Scheduler, k int) {
	var seq int
	key := strconv.Itoa(k)
	ch := make(chan int, 20)

	go client_recv(key, ch)

	time.Sleep(time.Millisecond * 5)

	t := time.Now()

	for {
		rd := rand.Intn(5) + 1
		time.Sleep(time.Duration(rd) * time.Second)

		seq++
		q.Schedule(&schedule.Job{
			Key: key,
			Jobber: func(j *schedule.Job) {
				ch <- j.Meta.(int)
			},
			Meta: seq,
		})
		if time.Now().Sub(t) > time.Minute*5 {
			time.Sleep(time.Second * 30)
			fmt.Printf("key quit:%s\n", key)
			close(ch)
			return
		}
	}
}

func client_recv(key string, ch chan int) {
	defer fmt.Printf("recv quit:%s\n", key)
	var i int
	for seq := range ch {
		i++
		if seq != i {
			s := fmt.Sprintf("seq is wrong, key:%s, i:%d, seq:%d\n", key, i, seq)
			panic(s)
		}
	}
}

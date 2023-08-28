package main

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func main() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}

	// 默认本地，也可不设置
	timeOpt := cron.WithLocation(loc)

	cronJob := cron.New(timeOpt)

	task := &Task{
		Name: "Foo Job",
	}

	_, err = cronJob.AddJob("1 * * * * *", task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Sprintln("start job")

	cronJob.Start()
	// cronJob.AddFunc()

	// time.Sleep(10 * time.Second)

	// cronJob.Stop()

	// fmt.Println("job stoped")
	// time.Sleep(4 * time.Second)
	// fmt.Println("finished")

	select {}

}

/*
type Job interface {
	Run()
}
*/

type Task struct {
	Name string
}

func (t *Task) Run() {
	log.Println("Run ", t.Name, " task.")
}

package benchmark

import (
	"sync"
	"time"
)

type taskFunc = func(int) error

// BenchmarkCount 基准测试统计
type BenchmarkCount struct {
	Begin      time.Time
	End        time.Time
	RoundCount map[int]*RoundCount
}

// RoundCount 单轮统计
type RoundCount struct {
	Begin      time.Time
	End        time.Time
	TaskCounts map[int]TaskCount
}

// TaskCount 单次统计
type TaskCount struct {
	Begin  time.Time
	End    time.Time
	Status bool
}

// 轮任务锁
var roundLock = sync.Mutex{}

// 单任务锁
var taskLock = sync.Mutex{}

// RunBenchmark 一次基准测试
func RunBenchmark(tps int, rounds int, interval time.Duration, task taskFunc) (count BenchmarkCount) {
	benchmarkBegin := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(rounds)
	roundCount := map[int]*RoundCount{}
	for r := 0; r < rounds; r++ {
		go runRound(r, roundCount, tps, wg, task)
		time.Sleep(interval * time.Millisecond)
	}
	wg.Wait()
	benchmarkEnd := time.Now()
	return BenchmarkCount{benchmarkBegin, benchmarkEnd, roundCount}
}

// runRound 一轮并发
func runRound(index int, countMap map[int]*RoundCount, tps int, wg *sync.WaitGroup, task taskFunc) {
	roundWG := &sync.WaitGroup{}
	roundWG.Add(tps)
	taskCount := map[int]TaskCount{}
	roundBegin := time.Now()
	for t := 0; t < tps; t++ {
		go runTask(t, taskCount, roundWG, task)
	}
	roundWG.Wait()
	roundEnd := time.Now()
	roundLock.Lock()
	countMap[index] = &RoundCount{roundBegin, roundEnd, taskCount}
	roundLock.Unlock()
	wg.Done()
}

// runTask 单个任务
func runTask(index int, countMap map[int]TaskCount, wg *sync.WaitGroup, task taskFunc) {
	taskBegin := time.Now()
	err := task(index)
	taskEnd := time.Now()
	isSuccess := true
	if err != nil {
		isSuccess = false
	}
	taskLock.Lock()
	countMap[index] = TaskCount{taskBegin, taskEnd, isSuccess}
	taskLock.Unlock()
	wg.Done()
}

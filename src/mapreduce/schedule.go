package mapreduce

import (
	"fmt"
	"sync"
)

//
// schedule() starts and waits for all tasks in the given phase (mapPhase
// or reducePhase). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}
	//fmt.Println("Number of tasks ===",ntasks)
	var wg sync.WaitGroup
	wg.Add(ntasks)
	for i:=0;i < ntasks; i++ {
		//fmt.Println("Processing ",i)
		///ad one to wait group.
		//wg.Add(1)

		//fmt.Println("curWorker-->",curWorker);
		///Prepare the work
		var  doTaskArgs DoTaskArgs
		if phase == mapPhase {
			doTaskArgs.File = mapFiles[i]
		}
		doTaskArgs.NumOtherPhase = n_other
		doTaskArgs.Phase = phase
		doTaskArgs.TaskNumber = i
		doTaskArgs.JobName =jobName
		go func() {
			for {
				curWorker := <-registerChan
				isSuccess := call(curWorker, "Worker.DoTask", doTaskArgs, nil)
				if isSuccess {
					/// Decrement the counter when the goroutine completes.
					wg.Done()
					///Return the worker back to channel.
					registerChan <- curWorker
					break
				}else{
					fmt.Println("FAILED -------------->",curWorker)
					//If I try to add worker back to the channel it fails. I am guessing that if a worker fails for some reason
					//we should not add it back to the channel.
				}
			}

		}()
	}
	//fmt.Println("============== OUT  ============")
	///Wait till every go routine completes.
	wg.Wait()
	//close(registerChan)

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers. Once all tasks
	// have completed successfully, schedule() should return.
	//
	// Your code here (Part III, Part IV).
	//
	fmt.Printf("Schedule: %v done\n", phase)
}

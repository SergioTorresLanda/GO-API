package main 
import (
	"fmt"
	//"math/rand"
	"time"
	"sync"
)
//GOROUTINES
//CONCURRENCY : MULTIPLE TAKS RUNNING AT THE SAME TIME (1 CPU core) 
//Parallel execution (2 CPU cores) one for each task.
var wg = sync.WaitGroup{} //atomic counter
//var m = sync.Mutex{} //mutual exclusion (m.lock , m.unlock)
var mPro = sync.RWMutex{} //mutual exclusion (m.lock , m.unlock, & m.Rlock , m.Runlock)

var dbData = []string{"id1", "id2", "id3", "id4", "id5"}
var results = []string{}

func mainCoroutines(){

	t0:=time.Now()
	for i:=0; i<len(dbData); i++{
		//dbCall(i) //waits until completion (2secs) for each iteration (task) to trigger queue!
		wg.Add(1)
		go dbCall(i) //triggers every iteration at the same time ! (concurrently) no queue!
	}

	wg.Wait()
	fmt.Printf("\nTotal execution time: %v", time.Since(t0)) //6 seconds with no concurrency,
	fmt.Printf("\nResults: %v", results) //6 seconds with no concurrency,
}

func dbCall(i int){
	//Simulate DB call delay
	var delay float32 = 2000 //rand.Float32()*
	time.Sleep(time.Duration(delay)*time.Millisecond)
	fmt.Printf("\nThe result from the database is: %v", dbData[i])
	//m.Lock() // to ensure no same data allocation is accesed twice or at the same time. //SIMPLE MUTEX
	//results = append(results, dbData[i])
	//m.Unlock()
	save(dbData[i])
	log()
	wg.Done()
}

func save(result string){
	mPro.Lock() //all locks must be cleared to proceed !! (Full & Read)
	results = append(results, result)
	mPro.Unlock()
}

func log(){
	mPro.RLock() ////only full locks must (writings) be cleared to proceed to read !! (Full)
	fmt.Printf("\nCurr Results are: %v", results)
	mPro.RUnlock()
}

//PERFORMANCE WILL BE LIMITED FOR THE AMOUNT OF CORES OF THE CPU !!
//8 CORES MEANS WE ARE ABLE TU RUN 8 GOROUTINES (ITERATIONS) AT THE SAME TIME . OTHERS WILL QUEUE

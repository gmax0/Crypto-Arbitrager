package main

import (
    "fmt"
    "time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
    //Block until jobs are available
    for j:= range jobs {
        fmt.Println("worker", id, "started job", j)
        time.Sleep(time.Second)
        fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}

func main() {
    const numberOfJobs = 6

    jobs := make(chan int, numberOfJobs)
    results := make(chan int, numberOfJobs)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numberOfJobs; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= numberOfJobs; a++ {
        <- results
    }
    close(results)
}
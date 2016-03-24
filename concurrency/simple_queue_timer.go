package main

import "fmt"

/*
#include <sys/time.h>
int GetTickCount()
{
        struct timeval tv;
        if(gettimeofday(&tv, 0) != 0)
                return 0;

        return (tv.tv_sec * 1000) + (tv.tv_usec / 1000);
}
*/
import "C"

type timer struct {
    id_ int
    start_tick_ int
    ticks_ int
    alive_ bool
}

const NUM_WORKERS = 8

//timer queue
var ch chan *timer = make (chan *timer, 1024)

func main() {
    fmt.Println("start")
    //start worker threads
    for i := 0; i < NUM_WORKERS; i++ {
        go timerWorker()
    }
    
    //start test
    benchmark()
    
    var input string
    fmt.Scanln(&input)
}

func benchmark() {
    for i := 0; i < 1024; i++ {
        t := timer{ id_:i }
        t.ticks_ = 100
        t.start()
    }
}

//starts the timer
func (this *timer) start() {
    this.start_tick_ = int(C.GetTickCount())
    this.alive_ = true
    this.enqueue()
}

//enqueues the timer to timer queue
func (this *timer) enqueue() {
    ch <- this
}

//stops the timer
func (this *timer) stop() {
    this.alive_ = false
}

//the timer's callback
func (this *timer) callback(elapsed int) {
    fmt.Println("timer", this.id_, "ticks", this.ticks_, "elapsed", elapsed)
}

func timerWorker() {
    for {
        var it_ *timer
        select {
            case it_ = <- ch:
                tick := int(C.GetTickCount())
                if tick - it_.start_tick_ >= it_.ticks_ && it_.alive_ {
                    it_.callback(tick - it_.start_tick_)
                    it_.start_tick_ = tick
                }
                ch <- it_
            default:
                it_ = nil
        }
    }
}
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
    id int
    startTick int
    ticks int
    alive bool
}

var ch chan *timer = make (chan *timer, 1024)

func main() {
    fmt.Println("wat")
    for i := 0; i < 8; i++ {
        go timerWorker()
    }
    
    for i := 0; i < 1024; i++ {
        t := timer{ id:i }
        t.ticks = 100
        t.start()
    }
    fmt.Println("allocated")
    var input string
    fmt.Scanln(&input)
}

func (this *timer) start() {
    this.startTick = int(C.GetTickCount())
    this.alive = true
    this.enqueue()
}

func (this *timer) enqueue() {
    ch <- this
}

func (this *timer) callback(elapsed int) {
    fmt.Println("timer", this.id, "ticks", this.ticks, "elapsed", elapsed)
}

func timerWorker() {
    for {
        var it *timer
        select {
            case it = <- ch:
                tick := int(C.GetTickCount())
                if tick - it.startTick >= it.ticks && it.alive {
                    it.callback(tick - it.startTick)
                    it.startTick = tick
                }
                ch <- it
            default:
                it = nil
        }
    }
}
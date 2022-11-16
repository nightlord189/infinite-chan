package infinite_test

import (
	"fmt"
	"github.com/nightlord189/infinite-chan"
	"runtime"
	"testing"
	"time"
)

const inputCount = 10000

func BenchmarkInfiniteChannel(b *testing.B) {
	ch := infinite.NewChan()
	allcount := inputCount
	go func() {
		for i := 0; i < allcount; i++ {
			ch.In() <- i
		}
		fmt.Println("infinite goroutines_count", runtime.NumGoroutine())
	}()
	count := 0
	resultMap := make(map[int]bool)
	for i := range ch.Out() {
		//fmt.Println("read", iVal)
		count++
		time.Sleep(1 * time.Millisecond)

		resultMap[i.(int)] = true

		if count == allcount {
			close(ch.In())
		}
	}
	if count != allcount {
		b.Fatalf("count is %d, should be %d\n", count, allcount)
	}
	for i := 0; i < allcount; i++ {
		if _, ok := resultMap[i]; !ok {
			b.Fatalf("%d doesn't exists in result\n", i)
		}
	}
}

package infinite_test

import (
	"github.com/nightlord189/infinite-chan"
	"runtime"
	"testing"
	"time"
)

func TestInfiniteChannel(t *testing.T) {
	ch := infinite.NewChan[int]()
	allCount := 10000
	go func() {
		for i := 0; i < allCount; i++ {
			ch.In() <- i
		}
	}()

	count := 0
	resultMap := make(map[int]bool)
	for i := range ch.Out() {
		//fmt.Println("read", iVal)
		count++
		time.Sleep(1 * time.Millisecond)

		resultMap[i] = true

		if count == allCount {
			close(ch.In())
		}
	}
	if count != allCount {
		t.Fatalf("count is %d, should be %d\n", count, allCount)
	}
	for i := 0; i < allCount; i++ {
		if _, ok := resultMap[i]; !ok {
			t.Fatalf("%d doesn't exists in result\n", i)
		}
	}
}

func TestMemory(t *testing.T) {
	printAlloc(t)

	ch := infinite.NewChan[[128]byte]()
	allCount := 1000000

	signalCh := make(chan bool)

	go func() {
		for i := 0; i < allCount; i++ {
			ch.In() <- [128]byte{}
		}
		t.Logf("ch size: %d\n", ch.Len())
		ch.Close()
		signalCh <- true
	}()

	<-signalCh

	count := 0
	t.Logf("starting reading\n")
	for range ch.Out() {
		//fmt.Println("read", iVal)
		count++
	}

	t.Logf("reading finished")

	t.Logf("ch size: %d\n", ch.Len())
	printAlloc(t)
	if count != allCount {
		t.Fatalf("count is %d, should be %d\n", count, allCount)
	}

	t.Logf("resize\n")

	ch.Resize()
	runtime.GC()
	t.Logf("ch size: %d\n", ch.Len())
	printAlloc(t)
}

func printAlloc(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	t.Logf("alloc: %d KB\n", m.Alloc/1024)
}

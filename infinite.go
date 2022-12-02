package infinite

import (
	"math"
	"sync"
)

const maxBufferSize = math.MaxInt
const halfBufferSize = math.MaxInt / 2

type Chan[T any] struct {
	in, out chan T
	signal  chan bool
	head    int
	tail    int
	buf     map[int]T
	lock    *sync.Mutex
	closed  bool
}

func NewChan[T any]() *Chan[T] {
	result := &Chan[T]{
		in:     make(chan T),
		out:    make(chan T),
		signal: make(chan bool, 1),
		head:   -1,
		tail:   -1,
		buf:    make(map[int]T),
		lock:   &sync.Mutex{},
		closed: false,
	}
	go result.processInput()
	go result.processOutput()
	return result
}

func (c *Chan[T]) In() chan T {
	return c.in
}

func (c *Chan[T]) Out() chan T {
	return c.out
}

func (c *Chan[T]) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.buf)
}

func (c *Chan[T]) Close() {
	close(c.in)
}

// Head - returns current head index, for example, to analyze do you need resize chan now
func (c *Chan[T]) Head() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.head
}

// Resize - resizes channel index, copies all values to new map to free memory, prevent memory leaks and reduce indexes
func (c *Chan[T]) Resize() {
	// someday tail index can be closer to maxInt
	// in that case we will create new buffer map, reduce head and tail both and move all elements to new buffer
	//fmt.Printf("resize started: head %d, tail %d, buf %d\n", c.head, c.end, len(c.buf))
	offset := c.head
	newBuf := make(map[int]T, len(c.buf))
	for i := 0; i <= c.tail-offset; i++ {
		newBuf[i] = c.buf[i+offset]
	}
	c.head = 0
	c.tail -= offset
	c.buf = newBuf
	//fmt.Printf("resize finished: head %d, tail %d, buf %d\n", c.head, c.end, len(c.buf))
}

func (c *Chan[T]) processInput() {
	for inVal := range c.in {
		c.lock.Lock()
		// increase index of last element
		c.tail++
		if c.head < 0 {
			c.head++
		}
		c.addToBuffer(inVal)
		c.lock.Unlock()

		// if buffer is empty - maybe processOutput is waiting new element
		if len(c.signal) == 0 {
			c.signal <- true
		}
	}
	c.lock.Lock()
	c.closed = true
	c.lock.Unlock()

	c.signal <- true
}

func (c *Chan[T]) addToBuffer(val T) {
	c.buf[c.tail] = val
	if c.tail == maxBufferSize && c.head > halfBufferSize {
		c.Resize()
	}
}

func (c *Chan[T]) processOutput() {
	var outVal T
	for {
		c.lock.Lock()
		if c.closed && c.head > c.tail {
			c.lock.Unlock()
			break
		}
		if c.head <= c.tail && c.head >= 0 {
			outVal = c.buf[c.head]
			delete(c.buf, c.head)
			c.head++

			c.lock.Unlock()

			c.out <- outVal
		} else {
			c.lock.Unlock()
			<-c.signal //wait until we go new value on input channel
		}
	}
	//fmt.Println("buffer_size", len(c.buf))
	close(c.out)
}

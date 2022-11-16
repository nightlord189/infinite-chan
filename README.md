# Infinite channel

Go struct for implement Go channel concept with infinite buffer.

## Usecase
+ If you don't know what size of buffer can require your channel in the future, but also you need non-locking writing and locking reading from channel.

## How it works
+ Two channels for input and output
+ buffer as map
+ two indexes - head and tail 
+ when producer writes to channel, tail is incrementing
+ when consumer reads from channel, head is incrementing
+ every element in buffer is stored as value in map, where keys are consecutive integer values (including head and tail)
+ if head and tail values are too large and close to max integer value, buffer will be resized (head and tail will be moved closer to zero with offset)

## Tips
+ to close all channel, close .In() channel

## Getting started

```
go get github.com/nightlord189/infinite-chan
```

```go
package main

import (
	"fmt"
	"github.com/nightlord189/infinite-chan"
)

func main () {
	val := 1
	ch := infinite.NewChan()
	ch.In() <- val
	value := <- ch.Out()
	fmt.Println(value.(int))

	close (ch.In())
}
```
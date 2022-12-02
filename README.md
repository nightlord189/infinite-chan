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
+ because map in Go are never shrinking even if all keys are deleted (see [article](https://teivah.medium.com/maps-and-memory-leaks-in-go-a85ebe6e7e69)), so if your channel is alive too long, you should periodically (or based on your application's logic) call Resize() to prevent memory leaks

## Getting started

```
go get github.com/nightlord189/infinite-chan@v1.1.0
```

```go
package main

import (
	"fmt"
	"github.com/nightlord189/infinite-chan"
)

func main () {
	val := 1
	ch := infinite.NewChan[int]()
	ch.In() <- val
	value := <- ch.Out()
	fmt.Println(value)

	close (ch.In())
}
```
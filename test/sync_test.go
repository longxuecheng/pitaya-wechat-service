package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	w := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func() {
			Log(os.Stdout, "path", "/search?q=flowers\n")
			w.Done()
		}()
	}
	w.Wait()
	fmt.Println("all options done!")
}

func TestSyncAddWithPool(t *testing.T) {
	w := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		w.Add(1)
		go func() {
			sum := intPool.Get().(int)
			fmt.Printf("sum is %d \n", sum)
			sum = sum + 1
			intPool.Put(sum)
			w.Done()
		}()
	}
	w.Wait()
	fmt.Printf("final sum is %d ", intPool.Get().(int))
}

var intPool = sync.Pool{
	New: func() interface{} {
		return 0
	},
}

var bufPool = sync.Pool{
	New: func() interface{} {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return new(bytes.Buffer)
	},
}

// timeNow is a fake version of time.Now for tests.
func timeNow() time.Time {
	return time.Unix(1136214245, 0)
}

func Log(w io.Writer, key, val string) {
	b := bufPool.Get().(*bytes.Buffer)
	fmt.Printf("point position is %p \n", b)
	b.Reset()
	// Replace this with time.Now() in a real logger.
	b.WriteString(timeNow().UTC().Format(time.RFC3339))
	b.WriteByte(' ')
	b.WriteString(key)
	b.WriteByte('=')
	b.WriteString(val)
	w.Write(b.Bytes())
	bufPool.Put(b)
}

package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/muesli/cache2go"
)

// Keys & values in cache2go can be of arbitrary types, e.g. a struct.
type myStruct struct {
	text     string
	moreData []byte
}

func TestCache(t *testing.T) {
	// Accessing a new cache table for the first time will create it.
	cache := cache2go.Cache("myCache")

	// We will put a new item in the cache. It will expire after
	// not being accessed via Value(key) for more than 5 seconds.
	val := myStruct{"This is a test!", []byte{}}
	cache.Add("someKey", 5*time.Second, &val)

	// Let's retrieve the item from the cache.
	res, err := cache.Value("someKey")
	if err == nil {
		t.Log("Found value in cache:", res.Data().(*myStruct).text)
	} else {
		t.Log("Error retrieving value from cache:", err)
	}

	// Wait for the item to expire in cache.
	time.Sleep(6 * time.Second)
	res, err = cache.Value("someKey")
	if err != nil {
		t.Log("Item is not cached (anymore).")
	}

	// Add another item that never expires.
	cache.Add("someKey", 0, &val)

	// cache2go supports a few handy callbacks and loading mechanisms.
	cache.SetAboutToDeleteItemCallback(func(e *cache2go.CacheItem) {
		t.Log("Deleting:", e.Key(), e.Data().(*myStruct).text, e.CreatedOn())
	})

	// Remove the item from the cache.
	cache.Delete("someKey")

	// And wipe the entire cache table.
	cache.Flush()
}

func TestMultiSameKeyCache(t *testing.T) {
	cache := cache2go.Cache("myCache")
	item := cache.Add("k1", 5*time.Second, "v1")
	fmt.Println(item.CreatedOn())
	time.Sleep(4 * time.Second)
	cache.Add("k1", 5*time.Second, "v2")
	item, err := cache.Value("k1")
	if err != nil {
		t.Error("item of k1 is not exists!")
	}
	fmt.Println(item.CreatedOn())
	fmt.Println("total item number is ", cache.Count())
	fmt.Println(item.Data())
}

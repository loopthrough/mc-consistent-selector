package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/loopthrough/mc-consistent-selector/selector"
	"log"
	"math/rand"
	"strconv"
)

func main() {
	// Assumes you have started these instances locally
	sl := selector.NewServerList()
	sl.Add("0.0.0.0:11211")
	sl.Add("0.0.0.0:11212")
	sl.Add("0.0.0.0:11213")

	cache := memcache.NewFromSelector(sl)

	err := cache.Set(&memcache.Item{
		Key:   "test",
		Value: []byte("test"),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = cache.Delete("test")
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 1000; i++ {
		key := rand.Int()
		cache.Set(&memcache.Item{
			Key:   strconv.Itoa(key),
			Value: []byte(strconv.Itoa(key)),
		})
	}
}

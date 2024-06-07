package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

// in etcd, the key space is organized like a filesystem hierarchy
// the forward slash (/) is used as a delimiter to create a hierarchical structure for the keys
// this allows for better organization and grouping of related keys

var (
	etcdEndpoints = flag.String("endpoints", "etcd:2379", "comma-separated list of etcd endpoints")
	lockName      = flag.String("lock-name", "/my-lock", "name of the lock")
	id            = flag.String("id", "", "unique ID of the locker")
)

func main() {
	flag.Parse()

	// create an etcd client
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{*etcdEndpoints},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	log.Printf("%s successfully connected to etcd", *id)

	// create a session to acquire locks
	session, err := concurrency.NewSession(client)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// create a mutex for the lock
	mutex := concurrency.NewMutex(session, *lockName)

	key := "/my-key"

	var loopCount int

	for {
		// wait for a random duration between 1 and 5 seconds before trying to acquire the lock
		waitDuration := time.Duration(1+rand.Intn(4)) * time.Second
		log.Printf("%s waiting for %s before attempting to acquire the lock...", *id, waitDuration)
		time.Sleep(waitDuration)

		log.Printf("%s ⏳ attempting to acquire the lock...", *id)

		// acquire the lock
		err := mutex.Lock(context.TODO())
		if err != nil {
			log.Fatalf("error: failed to acquire the lock: %v", err)
		}
		log.Printf("%s 🔒 ACQUIRED the lock", *id)

		log.Printf("%s ⏳ attempting to update the value...", *id)

		newValue := fmt.Sprintf("%s-value-%d", *id, loopCount)

		// put the key-value pair
		putResponse, err := client.Put(context.TODO(), key, newValue)
		if err != nil {
			log.Fatalf("error: failed to update value: %v", err)
		}

		log.Printf("%s put response: %v", *id, putResponse)

		log.Printf("%s ✅ updated the value to: %s", *id, newValue)

		// hold the lock for a random duration between 1 and 5 seconds
		waitDuration = time.Duration(1+rand.Intn(4)) * time.Second
		log.Printf("%s HOLDING the lock for %s seconds...", *id, waitDuration)
		time.Sleep(waitDuration)

		// release the lock
		err = mutex.Unlock(context.TODO())
		if err != nil {
			log.Fatalf("error: failed to release the lock: %v", err)
		}
		log.Printf("%s 🔓 RELEASED the lock", *id)

		loopCount++
	}
}

package main

import (
	"sync"
)

// Bucket is a single spot in the table
type Bucket struct {
	mutex *sync.RWMutex
	entry *Entry
}

// Entry is a single node in the bucket
type Entry struct {
	key   [2]float32
	value [2]float32
	next  *Entry
	mutex *sync.Mutex
}

// SafeMap is a custom concurrent map
type SafeMap struct {
	con              []Bucket
	size             int
	insertions       int // Debug
	wastedInsertions int // Debug -- waited for writer lock but someone else wrote key
}

// Initialize map to size - won't ever resize
func (sm *SafeMap) init(size int) {
	sm.con = make([]Bucket, size)
	for i := 0; i < size; i++ {
		sm.con[i].mutex = &sync.RWMutex{}
	}
	sm.size = size
}

func hash(key [2]float32) int {
	return int(key[0])
}

// Insert a key into the map
func (sm *SafeMap) put(key [2]float32, value [2]float32) {
	// Find bucket, acquire read lock
	// See if my key is there
	// If my key is there
	// Acquire mutex for key, add my key, release mutex for key
	// Else
	// Release read lock, acquire write lock for bucket
	// Acquire write lock for bucket
	// Go to end of linked list and add my key
	bucketNum := hash(key) % sm.size
	rwLock := sm.con[bucketNum].mutex
	rwLock.RLock()
	entry := findEntry(&sm.con[bucketNum], key)
	if entry != nil {
		// My key is here, acquire the mutex and add
		entry.mutex.Lock()
		entry.value = value
		entry.mutex.Unlock()
		rwLock.RUnlock()
	} else {
		// Key hasn't been inserted before
		rwLock.RUnlock()

		// Acquire write lock so we can add a new linked list node
		rwLock.Lock()

		// Make sure someone didn't add it while we were waiting for lock
		entry = findEntry(&sm.con[bucketNum], key)
		if entry != nil {
			// Someone did it for us, we are the only one in the whole bucket so just append
			entry.value = value
			sm.wastedInsertions++
		} else {
			addEntry(&sm.con[bucketNum], key, value)
			sm.insertions++
		}
		rwLock.Unlock()
	}

}

// Find a entry in a bucket - return nil if does not exist
func findEntry(bucket *Bucket, key [2]float32) *Entry {
	entry := bucket.entry
	for entry != nil && entry.key != key {
		entry = entry.next
	}
	return entry
}

// Add an entry to the bucket - assumes the lock for the bucket has been acquired
func addEntry(bucket *Bucket, key [2]float32, value [2]float32) {
	// See if the bucket has anything
	if bucket.entry == nil {
		bucket.entry = &Entry{key, value, nil, &sync.Mutex{}}
	} else {
		// Add to end of linked list!
		trail := bucket.entry
		cur := bucket.entry.next
		for cur != nil {
			trail = cur
			cur = cur.next
		}
		trail.next = &Entry{key, value, nil, &sync.Mutex{}}
	}
}

// NOT concurrent
func (sm *SafeMap) get(key [2]float32) [2]float32 {
	bucketNum := hash(key) % sm.size
	rwLock := sm.con[bucketNum].mutex
	rwLock.RLock()
	entry := findEntry(&sm.con[bucketNum], key)
	var res [2]float32
	entry.mutex.Lock()
	res = entry.value
	entry.mutex.Unlock()
	rwLock.RUnlock()
	return res
}

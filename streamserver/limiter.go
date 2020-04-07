package main

import "log"

type ConnLimiter struct {
	bucket        chan int
	concurrentCap int
}

func NewLimiter(concurrentCap int) *ConnLimiter {
	bucket := make(chan int, concurrentCap)
	return &ConnLimiter{bucket: bucket, concurrentCap: concurrentCap}
}

func (l ConnLimiter) openConn() bool {
	if len(l.bucket) >= l.concurrentCap {
		log.Printf("Reach the rate limitation")
		return false
	}

	l.bucket <- 1
	log.Printf("Open new connection")
	return true
}

func (l ConnLimiter) releaseConn() {
	<-l.bucket
	log.Printf("Release a connection")
}

package main

import (
	"os"
	"strconv"
	"time"
)

var (
	publicMemberID string
	publicPassHash string
	rootUrl        string
	connTimeout    time.Duration
	connKeepAlive  time.Duration
	connMaxIdle    int
)

func init() {
	publicMemberID = os.Getenv("PANDA_PUBLIC_MEMBER_ID")
	publicPassHash = os.Getenv("PANDA_PUBLIC_PASS_HASH")
	rootUrl = os.Getenv("PANDA_ROOT_URL")
	timeoutEnv := os.Getenv("PANDA_CONN_TIMEOUT")
	if timeoutEnv == "" {
		connTimeout = time.Second * 30
	} else {
		timeout, err := strconv.Atoi(timeoutEnv)
		if err != nil {
			panic(err)
		}
		connTimeout = time.Second * time.Duration(timeout)
	}
	keepAliveEnv := os.Getenv("PANDA_CONN_KEEP_ALIVE")
	if keepAliveEnv == "" {
		connKeepAlive = time.Second * 30
	} else {
		keepAlive, err := strconv.Atoi(keepAliveEnv)
		if err != nil {
			panic(err)
		}
		connKeepAlive = time.Second * time.Duration(keepAlive)
	}
	maxIdleEnv := os.Getenv("PANDA_CONN_MAX_IDLE")
	if maxIdleEnv == "" {
		connMaxIdle = 16
	} else {
		maxIdle, err := strconv.Atoi(maxIdleEnv)
		if err != nil {
			panic(err)
		}
		connMaxIdle = maxIdle
	}
}

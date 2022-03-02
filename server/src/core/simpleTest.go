package main

import (
	"encoding/json"
	"fmt"
)

type SyncPID struct {
	PID     int `json:"PID"`
	TestVal int
}

func main() {
	Map := make(map[int]interface{})
	s := &SyncPID{PID: 99999, TestVal: 1111}
	Map[1] = s

	newS := Map[1]

	tests := SyncPID{PID: 88888}
	msg1, _ := json.Marshal(s)
	msg2, _ := json.Marshal(tests)
	//json.Unmarshal(msg, newS)

	fmt.Println(string(msg1))
	fmt.Println(string(msg2))
	json.Unmarshal(msg2, newS)
	fmt.Println(newS)
}

package queue

import (
	"testing"
)

func TestPush(t *testing.T) {
	// Cleanup at first
	allQueues = make(map[string][]interface{})

	Push("test", "value")
	Push("test", "value2")
	Push("test2", "value")

	// Test for key test
	res := allQueues["test"]
	expect := []string{"value", "value2"}
	if len(res) != len(expect) {
		t.Errorf("expect %v for test key, but got %v", expect, res)
	}
	for i, v := range res {
		if v != expect[i] {
			t.Errorf("expect %v for test key, but got %v", expect, res)
		}
	}

	// Test for invalid key
	res = allQueues["invalid_key"]
	if len(res) != 0 {
		t.Errorf("expect empty for invalid key, but got %v", res)
	}
}

func TestPop(t *testing.T) {
	// Set data
	allQueues = make(map[string][]interface{})
	allQueues["test"] = []interface{}{"value", "value2"}

	res := Pop("test")
	if res != "value" {
		t.Errorf("expect value, but got %v", res)
	}
	remaining := allQueues["test"]
	if len(remaining) != 1 || remaining[0] != "value2" {
		t.Errorf("expect remaining value2 only at allQueues, but got %v", remaining)
	}
}

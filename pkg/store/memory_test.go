package store

import (
	"context"
	"testing"
)

func TestGet(t *testing.T) {
	dataa := make(map[string]interface{})
	datab := make(map[string]interface{})
	datac := make(map[string]interface{})
	datac["c"] = "hello"
	datab["b"] = datac
	dataa["a"] = datab
	memStore := NewMemoryStoreWithData(dataa)
	val := memStore.GetString(context.Background(), "a.b.c")
	if val != "hello" {
		t.Fail()
	}
	// println("hello")
	// out, _ := yaml.Marshal(dataa)
	// println(string(out))
}

func TestSetGet(t *testing.T) {
	memStore := NewMemoryStore()
	memStore.Set(context.Background(), "a.b.c", "hello")
	val := memStore.GetString(context.Background(), "a.b.c")
	if val != "hello" {
		t.Fail()
	}
	// println("hello")
	// out, _ := yaml.Marshal(memStore.DataMap(context.Background()))
	// println(string(out))
}

package storage

import (
	"testing"
)

func TestAdd(t *testing.T) {
	ss := NewSessionStorage()
	key := "testKey"
	value := "testValue"

	ss.Add(key, value)

	if _, ok := ss.Values[key]; !ok {
		t.Errorf("Expected key %s to be in Values map", key)
	}
	if ss.Values[key] != value {
		t.Errorf("Expected value %s for key %s, but got %s", value, key, ss.Values[key])
	}
}

func TestGet(t *testing.T) {
	ss := NewSessionStorage()
	key := "testKey"
	value := "testValue"

	ss.Add(key, value)

	retrievedValue, ok := ss.Get(key)
	if !ok {
		t.Errorf("Expected key %s to be in Values map", key)
	}
	if retrievedValue != value {
		t.Errorf("Expected value %s for key %s, but got %s", value, key, retrievedValue)
	}
}

func TestGetAll(t *testing.T) {
	ss := NewSessionStorage()
	key1 := "testKey1"
	value1 := "testValue1"
	key2 := "testKey2"
	value2 := "testValue2"

	ss.Add(key1, value1)
	ss.Add(key2, value2)

	allValues := ss.GetAll()
	if len(allValues) != 2 {
		t.Errorf("Expected 2 items in Values map, but got %d", len(allValues))
	}
	if allValues[key1] != value1 {
		t.Errorf("Expected value %s for key %s, but got %s", value1, key1, allValues[key1])
	}
	if allValues[key2] != value2 {
		t.Errorf("Expected value %s for key %s, but got %s", value2, key2, allValues[key2])
	}
}

func BenchmarkAdd(b *testing.B) {
	ss := NewSessionStorage()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			ss.Add(string(rune(j+200)), "test_value")
		}
	}
}

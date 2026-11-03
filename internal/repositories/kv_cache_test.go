package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleKVCache(t *testing.T) {
	t.Parallel()

	type operation struct {
		name      string
		action    string // "put", "get", "delete", or "clear"
		key       string
		value     any
		wantValue any
		wantFound bool
	}

	tests := []struct {
		name string
		ops  []operation
	}{
		{
			name: "returns cache miss for unknown key",
			ops: []operation{
				{name: "get missing key", action: "get", key: "missing", wantFound: false},
			},
		},
		{
			name: "returns stored values",
			ops: []operation{
				{name: "put integer", action: "put", key: "count", value: 1},
				{name: "put string", action: "put", key: "name", value: "Acme"},
				{name: "get integer", action: "get", key: "count", wantValue: 1, wantFound: true},
				{name: "get string", action: "get", key: "name", wantValue: "Acme", wantFound: true},
			},
		},
		{
			name: "updates an existing key",
			ops: []operation{
				{name: "put initial value", action: "put", key: "count", value: 1},
				{name: "replace value", action: "put", key: "count", value: 2},
				{name: "get replacement", action: "get", key: "count", wantValue: 2, wantFound: true},
			},
		},
		{
			name: "delete removes only the requested key",
			ops: []operation{
				{name: "put first value", action: "put", key: "first", value: 1},
				{name: "put second value", action: "put", key: "second", value: 2},
				{name: "delete first value", action: "delete", key: "first"},
				{name: "first key is missing", action: "get", key: "first", wantFound: false},
				{name: "second key remains", action: "get", key: "second", wantValue: 2, wantFound: true},
			},
		},
		{
			name: "clear removes all entries",
			ops: []operation{
				{name: "put first value", action: "put", key: "first", value: 1},
				{name: "put second value", action: "put", key: "second", value: 2},
				{name: "clear cache", action: "clear"},
				{name: "first key is missing", action: "get", key: "first", wantFound: false},
				{name: "second key is missing", action: "get", key: "second", wantFound: false},
			},
		},
		{
			name: "distinguishes a stored nil value from a cache miss",
			ops: []operation{
				{name: "put nil", action: "put", key: "nil-value", value: nil},
				{name: "get stored nil", action: "get", key: "nil-value", wantValue: nil, wantFound: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewSimpleKVCache()
			for _, op := range tt.ops {
				switch op.action {
				case "put":
					cache.Put(op.key, op.value)
				case "get":
					got, found := cache.Get(op.key)
					assert.Equal(t, op.wantFound, found, op.name)
					if found {
						assert.Equal(t, op.wantValue, got, op.name)
					}
				case "delete":
					cache.Delete(op.key)
				case "clear":
					cache.Clear()
				default:
					t.Fatalf("unknown action: %q", op.action)
				}
			}
		})
	}
}

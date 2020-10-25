package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetAndExists(t *testing.T) {
	store := NewMemoryStorage()

	tests := []string{"foobar", "123123", "", "ABCD", "!§$&$%&$%/%&(", "this is an example", "hello world!"}
	for _, test := range tests {
		err := store.Set(context.Background(), test)
		assert.NoError(t, err)

		exists, err := store.Exists(context.Background(), test)
		assert.NoError(t, err)
		assert.True(t, exists)
	}
}

func Test_NotExists(t *testing.T) {
	store := NewMemoryStorage()

	err := store.Set(context.Background(), "foo bar")
	assert.NoError(t, err)

	tests := []string{"foobar", "123123", "", "ABCD", "!§$&$%&$%/%&(", "this is an example", "hello world!"}
	for _, test := range tests {
		exists, err := store.Exists(context.Background(), test)
		assert.NoError(t, err)
		t.Log(test)
		assert.False(t, exists)
	}
}

func Test_Limit(t *testing.T) {
	store := NewMemoryStorage()

	limit = 5
	deleteSize = 2

	tests := []string{"1", "2", "3", "4", "5", "6", "7"}
	for _, test := range tests {
		err := store.Set(context.Background(), test)
		assert.NoError(t, err)
		t.Log(store.list)
	}
	assert.Equal(t, limit, len(store.list))
	assert.Equal(t, []string{"3", "4", "5", "6", "7"}, store.list)
}

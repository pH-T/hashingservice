package service

import (
	"context"
	"hashservice/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashAndVerify(t *testing.T) {
	store := storage.NewMemoryStorage()
	s := NewBcryptHashService(store)

	passwords := []string{"foobar", "123123", "", "ABCD", "!§$&$%&$%/%&("}

	for _, pw := range passwords {
		hash, err := s.Hash(context.Background(), pw)
		assert.NoError(t, err)
		assert.True(t, hash != "")

		v, self, err := s.Verify(context.Background(), pw, hash)
		assert.NoError(t, err)
		assert.True(t, v)
		assert.True(t, self)
	}

}

func Test_Verify(t *testing.T) {
	store := storage.NewMemoryStorage()
	s := NewBcryptHashService(store)

	testCases := []struct {
		pw   string
		hash string
	}{
		{
			pw:   "foobar",
			hash: "",
		},
		{
			pw:   "foobar",
			hash: "123123123",
		},
		{
			pw:   "foobar",
			hash: "asdfasdfoiasjdf",
		},
		{
			pw:   "foobar",
			hash: "'*_:;§$%$&/",
		},
	}

	for _, elem := range testCases {
		v, self, err := s.Verify(context.Background(), elem.pw, elem.hash)
		assert.Error(t, err)
		assert.False(t, v)
		assert.False(t, self)
	}
}

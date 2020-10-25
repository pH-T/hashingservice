package http

import (
	"encoding/json"
	"hashservice/service"
	"hashservice/storage"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HTTPServerHealth(t *testing.T) {
	store := storage.NewMemoryStorage()
	hasher := service.NewBcryptHashService(store)

	server := httptest.NewServer(initRoutes(hasher))
	defer server.Close()

	resp, err := http.Get(server.URL + "/health")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	actual, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "up", string(actual))
}

type hashResp struct {
	Hash string `json:"hash"`
}

type verifyResp struct {
	Verified bool `json:"verified"`
	Selfmade bool `json:"selfmade"`
}

func Test_HTTPServerHashAndVerify(t *testing.T) {
	store := storage.NewMemoryStorage()
	hasher := service.NewBcryptHashService(store)

	server := httptest.NewServer(initRoutes(hasher))
	defer server.Close()

	r := strings.NewReader(`{"pw":"foobar"}`)
	resp, err := http.Post(server.URL+"/hash", "application/json", r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	b, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var hresp hashResp
	err = json.Unmarshal(b, &hresp)
	assert.NoError(t, err)

	r = strings.NewReader(`{"pw":"foobar", "hash": "` + hresp.Hash + `"}`)
	resp, err = http.Post(server.URL+"/verify", "application/json", r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	b, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()

	var vresp verifyResp
	err = json.Unmarshal(b, &vresp)
	assert.NoError(t, err)
	assert.True(t, vresp.Verified)
	assert.True(t, vresp.Selfmade)

}

func Test_HTTPServerVerify1(t *testing.T) {
	store := storage.NewMemoryStorage()
	hasher := service.NewBcryptHashService(store)

	server := httptest.NewServer(initRoutes(hasher))
	defer server.Close()

	r := strings.NewReader(`{"pw":"foobar", "hash": "1234"}`)
	resp, err := http.Post(server.URL+"/verify", "application/json", r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

func Test_HTTPServerVerify2(t *testing.T) {
	store := storage.NewMemoryStorage()
	hasher := service.NewBcryptHashService(store)

	server := httptest.NewServer(initRoutes(hasher))
	defer server.Close()

	r := strings.NewReader(`{"pw":"foo", "hash": "$2a$14$DewuCqBaOSjOVwQ3bhBsnORYdZUeXQ5i00D5b9l1NYgd1ND6zisq2"}`)
	resp, err := http.Post(server.URL+"/verify", "application/json", r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

}

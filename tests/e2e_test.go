// +build e2e

package tests

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jonfriesen/redis-proxy/api"
	"github.com/jonfriesen/redis-proxy/cache"
	"github.com/jonfriesen/redis-proxy/storage/redis"
)

func Test_SimpleGet(t *testing.T) {
	cache := cache.New(int32(2), time.Duration(60000))

	hm := map[string]string{
		"han":      "solo",
		"princess": "leia",
	}

	r, err := redis.New("redis", "6379")
	if err != nil {
		log.Fatalf("Error creating Redis connection: %+v", err)
	}
	seedDatasource(hm, r)

	dataSource := NewStorageMiddleware(r)

	handler := api.New(dataSource.asStorage(), cache)

	s := httptest.NewServer(handler)
	defer s.Close()

	res, err := http.Get(fmt.Sprintf("%s/v1/get/princess", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	expected := "leia"

	if expected != string(body) {
		t.Fatalf("Expected %s but received %s", expected, body)
	}

	if dataSource.callList[0] != "GET$princess" {
		t.Fatalf("Call did not reach data source for: %s", "GET$princess")
	}
}

func Test_TimeExpiry(t *testing.T) {
	cache := cache.New(int32(2), time.Duration(1))

	hm := map[string]string{
		"han":      "solo",
		"princess": "leia",
	}

	r, err := redis.New("redis", "6379")
	if err != nil {
		log.Fatalf("Error creating Redis connection: %+v", err)
	}
	seedDatasource(hm, r)

	dataSource := NewStorageMiddleware(r)

	handler := api.New(dataSource.asStorage(), cache)

	s := httptest.NewServer(handler)
	defer s.Close()

	_, err = http.Get(fmt.Sprintf("%s/v1/get/princess", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10)

	_, err = http.Get(fmt.Sprintf("%s/v1/get/princess", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	if dataSource.table["GET$princess"] != 2 {
		t.Fatalf("Cache used when record *should* have been expired, expected 2 calls actually: %v", dataSource.table["GET$princess"])
	}
}

func Test_LRUEviction(t *testing.T) {
	cache := cache.New(int32(2), time.Duration(60000))

	hm := map[string]string{
		"han":      "solo",
		"princess": "leia",
		"malcolm":  "reynolds",
	}

	r, err := redis.New("redis", "6379")
	if err != nil {
		log.Fatalf("Error creating Redis connection: %+v", err)
	}
	seedDatasource(hm, r)

	dataSource := NewStorageMiddleware(r)

	handler := api.New(dataSource.asStorage(), cache)

	s := httptest.NewServer(handler)
	defer s.Close()

	// first call populates our cache (size of 2)
	_, err = http.Get(fmt.Sprintf("%s/v1/get/princess", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	// second call fills the cache
	_, err = http.Get(fmt.Sprintf("%s/v1/get/han", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	// third call evicts first call (princess leia)
	_, err = http.Get(fmt.Sprintf("%s/v1/get/malcolm", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	// fourth call for princess leia causes redis hit
	_, err = http.Get(fmt.Sprintf("%s/v1/get/princess", s.URL))
	if err != nil {
		t.Fatal(err)
	}

	expectedStack := []string{
		"GET$princess",
		"GET$han",
		"GET$malcolm",
		"GET$princess",
	}

	for i, v := range expectedStack {
		if v != dataSource.callList[i] {
			t.Fatalf("Expected %s but stack had %s", v, dataSource.callList[i])
		}
	}
}

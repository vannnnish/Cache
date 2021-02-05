package test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var url = "http://localhost:5837/v1/"

func Test_putEntry(t *testing.T) {
	body := strings.NewReader("value")

	request, err := http.NewRequest("PUT", url+"cache/key", body)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp)
}

func Test_getEntry(t *testing.T) {
	request, err := http.NewRequest("GET", url+"cache/key", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(value))
	fmt.Println(resp)

}

func Test_getCacheStatus(t *testing.T) {
	request, err := http.NewRequest("GET", url+"status", nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	value, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(value))
}

func Test_deleteEntry(t *testing.T) {
	request, err := http.NewRequest("DELETE", url+"cache/key", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp)
}

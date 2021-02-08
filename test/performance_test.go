package test

import (
	"cache/services"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHttpServer(t *testing.T) {
	writeTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		request, err := http.NewRequest("PUT", "http://localhost:5837/v1/cache"+data, strings.NewReader(data))
		if err != nil {
			t.Fatal(err)
		}
		resp, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
	})
	t.Logf("写入消耗时间为 %s", writeTime)
}

func TestTcpServer(t *testing.T) {
	client, err := services.NewTCPClient(":5837")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	writeTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		err = client.Set(data, []byte(data), 0)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Logf("写入时间消耗为 %s", writeTime)
	time.Sleep(3 * time.Second)

	readTime := testTask(func(no int) {
		data := strconv.Itoa(no)
		_, err := client.Get(data)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Logf("读取时间消耗为: %s", readTime)
}


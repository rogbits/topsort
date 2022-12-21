package lib

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type TestHandler struct{}

func (t *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v := r.Header.Get("x-test")
	_, err := w.Write([]byte(v))
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func TestNewHttpServer(t *testing.T) {
	api := new(TestHandler)
	s := NewHttpServer(nil, api)
	if s.Server.Handler != api {
		t.FailNow()
	}
}

func TestServer_Start_Stop(t *testing.T) {
	testPort := 9998
	logger := NewLogger()
	logger.TestMode = true
	api := new(TestHandler)
	s := NewHttpServer(logger, api)

	addr := fmt.Sprintf(":%d", testPort)
	url := fmt.Sprintf("http://localhost:%d", testPort)
	go func() {
		_ = s.Start(addr)
	}()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte("test123")
	req.Header.Add("x-test", string(expected))
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	text, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(text, expected) != 0 {
		t.Fatal("text != expected")
	}

	err = s.Stop()
	if err != nil {
		t.Fatal(err)
	}

	l1Actual := fmt.Sprintf("%s %s", logger.TestMessages[0], logger.TestMessages[1])
	l1Expected := fmt.Sprintf("starting server on :%d", testPort)
	l2Actual := logger.TestMessages[2]
	l2Expected := "stopping server"
	if l1Actual != l1Expected {
		t.Fatal("l1Actual != l1Expected")
	}
	if l2Actual != l2Expected {
		t.Fatal("l2Actual != l2Expected")
	}
}

func TestServer_Start_Error(t *testing.T) {
	logger := NewLogger()
	logger.TestMode = true

	api := new(TestHandler)
	s := NewHttpServer(logger, api)
	err := s.Start(":123456")
	expectedMsg := "listen tcp: address 123456: invalid port"
	if err == nil && err.Error() != expectedMsg {
		t.Fatal("expecting error due to port")
	}
	_ = s.Stop()
}

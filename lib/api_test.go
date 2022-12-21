package lib

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestApi_ServeHTTP_UriParseFailure(t *testing.T) {
	req := httptest.NewRequest("GET", "/willRemove", nil)
	req.RequestURI = ""
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	if res.Body.String() != "bad URI\n" {
		t.Fatal("expecting message re: bad URI")
	}
}

func TestApi_ServeHTTP_Health(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatal("expecting 200")
	}
	actual, err := io.ReadAll(res.Body)
	expected := []byte("ok\n")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Compare(actual, expected) != 0 {
		t.Fatal("actual != expected")
	}
}

func TestApi_ServeHTTP_Health_MethodError(t *testing.T) {
	req := httptest.NewRequest("POST", "/health", nil)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	expected := "unsupported method for /health\n"
	if res.Body.String() != expected {
		t.Fatal("expecting msg", expected)
	}
}

func TestApi_ServeHTTP_Sort_MethodError(t *testing.T) {
	req := httptest.NewRequest("GET", "/sort", nil)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	expected := "unsupported method for /sort\n"
	if res.Body.String() != expected {
		t.Fatal("expecting msg", expected)
	}
}

func TestApi_ServeHTTP_Sort_JsonDecodeError(t *testing.T) {
	reader := bytes.NewReader([]byte("{badjson}"))
	req := httptest.NewRequest("POST", "/sort", reader)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	expected := "error decoding edges input\n"
	if res.Body.String() != expected {
		t.Fatal("expecting msg", expected)
	}
}

func TestApi_ServeHTTP_Sort_UndirectedGraph(t *testing.T) {
	edges := [][]string{
		{"a", "b"},
		{"b", "a"},
	}
	marshalled, err := json.Marshal(edges)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(marshalled)
	req := httptest.NewRequest("POST", "/sort", reader)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	expected := "seeing an undirected graph\n"
	if res.Body.String() != expected {
		t.Fatal("expecting msg", expected)
	}
}

func TestApi_ServeHTTP_Sort_EmptyGraph(t *testing.T) {
	reader := bytes.NewReader([]byte("[]"))
	req := httptest.NewRequest("POST", "/sort", reader)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 400 {
		t.Fatal("expecting 400")
	}
	expected := "seeing empty graph\n"
	if res.Body.String() != expected {
		t.Fatal("expecting msg", expected)
	}
}

func TestApi_ServeHTTP_Sort_CycleDetection(t *testing.T) {
	edges := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "a"},
	}
	marshalled, err := json.Marshal(edges)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(marshalled)
	req := httptest.NewRequest("POST", "/sort", reader)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)

	if res.Code != 500 {
		t.Fatal("expecting 500")
	}
	expected := "cycle detected\n"
	if res.Body.String() != expected {
		t.Fatal("expecting", expected)
	}
}

func TestApi_ServeHTTP_Sort(t *testing.T) {
	body := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"d", "e"},
		{"e", "f"},
		{"f", "g"},
	}

	marshalled, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	reader := bytes.NewReader(marshalled)
	req := httptest.NewRequest("POST", "/sort", reader)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Fatal("expecting 200")
	}
	var payload []string
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&payload)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"a", "b", "c", "d", "e", "f", "g"}
	if strings.Join(payload, "") != strings.Join(expected, "") {
		t.Fatal("failed sort")
	}
}

func TestApi_ServeHTTP_404(t *testing.T) {
	req := httptest.NewRequest("GET", "/does-not-exist", nil)
	res := httptest.NewRecorder()

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	api.ServeHTTP(res, req)
	if res.Code != 404 {
		t.Fatal("expecting 404")
	}
}

func BenchmarkApi_ServeHTTP(b *testing.B) {
	edges := [][]string{
		{"a", "b"},
		{"b", "c"},
		{"c", "d"},
		{"d", "e"},
		{"e", "f"},
	}
	marshalled, err := json.Marshal(edges)
	if err != nil {
		b.Fatal(err)
	}

	logger := NewLogger()
	logger.TestMode = true
	api := NewApi(logger)
	expected := "abcdef"
	for i := 0; i < b.N; i++ {
		reader := bytes.NewReader(marshalled)
		req := httptest.NewRequest("POST", "/sort", reader)
		res := httptest.NewRecorder()

		api.ServeHTTP(res, req)
		if res.Code != 200 {
			b.Fatal("expecting 200")
		}
		var actual []string
		dec := json.NewDecoder(res.Body)
		_ = dec.Decode(&actual)
		if strings.Join(actual, "") != expected {
			b.Fatal("unexpected resp")
		}
	}
}

package goplay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestClient_Run(t *testing.T) {
	saveDelay := delay
	defer func() { delay = saveDelay }()

	events := []*Event{
		{Message: "out1", Kind: "stdout", Delay: 5},
		{Message: "out2", Kind: "stdout", Delay: 5},
		{Message: "err1", Kind: "stderr", Delay: 5},
	}

	delayCalledN := 0
	delay = func(d time.Duration) {
		delayCalledN++
		if d != 5 {
			t.Errorf("delay func got unexpected value: %v", d)
		}
	}
	defer func() {
		if got, want := delayCalledN, len(events); got != want {
			t.Errorf("delay func called %v times, want %v times", got, delayCalledN)
		}
	}()

	mu := http.NewServeMux()
	mu.HandleFunc("/compile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: got %v, want %v", r.Method, http.MethodPost)
		}
		if err := json.NewEncoder(w).Encode(&Response{Events: events}); err != nil {
			t.Error(err)
		}
	})
	ts := httptest.NewServer(mu)
	defer ts.Close()

	cli := &Client{
		BaseURL: ts.URL,
	}

	stdout, stderr := new(bytes.Buffer), new(bytes.Buffer)

	if err := cli.Run(bytes.NewReader([]byte("code")), stdout, stderr); err != nil {
		t.Fatal(err)
	}

	wantStdout := ""
	wantStderr := ""
	for _, e := range events {
		if e.Kind == "stdout" {
			wantStdout += e.Message
		} else {
			wantStderr += e.Message
		}
	}

	if got := stdout.String(); got != wantStdout {
		t.Errorf("Run() writes %v to stdout, want %v", got, wantStdout)
	}
	if got := stderr.String(); got != wantStderr {
		t.Errorf("Run() writes %v to stderr, want %v", got, wantStderr)
	}
}

func TestClient_Compile(t *testing.T) {
	wantResp := &Response{
		Errors: "",
		Events: []*Event{
			{Message: "test1", Kind: "stdout", Delay: 0},
			{Message: "test2", Kind: "stdout", Delay: 2},
		},
	}

	mu := http.NewServeMux()
	mu.HandleFunc("/compile", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: got %v, want %v", r.Method, http.MethodPost)
		}
		if err := json.NewEncoder(w).Encode(wantResp); err != nil {
			t.Error(err)
		}
	})
	ts := httptest.NewServer(mu)
	defer ts.Close()

	cli := &Client{
		BaseURL: ts.URL,
	}

	got, err := cli.Compile(bytes.NewReader([]byte("code")))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, wantResp) {
		t.Errorf("Compile(code) == %v, want %v", got, wantResp)
	}
}

func TestClient_Share(t *testing.T) {
	mu := http.NewServeMux()
	mu.HandleFunc("/share", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method: got %v, want %v", r.Method, http.MethodPost)
		}
		fmt.Fprint(w, "xxx") // token
	})
	ts := httptest.NewServer(mu)
	defer ts.Close()

	cli := &Client{
		BaseURL: ts.URL,
	}

	link, err := cli.Share(bytes.NewReader([]byte("test")))
	if err != nil {
		t.Fatal(err)
	}
	if want := fmt.Sprintf("%s/p/xxx", ts.URL); link != want {
		t.Errorf("Share(code) == %v, want %v", link, want)
	}
}

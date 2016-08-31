// Package goplay provides The Go Playground (https://play.golang.org/) client
package goplay

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var (
	baseURL         = "https://play.golang.org"
	shareEndpoint   = baseURL + "/share"
	compileEndpoint = baseURL + "/compile"
)

var httpClient = &http.Client{}

// Run runs go code on playground and output stdout/stderr.
func Run(code io.Reader, stdout io.Writer, stderr io.Writer) error {
	resp, err := Compile(code)
	if err != nil {
		return err
	}
	if resp.Errors != "" {
		return errors.New(resp.Errors)
	}
	for _, event := range resp.Events {
		time.Sleep(event.Delay)
		w := stderr
		if event.Kind == "stdout" {
			w = stdout
		}
		fmt.Fprint(w, event.Message)
	}
	return nil
}

// Share creates go playground share link.
func Share(code io.Reader) (string, error) {
	req, err := http.NewRequest("POST", shareEndpoint, code)
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/p/%s", baseURL, string(b)), nil
}

// Compile compiles code on the go playground.
func Compile(code io.Reader) (*Response, error) {
	b, err := ioutil.ReadAll(code)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Set("version", "2")
	v.Set("body", string(b))
	resp, err := httpClient.PostForm(compileEndpoint, v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var r Response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// Response represensts response type of /compile.
// Licence: Copyright (c) 2014 The Go Authors. All rights reserved.
// https://github.com/golang/playground/blob/816964eae74f7612221c13ab73f2a8021c581010/sandbox/sandbox.go#L35-L38
type Response struct {
	Errors string
	Events []Event
}

// Event represensts event of /compile result.
// Licence: Copyright (c) 2014 The Go Authors. All rights reserved.
// https://github.com/golang/playground/blob/816964eae74f7612221c13ab73f2a8021c581010/sandbox/play.go#L76-L80
type Event struct {
	Message string
	Kind    string        // "stdout" or "stderr"
	Delay   time.Duration // time to wait before printing Message
}

package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func Download(filepath, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
		if err != nil {
			Removes(filepath)
		}
	}()

	var resp *http.Response
	transport := &http.Transport{
		TLSHandshakeTimeout: 20 * time.Second,
	}

	client := http.Client{Transport: transport}
	resp, err = client.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	defer func() { _ = resp.Body.Close() }()
	_, err = io.Copy(out, resp.Body)
	return err
}

func Progress() (chan<- struct{}, <-chan struct{}) {
	var ch, done = make(chan struct{}), make(chan struct{})

	var fn = func() {
		now := time.Now()
		fmt.Printf("Start at %s\n", FormatTime(now))
		defer func() {
			fmt.Printf("Finish at %s used: %fs\n\n", FormatTime(time.Now()), time.Since(now).Seconds())
			done <- struct{}{}
		}()
		ticker := time.NewTicker(time.Second * 3)
		for {
			select {
			case <-ticker.C:
				fmt.Print(".")
			case <-ch:
				fmt.Println()
				return
			}
		}
	}

	go fn()
	return ch, done
}

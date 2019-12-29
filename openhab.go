package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	defaultOpenhabBase = "http://openhab:8080/rest"
	togglePath         = "items/Fan"
	statePath          = "items/Fan/state"
	eventsPath         = "events?topics=smarthome/items/Fan/statechanged"
)

type openhab struct {
	base string
}

func (oh *openhab) url(p string) string {
	return fmt.Sprintf("%s/%s", oh.base, p)
}

func (oh *openhab) toggleState(ctx context.Context) error {
	_, err := postContext(ctx, oh.url(togglePath), []byte("TOGGLE"))
	return err
}

func (oh *openhab) currentState(ctx context.Context) (fanState, error) {
	resp, err := getContext(ctx, oh.url(statePath))
	if err != nil {
		return fsNone, err
	}

	defer resp.Body.Close()

	b := new(bytes.Buffer)
	if _, err := io.Copy(b, resp.Body); err != nil {
		return fsNone, err
	}

	switch strings.TrimSpace(b.String()) {
	case "OFF":
		return fsOff, nil
	case "ON":
		return fsOn, nil
	default:
		return fsNone, nil
	}
}

func (oh *openhab) events(ctx context.Context, ch chan<- fanState) error {
	if ch == nil {
		return errors.New("nil channel")
	}
	defer close(ch)

	resp, err := getContext(ctx, oh.url(eventsPath))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		txt := s.Text()

		if !strings.HasPrefix(txt, "data: ") {
			continue
		}

		var evt event

		if err := json.Unmarshal([]byte(strings.TrimPrefix(txt, "data: ")), &evt); err != nil {
			fmt.Fprintf(os.Stderr, "JSON ERROR: %v\n", err)
			continue
		}

		if evt.Type != "ItemStateChangedEvent" || evt.Topic != "smarthome/items/Fan/statechanged" || evt.Payload == nil {
			continue
		}

		var fs fanState
		switch evt.Payload.Value {
		case "OFF":
			fs = fsOff
		case "ON":
			fs = fsOn
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case ch <- fs:
		}
	}

	return s.Err()
}

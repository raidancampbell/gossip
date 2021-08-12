package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

var host string
var formatter logrus.JSONFormatter

func init() {
	h, err := os.Hostname()
	if err != nil {
		h = "UNKNOWN"
	}
	host = h
	formatter = logrus.JSONFormatter{}
}

type SplunkHook struct {
	Client       *http.Client
	Endpoint     string
	Token        string
	FlushFreq    time.Duration
	MaxFlushSize int

	msg chan *SplunkLog
}

func NewSplunkHook(client *http.Client, endpoint string, token string, flushFreq time.Duration, maxFlushSize int) logrus.Hook {
	hook := &SplunkHook{
		Client:       client,
		Endpoint:     endpoint,
		Token:        token,
		FlushFreq:    flushFreq,
		MaxFlushSize: maxFlushSize,
	}

	// prevent a busy loop if the user gives bad data
	if hook.FlushFreq == 0 {
		hook.FlushFreq = 1 * time.Second
	}
	hook.msg = make(chan *SplunkLog, hook.MaxFlushSize)
	go hook.manageBuffer()
	return hook
}

func (s *SplunkHook) manageBuffer() {
	ticker := time.NewTicker(s.FlushFreq)

	buf := make([]*SplunkLog, 0)
	flush := func() {
		go s.doSend(buf)
		buf = make([]*SplunkLog, 0)
	}
	for {
		select {
		case <-ticker.C:
			if len(buf) > 0 {
				flush()
			}
		case d := <-s.msg:
			buf = append(buf, d)
			if len(buf) >= s.MaxFlushSize {
				flush()
			}
		}
	}
}

// doSend is synchronous with the actual HTTP send
// errors are all ignored.
func (s *SplunkHook) doSend(logs []*SplunkLog) {
	_b := []byte{}
	outputBuf := bytes.NewBuffer(_b)
	encoder := json.NewEncoder(outputBuf)
	for i := range logs {
		_ = encoder.Encode(logs[i])
	}

	req, _ := http.NewRequest(http.MethodPost, s.Endpoint, outputBuf)

	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", s.Token))
	res, err := s.Client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	if false { // dumb, but here for logging
		response, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(response))
	}
	_ = res.Body.Close()
}

func (s *SplunkHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (s *SplunkHook) Fire(e *logrus.Entry) error {
	b, err := formatter.Format(e)

	s.msg <- &SplunkLog{
		Time:  e.Time.UnixNano(),
		Host:  host,
		Event: b,
	}

	return err
}

type SplunkLog struct {
	Time  int64        `json:"time"`
	Host  string       `json:"host"`
	Event eventWrapper `json:"event"`
}
type eventWrapper []byte

func (e eventWrapper) MarshalJSON() ([]byte, error) {
	return e, nil
}

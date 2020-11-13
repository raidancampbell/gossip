package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"os"
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
	Client *http.Client
	Endpoint string
	Token string
}

func (s *SplunkHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

type SplunkLog struct {
	Time int64 `json:"time"`
	Host string `json:"host"`
	Event eventWrapper `json:"event"`
}
type eventWrapper []byte

func (e eventWrapper) MarshalJSON() ([]byte, error) {
	return e, nil
}

func (s *SplunkHook) Fire(e *logrus.Entry) error {
	b, err := formatter.Format(e)
	if err != nil {
		return err
	}
	b, err = json.Marshal(SplunkLog{
		Time: e.Time.UnixNano(),
		Host: host,
		Event: b,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, s.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Splunk %s", s.Token))
	go func() {
		res, err := s.Client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer res.Body.Close()
		if false { // dumb, but here for logging
			response, _ := httputil.DumpResponse(res, true)
			fmt.Println(string(response))
		}
	}()
	return nil
}

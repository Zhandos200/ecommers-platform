package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetReportCaller(true)
	Log.SetLevel(logrus.InfoLevel)
	Log.Out = &LokiHook{}
}

type LokiHook struct{}

func (hook *LokiHook) Write(p []byte) (n int, err error) {
	payload := map[string]interface{}{
		"streams": []map[string]interface{}{
			{
				"stream": map[string]string{
					"level":   "info",
					"service": "consumer-service",
				},
				"values": [][]string{
					{fmt.Sprintf("%d", time.Now().UnixNano()), string(p)},
				},
			},
		},
	}

	data, _ := json.Marshal(payload)
	resp, err := http.Post("http://loki:3100/loki/api/v1/push", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return len(p), nil
}

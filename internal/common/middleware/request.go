package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func RequestLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestIn(c, l)
		defer requestOut(c, l)
		c.Next()
	}
}

func requestOut(c *gin.Context, l *logrus.Entry) {
	response, _ := c.Get("response")
	var unMarshalledResp interface{}
	_ = json.Unmarshal(response.([]byte), &unMarshalledResp)
	start, _ := c.Get("start_time")
	startTime := start.(time.Time)
	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"proc_time_ms": time.Since(startTime).Milliseconds(),
		"response":     unMarshalledResp,
	}).Info("__request_out")
}

func requestIn(c *gin.Context, l *logrus.Entry) {
	c.Set("start_time", time.Now())
	body := c.Request.Body
	bodyBytes, _ := io.ReadAll(body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var jsonCompact bytes.Buffer
	_ = json.Compact(&jsonCompact, bodyBytes)
	l.WithContext(c.Request.Context()).WithFields(logrus.Fields{
		"start": time.Now().Unix(),
		"args":  jsonCompact.String(),
		"from":  c.RemoteIP(),
		"uri":   c.Request.RequestURI,
	}).Info("__request_in")
}

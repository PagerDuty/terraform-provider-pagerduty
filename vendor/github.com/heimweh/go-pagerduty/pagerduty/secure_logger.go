package pagerduty

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	secureLogRequestHeading = `[SECURE] PagerDuty API Request Details:
---[ REQUEST ]---------------------------------------`
	secureLogResponseHeading = `[SECURE] PagerDuty API Response Details:
---[ RESPONSE ]--------------------------------------`
	secureLogBottomDelimiter = `-----------------------------------------------------`
	obscuredLogTag           = `<OBSCURED>`
)

type secureLogger struct {
	logger         *log.Logger
	headersContent string
	bodyContent    string
	logsContent    string
	canLog         bool
}

func (l *secureLogger) handleHeadersLogsContent(h http.Header) {
	l.headersContent = ""
	headers := make(http.Header)
	for k, v := range h {
		headers[k] = v
	}

	if _, ok := headers["Authorization"]; ok {
		authHeader := headers["Authorization"][0]
		last4AuthChars := authHeader
		if len(authHeader) > 4 {
			last4AuthChars = authHeader[len(authHeader)-4:]
		}
		headers["Authorization"] = []string{fmt.Sprintf("%s%s", obscuredLogTag, last4AuthChars)}
	}

	for k, v := range headers {
		h := fmt.Sprintf("%s: %s", k, strings.Join(v, ";"))
		l.headersContent = fmt.Sprintf("%s%s\n", l.headersContent, h)
	}
}

func (l *secureLogger) handleBodyLogsContent(body io.ReadCloser) io.ReadCloser {
	l.bodyContent = ""
	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err != nil {
			log.Printf("[ERROR] Error reading body: %v\n", err)
			return body
		}

		var jsonObj map[string]interface{}
		err = json.Unmarshal(bodyBytes, &jsonObj)
		if err != nil {
			l.bodyContent = fmt.Sprintf("%s\n", string(bodyBytes))
		} else {
			prettyBody, err := json.MarshalIndent(jsonObj, "", " ")
			if err != nil {
				log.Printf("[ERROR] Error pretty-printing body: %v\n", err)
			} else {
				l.bodyContent = fmt.Sprintf("%s\n", prettyBody)
			}
		}

		body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return body
}

func (l *secureLogger) putTogetherLogsContent(logsContent *string, heading string) {
	content := *logsContent
	content = fmt.Sprintf("%s\n%s\n%s\n", heading, content, l.headersContent)
	if l.bodyContent != "" {
		content = fmt.Sprintf("%s%s\n", content, l.bodyContent)
	}

	content = fmt.Sprintf("%s%s", content, secureLogBottomDelimiter)
	*logsContent = content
}

func (l *secureLogger) LogReq(req *http.Request) {
	if !l.canLog {
		return
	}

	logsContent := fmt.Sprintf("%s %s %s", req.Method, req.URL.Path, req.Proto)
	l.handleHeadersLogsContent(req.Header)
	req.Body = l.handleBodyLogsContent(req.Body)
	l.putTogetherLogsContent(&logsContent, secureLogRequestHeading)

	l.logger.Print(logsContent)
}

func (l *secureLogger) LogRes(res *http.Response) {
	if !l.canLog {
		return
	}

	logsContent := fmt.Sprintf("%s %d %s", res.Proto, res.StatusCode, res.Status)
	l.handleHeadersLogsContent(res.Header)
	res.Body = l.handleBodyLogsContent(res.Body)
	l.putTogetherLogsContent(&logsContent, secureLogResponseHeading)

	l.logger.Print(logsContent)
}

func (l *secureLogger) SetCanLog(flag bool) {
	l.canLog = flag
}

func newSecureLogger() *secureLogger {
	pdLogFlag := os.Getenv("TF_LOG_PROVIDER_PAGERDUTY")
	pdLogFlag = strings.ToUpper(pdLogFlag)
	tfLogFlag := os.Getenv("TF_LOG")
	tfLogFlag = strings.ToUpper(tfLogFlag)

	secLogger := secureLogger{
		logger: log.Default(),
		canLog: tfLogFlag == "INFO" && pdLogFlag == "SECURE",
	}
	secLogger.logger.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	return &secLogger
}

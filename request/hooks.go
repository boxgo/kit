package request

import (
	"github.com/BiteBit/gorequest"
	"github.com/boxgo/kit/logger"
)

func logBefore(agent *gorequest.SuperAgent) {
	if GlobalOptions.UserAgent != "" {
		agent.AppendHeader("user-agent", GlobalOptions.UserAgent)
	}

	if GlobalOptions.ShowLog {
		curl, _ := agent.AsCurlCommand()
		logger.Default.Infow("request_start", "curl", curl)
	}
}

func logAfter(agent *gorequest.SuperAgent, resp gorequest.Response, body []byte, errs []error) {
	if GlobalOptions.ShowLog {
		curl, _ := agent.AsCurlCommand()
		logger.Default.Infow("request_end", "curl", curl, "errs", errs, "resp.status", resp.StatusCode, "body", string(body[:]), "resp.header", resp.Header)
	}
}

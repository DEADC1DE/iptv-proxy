package server

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func channelAwareFormatter(param gin.LogFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency.Truncate(time.Second)
	}

	channelName := "n/a"
	if param.Keys != nil {
		if rawName, ok := param.Keys[channelNameContextKey]; ok {
			if name, ok := rawName.(string); ok && name != "" {
				channelName = name
			}
		}
	}
	channelFragment := fmt.Sprintf(" | channel=%s", channelName)

	return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v%s\n%s",
		param.TimeStamp.Format("2006/01/02 - 15:04:05"),
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		channelFragment,
		param.ErrorMessage,
	)
}

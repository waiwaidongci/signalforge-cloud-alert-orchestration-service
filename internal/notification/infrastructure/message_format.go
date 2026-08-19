package infrastructure

import "strings"

func ChannelLabel(channel string) string {
	switch channel {
	case "email":
		return "邮件"
	case "webhook":
		return "回调"
	case "log":
		return "日志"
	default:
		return "未知通道"
	}
}

func DestinationSummary(destination string) string {
	if strings.TrimSpace(destination) == "" {
		return "未配置"
	}
	if len(destination) > 24 {
		return destination[:24]
	}
	return destination
}

func BuildFailureMessage(channel string, destination string) string {
	label := ChannelLabel(channel)
	if label == "未知通道" {
		return label
	}
	return strings.TrimSpace(label + " " + DestinationSummary(destination))
}

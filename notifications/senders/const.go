package senders

type Send func(msg, dst string, status MessageStatus) error

type MessageStatus string

var (
	Success       MessageStatus = "SUCCESS"
	Failure       MessageStatus = "FAILURE"
	Timeout       MessageStatus = "TIMEOUT"
	InternalError MessageStatus = "INTERNAL_ERROR"
)

var StatusColor = map[MessageStatus]string{
	Success:       "#34A853", // green
	Failure:       "#EA4335", // red
	Timeout:       "#FBBC05", // yellow
	InternalError: "#EA4335", // red
}

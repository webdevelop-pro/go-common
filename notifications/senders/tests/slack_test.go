package senders

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdevelop-pro/go-common/notifications/senders"
)

func TestSend_ToSlack_SuccessReturnNil(t *testing.T) {
	t.Parallel()

	var send senders.Send
	send = senders.SlackSender{
		Token: "xoxb-2523150300577-2495865134151-7yze2ClHQ1zLGILNvpGOYJOm",
	}.SendToSlack

	for status, _ := range senders.StatusColor {
		require.Nil(t, send("Test message - "+string(status), "tests", status))
	}
}

func TestSend_ToSlack_InvalidToken_ReturnError(t *testing.T) {
	t.Parallel()

	var send senders.Send
	send = senders.SlackSender{
		Token: "invalid_token",
	}.SendToSlack

	require.NotNil(t, send("Test message", "tests", senders.Success))
}

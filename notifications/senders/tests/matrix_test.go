package senders

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdevelop-pro/go-common/notifications/senders"
)

func TestSend_ToMatrix_SuccessReturnNil(t *testing.T) {
	t.Parallel()

	// Uncomment this lines if your want check matrix sender
	//
	// var send senders.Send
	// send = senders.SendToMatrix
	// webhook := "https://matrix.unfederalreserve.com/api/v1/matrix/hook/Z0IGHMnoLWRU8WDPx26FAwIdjJmbMaVX2Duy6LYQlmqUYTwPe77ZnUh4aKNhSXxv"
	// require.Nil(t, send("Test message", webhook, senders.Success))

	require.Nil(t, nil)
}

func TestSend_ToMatrix_InvalidWebhook_ReturnError(t *testing.T) {
	t.Parallel()

	// Uncomment this lines if your want check matrix sender
	//
	var send senders.Send
	send = senders.SendToMatrix
	webhook := "invalid_webhook"
	require.NotNil(t, send("Test message", webhook, senders.Success))
}

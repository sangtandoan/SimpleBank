package email

import (
	"testing"

	"github.com/FrostJ143/simplebank/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := utils.LoadConfig("../..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPasswrod)

	subject := "A test mail"
	content := `
        <h1>Hello World<h1>
        <p>This is a test message<p>
    `
	to := []string{"16025@student.vgu.edu.vn"}
	attachFiles := []string{"../../app.env"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}

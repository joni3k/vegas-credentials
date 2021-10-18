package provider

import (
	"context"

	"github.com/aripalo/aws-mfa-credential-process/internal/data"
	"github.com/aripalo/aws-mfa-credential-process/internal/prompt"
)

func (t *TokenProvider) QueryGUI(ctx context.Context, d data.Provider) {
	var token Token
	var err error

	token.Provider = TOKEN_PROVIDER_GUI_DIALOG_PROMPT

	value, err := prompt.Dialog(ctx, "Multifactor Authentication", "Enter TOPT MFA Token Code:")
	if err != nil {
		t.errorChan <- &err
	} else {
		token.Value = value
		t.tokenChan <- &token
	}
}

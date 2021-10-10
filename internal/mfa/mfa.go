package mfa

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aripalo/aws-mfa-credential-process/internal/profile"
	"github.com/aripalo/aws-mfa-credential-process/internal/utils"
)

// The default timeout of Yubikey operations
const YUBIKEY_TIMEOUT_SECONDS = 15

func GetTokenResult(config profile.Profile, hideArns bool) (Result, error) {
	resultChan := make(chan *Result, 1)
	errorChan := make(chan *error, 1)

	ctx, cancel := context.WithTimeout(context.Background(), YUBIKEY_TIMEOUT_SECONDS*time.Second)
	defer cancel()

	if hideArns == false {

		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_DEBUG, "👷 ", "Role", config.AssumeRoleArn))
		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_DEBUG, "🔒 ", "MFA", config.MfaSerial))

	}

	hasYubikey := (config.YubikeySerial != "" && config.YubikeyLabel != "")

	if hasYubikey {
		go getYubikeyToken(ctx, config.YubikeySerial, config.YubikeyLabel, resultChan, errorChan)
	}
	go getCliToken(ctx, resultChan, errorChan)

	if hasYubikey {
		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_IMPORTANT, "🔑 ", "MFA", "Touch Yubikey or enter TOPT MFA Token Code..."))
	} else {
		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_IMPORTANT, "🔑 ", "MFA", "Enter TOPT MFA Token Code..."))
	}

	utils.SafeLog(utils.FormatMessage(utils.COLOR_INPUT_EXPECTED, "🔑 ", "MFA", "> "))

	select {
	case i := <-resultChan:
		result := *i

		err := validateToken(result.Value)
		if err != nil {
			utils.SafeLogLn()
			utils.SafeLogLn(utils.FormatMessage(utils.COLOR_ERROR, "❌ ", "MFA", fmt.Sprintf("Invalid Token Code \"%s\" received via %s", result.Value, result.Provider)))
			return result, err
		}

		if result.Provider == TOKEN_PROVIDER_YUBIKEY {
			utils.SafeLogLn(result.Value)
		}
		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_IMPORTANT, "🔓 ", "MFA", fmt.Sprintf("Token Code \"%s\" received via %s", result.Value, result.Provider)))
		return result, nil
	case <-ctx.Done():
		utils.SafeLogLn()
		if ctx.Err() == context.DeadlineExceeded {
			utils.SafeLogLn(utils.FormatMessage(utils.COLOR_ERROR, "❌ ", "MFA", "Operation Timeout"))
			return Result{}, errors.New("MFA Operation Timeout")
		}

		utils.SafeLogLn(utils.FormatMessage(utils.COLOR_ERROR, "❌ ", "MFA", ctx.Err().Error()))
		return Result{}, ctx.Err()
	}
}

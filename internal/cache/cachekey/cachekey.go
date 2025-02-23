package cachekey

import (
	"encoding/json"
	"strings"

	"github.com/aripalo/vegas-credentials/internal/interfaces"
	"github.com/aripalo/vegas-credentials/internal/profile"
	"github.com/aripalo/vegas-credentials/internal/utils"
)

const Separator = "__"

// Get is responsible for creating a unique cache key for given profile configuration, therefore ensuring mutated profile configuration will not use previous cached data
func Get(a interfaces.AssumeCredentialProcess) (string, error) {
	f := a.GetFlags()
	p := a.GetProfile()

	configString, err := configToString(*p)
	if err != nil {
		return "", err
	}

	hash, err := utils.GenerateSHA1(configString)
	if err != nil {
		return "", err
	}

	var key strings.Builder
	key.WriteString(f.Profile)
	key.WriteString(Separator)
	key.WriteString(hash)

	return key.String(), err
}

// configToString convertts profile config into stringified JSON
func configToString(p profile.Profile) (string, error) {
	result, err := json.Marshal(p)
	return string(result), err
}

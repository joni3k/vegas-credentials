package response

import (
	"errors"
	"fmt"
	"time"

	"github.com/aripalo/vegas-credentials/internal/interfaces"
	"github.com/dustin/go-humanize"
)

// Validate ensures the response is of correct format
func (r *Response) Validate(a interfaces.AssumeCredentialProcess) error {

	if r.Version != AWS_CREDENTIAL_PROCESS_VERSION {
		return errors.New("Incorrect Version")
	}

	if r.AccessKeyID == "" {
		return errors.New("Missing AccessKeyID")
	}

	if r.SecretAccessKey == "" {
		return errors.New("Missing SecretAccessKey")
	}

	if r.SessionToken == "" {
		return errors.New("Missing SessionToken")
	}

	now := time.Now()

	if r.Expiration.Before(now) {
		return fmt.Errorf("Expired %s", humanize.RelTime(r.Expiration, now, "ago", "in future"))
	}

	return nil
}

// ValidateForMandatoryRefresh ensures response is within "mandatory refresh" duration as per BotoCore
// https://github.com/boto/botocore/blob/221ffa67a567df99ee78d7ae308c0e12d7eeeea7/botocore/credentials.py#L350-L355
func (r *Response) ValidateForMandatoryRefresh(a interfaces.AssumeCredentialProcess) error {

	f := a.GetFlags()

	if f.DisableMandatoryRefresh {
		return nil
	}

	now := time.Now()
	count := 10 * 60
	limit := now.Add(time.Duration(-count) * time.Second)

	if r.Expiration.Before(limit) {
		return fmt.Errorf("Mandatory refresh required because expiration in %s", humanize.Time(r.Expiration))
	}

	return nil
}

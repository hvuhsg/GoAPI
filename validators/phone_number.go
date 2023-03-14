package validators

import (
	"fmt"
	"regexp"

	"github.com/hvuhsg/goapi/request"
)

type VPhoneNumber struct {
	Prefix string
}

func (v VPhoneNumber) Validate(r *request.Request, paramName string) error {
	phoneNumber := r.GetString(paramName)

	// Regular expression for phone numbers with optional prefix and extension
	// The prefix is specified in the validator
	// Source: https://www.oreilly.com/library/view/regular-expressions-cookbook/9781449327453/ch04s02.html
	re := regexp.MustCompile(fmt.Sprintf(`^%s?\(?\d{3}\)?[- ]?\d{3}[- ]?\d{4}(?: *x\d+)?$`, v.Prefix))

	if !re.MatchString(phoneNumber) {
		return fmt.Errorf("parameter %s must be a valid phone number with prefix %s", paramName, v.Prefix)
	}

	return nil
}

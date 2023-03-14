package validators

import (
	"fmt"
	"net"

	"github.com/hvuhsg/goapi/request"
)

type VIPAddress struct{}

func (v VIPAddress) Validate(r *request.Request, paramName string) error {
	ipStr := r.GetString(paramName)

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("parameter %s must be a valid IP address", paramName)
	}

	return nil
}

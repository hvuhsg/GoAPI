package goapi

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

type Validator interface {
	Validate(r *Request, paramName string) error
	updateOpenAPISchema(schema *openapi3.Schema)
}

type VRequired struct{}

func (v VRequired) updateOpenAPISchema(schema *openapi3.Schema) {}
func (VRequired) Validate(r *Request, paramName string) error {
	_, ok := r.parameters[paramName]
	if !ok {
		return fmt.Errorf("parameter %s is required", paramName)
	}

	return nil
}

type VIsString struct{}

func (v VIsString) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "string" }
func (VIsString) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'string'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetString(paramName)
	}

	return nil
}

type VIsInt struct{}

func (v VIsInt) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "integer"
	schema.Format = "int64"
}
func (VIsInt) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'int'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetInt(paramName)
	}

	return nil
}

type VIsFloat struct{}

func (v VIsFloat) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Type = "number"
	schema.Format = "float"
}
func (VIsFloat) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'float'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetFloat(paramName)
	}

	return nil
}

type VIsBool struct{}

func (v VIsBool) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "boolean" }
func (VIsBool) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'bool'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetBool(paramName)
	}

	return nil
}

type VIsArray struct{}

func (v VIsArray) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "array" }
func (VIsArray) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'array'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetArray(paramName)
	}

	return nil
}

type VIsMap struct{}

func (v VIsMap) updateOpenAPISchema(schema *openapi3.Schema) { schema.Type = "object" }
func (VIsMap) Validate(r *Request, paramName string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("parameter %s must be of type 'map'", paramName)
		}
	}()

	_, ok := r.parameters[paramName]

	if ok {
		r.GetMap(paramName)
	}

	return nil
}

type VStringLength struct {
	Min int
	Max int
}

func (v VStringLength) updateOpenAPISchema(schema *openapi3.Schema) {
	VIsString{}.updateOpenAPISchema(schema)
	schema.MinLength = uint64(v.Min)
	maxL := uint64(v.Max)
	schema.MaxLength = &maxL
}
func (v VStringLength) Validate(r *Request, paramName string) error {
	vr := VIsString{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	strValue := r.GetString(paramName)

	if len(strValue) < v.Min || len(strValue) > v.Max {
		return fmt.Errorf("parameter %s length must be between %d and %d characters", paramName, v.Min, v.Max)
	}

	return nil
}

type VRange struct {
	Min float64
	Max float64
}

func (v VRange) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Min = &v.Min
	schema.Max = &v.Max
}
func (v VRange) Validate(r *Request, paramName string) error {
	vr := VIsFloat{}
	err := vr.Validate(r, paramName)
	if err != nil {
		return err
	}

	fVal := r.GetFloat(paramName)

	if fVal < v.Min || fVal > v.Max {
		return fmt.Errorf("parameter %s must be between %f and %f", paramName, v.Min, v.Max)
	}

	return nil
}

type VIsEmail struct{}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (v VIsEmail) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "email"
}
func (v VIsEmail) Validate(r *Request, paramName string) error {
	email := r.GetString(paramName)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("parameter %s must be a valid email address", paramName)
	}
	return nil
}

type VIsURL struct{}

func (v VIsURL) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "url"
}
func (v VIsURL) Validate(r *Request, paramName string) error {
	urlStr := r.GetString(paramName)

	// Parse the URL and check for any errors
	_, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid URL: %v", paramName, err)
	}

	return nil
}

type VIsUUID struct{}

var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

func (v VIsUUID) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "uuid"
}
func (v VIsUUID) Validate(r *Request, paramName string) error {
	uuidStr := r.GetString(paramName)

	// Check if the UUID string matches the expected format
	if !uuidRegex.MatchString(uuidStr) {
		return fmt.Errorf("parameter %s must be a valid UUID", paramName)
	}

	// Convert the UUID string to lowercase and check for any non-hexadecimal characters
	uuidStr = strings.ReplaceAll(strings.ToLower(uuidStr), "-", "")
	if _, err := fmt.Sscanf(uuidStr, "%x", &struct{}{}); err != nil {
		return fmt.Errorf("parameter %s must be a valid UUID", paramName)
	}

	return nil
}

type VIsDate struct {
	// The expected date format. This should be a string that specifies the expected format of the date, using the following standard Go date and time format codes: https://golang.org/pkg/time/#pkg-constants.
	// For example, to specify a date in the format "YYYY-MM-DD", you can set Format to "2006-01-02".
	Format string
}

func (v VIsDate) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "date"
}
func (v VIsDate) Validate(r *Request, paramName string) error {
	dateStr := r.GetString(paramName)

	// Parse the date string based on the specified format
	parsedDate, err := time.Parse(v.Format, dateStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid date in the format %s", paramName, v.Format)
	}

	// Check if the parsed date is valid
	if parsedDate.IsZero() {
		return fmt.Errorf("parameter %s must be a valid date in the format %s", paramName, v.Format)
	}

	return nil
}

type VRegex struct {
	Regex string
}

func (v VRegex) updateOpenAPISchema(schema *openapi3.Schema) {
	// Not applicable
}

func (v VRegex) Validate(r *Request, paramName string) error {
	str := r.GetString(paramName)

	// Compile the regular expression
	re, err := regexp.Compile(v.Regex)
	if err != nil {
		return fmt.Errorf("failed to compile regular expression: %v", err)
	}

	// Match the string against the regular expression
	if !re.MatchString(str) {
		return fmt.Errorf("parameter %s must match the regular expression %s", paramName, v.Regex)
	}

	return nil
}

type VIsTime struct {
	// The expected format of the time string. This should be a string that specifies the expected format of the time, using the following standard Go date and time format codes: https://golang.org/pkg/time/#pkg-constants.
	// For example, to specify a time in the format "15:04:05", you can set Format to "15:04:05".
	Format string
}

func (v VIsTime) updateOpenAPISchema(schema *openapi3.Schema) {
	schema.Format = "time"
}

func (v VIsTime) Validate(r *Request, paramName string) error {
	timeStr := r.GetString(paramName)

	// Parse the time string based on the specified format
	parsedTime, err := time.Parse(v.Format, timeStr)
	if err != nil {
		return fmt.Errorf("parameter %s must be a valid time in the format %s", paramName, v.Format)
	}

	// Check if the parsed time is valid
	if parsedTime.IsZero() {
		return fmt.Errorf("parameter %s must be a valid time in the format %s", paramName, v.Format)
	}

	return nil
}

type VPassword struct {
	MinLength      int
	MaxLength      int
	RequireSymbols bool
	RequireNumbers bool
	RequireUpper   bool
}

func (v VPassword) updateOpenAPISchema(schema *openapi3.Schema) {
	// Not applicable
}

func (v VPassword) Validate(r *Request, paramName string) error {
	password := r.GetString(paramName)

	// Check length requirements
	if len(password) < v.MinLength || len(password) > v.MaxLength {
		return fmt.Errorf("parameter %s must have length between %d and %d", paramName, v.MinLength, v.MaxLength)
	}

	// Check symbol requirements
	if v.RequireSymbols && !regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\{\}\[\]:;<>\?,\.\|\\\/~]`).MatchString(password) {
		return errors.New("password must contain at least one symbol")
	}

	// Check number requirements
	if v.RequireNumbers && !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	// Check uppercase letter requirements
	if v.RequireUpper && !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	return nil
}

type VIPAddress struct{}

func (v VIPAddress) Validate(r *Request, paramName string) error {
	ipStr := r.GetString(paramName)

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("parameter %s must be a valid IP address", paramName)
	}

	return nil
}

type VPhoneNumber struct {
	Prefix string
}

func (v VPhoneNumber) Validate(r *Request, paramName string) error {
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

package hotspot

import (
	"fmt"
	"math/rand"
	"time"
)

// GenerateVouchers generates batch of vouchers
func (c *Client) GenerateVouchers(gen *VoucherGenerator) (*VoucherResult, error) {
	// Validate
	if gen.Profile == "" {
		return nil, NewError("generate vouchers", fmt.Errorf("profile required"))
	}
	if gen.Quantity <= 0 {
		return nil, NewError("generate vouchers", fmt.Errorf("quantity must be > 0"))
	}

	// Set defaults
	if gen.Charset == "" {
		gen.Charset = DefaultCharset
	}
	if gen.LengthUsername < MinUsername {
		gen.LengthUsername = 8
	}
	if gen.LengthPassword < MinPassword {
		gen.LengthPassword = 8
	}

	result := &VoucherResult{
		Vouchers: make([]User, 0, gen.Quantity),
		Errors:   make([]string, 0),
	}

	// Generate vouchers
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < gen.Quantity; i++ {
		// Generate username
		username := gen.Prefix + "-" + GenerateRandomString(gen.LengthUsername, gen.Charset)

		// Generate password (same as username for voucher mode)
		password := GenerateRandomString(gen.LengthPassword, gen.Charset)

		user := User{
			Name:            username,
			Password:        password,
			Profile:         gen.Profile,
			Comment:         "vc-" + gen.Prefix,
			LimitUptime:     gen.TimeLimit,
			LimitBytesTotal: gen.DataLimit,
			Disabled:        false,
			Server:          "all",
		}

		// Create user on RouterOS
		err := c.CreateUser(&user)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: %v", username, err))
		} else {
			result.Success++
			result.Vouchers = append(result.Vouchers, user)
		}
	}

	return result, nil
}

// GenerateUserPasswordMode generates user-password mode (not voucher)
func (c *Client) GenerateUserPasswordMode(gen *VoucherGenerator) (*VoucherResult, error) {
	// Validate
	if gen.Profile == "" {
		return nil, NewError("generate users", fmt.Errorf("profile required"))
	}
	if gen.Quantity <= 0 {
		return nil, NewError("generate users", fmt.Errorf("quantity must be > 0"))
	}

	// Set defaults
	if gen.Charset == "" {
		gen.Charset = DefaultCharset
	}
	if gen.LengthUsername < MinUsername {
		gen.LengthUsername = 8
	}
	if gen.LengthPassword < MinPassword {
		gen.LengthPassword = 8
	}

	result := &VoucherResult{
		Vouchers: make([]User, 0, gen.Quantity),
		Errors:   make([]string, 0),
	}

	// Generate users
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < gen.Quantity; i++ {
		username := gen.Prefix + "-" + GenerateRandomString(gen.LengthUsername, gen.Charset)
		password := GenerateRandomString(gen.LengthPassword, gen.Charset)

		user := User{
			Name:            username,
			Password:        password,
			Profile:         gen.Profile,
			Comment:         "up-" + gen.Prefix,
			LimitUptime:     gen.TimeLimit,
			LimitBytesTotal: gen.DataLimit,
			Disabled:        false,
			Server:          "all",
		}

		err := c.CreateUser(&user)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: %v", username, err))
		} else {
			result.Success++
			result.Vouchers = append(result.Vouchers, user)
		}
	}

	return result, nil
}

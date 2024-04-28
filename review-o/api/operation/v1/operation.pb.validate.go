// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/operation/v1/operation.proto

package v1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on AuditReviewRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AuditReviewRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuditReviewRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuditReviewRequestMultiError, or nil if none found.
func (m *AuditReviewRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *AuditReviewRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetReviewID() <= 0 {
		err := AuditReviewRequestValidationError{
			field:  "ReviewID",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetStatus() <= 0 {
		err := AuditReviewRequestValidationError{
			field:  "Status",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetOpUser()) < 2 {
		err := AuditReviewRequestValidationError{
			field:  "OpUser",
			reason: "value length must be at least 2 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetOpReason()) < 2 {
		err := AuditReviewRequestValidationError{
			field:  "OpReason",
			reason: "value length must be at least 2 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.OpRemarks != nil {
		// no validation rules for OpRemarks
	}

	if len(errors) > 0 {
		return AuditReviewRequestMultiError(errors)
	}

	return nil
}

// AuditReviewRequestMultiError is an error wrapping multiple validation errors
// returned by AuditReviewRequest.ValidateAll() if the designated constraints
// aren't met.
type AuditReviewRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuditReviewRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuditReviewRequestMultiError) AllErrors() []error { return m }

// AuditReviewRequestValidationError is the validation error returned by
// AuditReviewRequest.Validate if the designated constraints aren't met.
type AuditReviewRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuditReviewRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuditReviewRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuditReviewRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuditReviewRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuditReviewRequestValidationError) ErrorName() string {
	return "AuditReviewRequestValidationError"
}

// Error satisfies the builtin error interface
func (e AuditReviewRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuditReviewRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuditReviewRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuditReviewRequestValidationError{}

// Validate checks the field values on AuditReviewReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *AuditReviewReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuditReviewReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuditReviewReplyMultiError, or nil if none found.
func (m *AuditReviewReply) ValidateAll() error {
	return m.validate(true)
}

func (m *AuditReviewReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ReviewID

	// no validation rules for Status

	if len(errors) > 0 {
		return AuditReviewReplyMultiError(errors)
	}

	return nil
}

// AuditReviewReplyMultiError is an error wrapping multiple validation errors
// returned by AuditReviewReply.ValidateAll() if the designated constraints
// aren't met.
type AuditReviewReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuditReviewReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuditReviewReplyMultiError) AllErrors() []error { return m }

// AuditReviewReplyValidationError is the validation error returned by
// AuditReviewReply.Validate if the designated constraints aren't met.
type AuditReviewReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuditReviewReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuditReviewReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuditReviewReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuditReviewReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuditReviewReplyValidationError) ErrorName() string { return "AuditReviewReplyValidationError" }

// Error satisfies the builtin error interface
func (e AuditReviewReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuditReviewReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuditReviewReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuditReviewReplyValidationError{}

// Validate checks the field values on AuditAppealRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *AuditAppealRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuditAppealRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuditAppealRequestMultiError, or nil if none found.
func (m *AuditAppealRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *AuditAppealRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if m.GetAppealID() <= 0 {
		err := AuditAppealRequestValidationError{
			field:  "AppealID",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetReviewID() <= 0 {
		err := AuditAppealRequestValidationError{
			field:  "ReviewID",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.GetStatus() <= 0 {
		err := AuditAppealRequestValidationError{
			field:  "Status",
			reason: "value must be greater than 0",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetOpUser()) < 2 {
		err := AuditAppealRequestValidationError{
			field:  "OpUser",
			reason: "value length must be at least 2 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if m.OpRemarks != nil {
		// no validation rules for OpRemarks
	}

	if len(errors) > 0 {
		return AuditAppealRequestMultiError(errors)
	}

	return nil
}

// AuditAppealRequestMultiError is an error wrapping multiple validation errors
// returned by AuditAppealRequest.ValidateAll() if the designated constraints
// aren't met.
type AuditAppealRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuditAppealRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuditAppealRequestMultiError) AllErrors() []error { return m }

// AuditAppealRequestValidationError is the validation error returned by
// AuditAppealRequest.Validate if the designated constraints aren't met.
type AuditAppealRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuditAppealRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuditAppealRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuditAppealRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuditAppealRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuditAppealRequestValidationError) ErrorName() string {
	return "AuditAppealRequestValidationError"
}

// Error satisfies the builtin error interface
func (e AuditAppealRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuditAppealRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuditAppealRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuditAppealRequestValidationError{}

// Validate checks the field values on AuditAppealReply with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *AuditAppealReply) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuditAppealReply with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuditAppealReplyMultiError, or nil if none found.
func (m *AuditAppealReply) ValidateAll() error {
	return m.validate(true)
}

func (m *AuditAppealReply) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AppealID

	// no validation rules for Status

	if len(errors) > 0 {
		return AuditAppealReplyMultiError(errors)
	}

	return nil
}

// AuditAppealReplyMultiError is an error wrapping multiple validation errors
// returned by AuditAppealReply.ValidateAll() if the designated constraints
// aren't met.
type AuditAppealReplyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuditAppealReplyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuditAppealReplyMultiError) AllErrors() []error { return m }

// AuditAppealReplyValidationError is the validation error returned by
// AuditAppealReply.Validate if the designated constraints aren't met.
type AuditAppealReplyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuditAppealReplyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuditAppealReplyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuditAppealReplyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuditAppealReplyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuditAppealReplyValidationError) ErrorName() string { return "AuditAppealReplyValidationError" }

// Error satisfies the builtin error interface
func (e AuditAppealReplyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuditAppealReply.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuditAppealReplyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuditAppealReplyValidationError{}

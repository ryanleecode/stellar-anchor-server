package internal

import (
	"regexp"
)

type BadRequestError struct {
	message string
}

func (e *BadRequestError) Error() string {
	return e.message
}

type ValidationError struct {
	errors interface{}
}

func (e *ValidationError) Error() string {
	return "validation error"
}

var horizonBadRequestRegex = regexp.MustCompile("bad_request")

type InvalidAccountID struct {
	message string
}

func (e *InvalidAccountID) Error() string {
	return e.message
}

type TransactionSourceAccountDoesntMatchAnchorPublicKey struct {
	Message string `json:"message"`
}

func (e *TransactionSourceAccountDoesntMatchAnchorPublicKey) Error() string {
	return e.Message
}

func NewTransactionSourceAccountDoesntMatchAnchorPublicKey(
	message string,
) *TransactionSourceAccountDoesntMatchAnchorPublicKey {
	return &TransactionSourceAccountDoesntMatchAnchorPublicKey{message}
}

type TransactionIsMissingTimeBounds struct {
	Message string `json:"message"`
}

func (e *TransactionIsMissingTimeBounds) Error() string {
	return e.Message
}

func NewTransactionIsMissingTimeBounds(
	message string,
) *TransactionIsMissingTimeBounds {
	return &TransactionIsMissingTimeBounds{message}
}

type TransactionChallengeExpired struct {
	Message string `json:"message"`
}

func (e *TransactionChallengeExpired) Error() string {
	return e.Message
}

func NewTransactionChallengeExpired(
	message string,
) *TransactionChallengeExpired {
	return &TransactionChallengeExpired{message}
}

type TransactionChallengeDoesNotHaveOnlyOneOperation struct {
	Message string `json:"message"`
}

func (e *TransactionChallengeDoesNotHaveOnlyOneOperation) Error() string {
	return e.Message
}

func NewTransactionChallengeDoesNotHaveOnlyOneOperation(
	message string,
) *TransactionChallengeDoesNotHaveOnlyOneOperation {
	return &TransactionChallengeDoesNotHaveOnlyOneOperation{message}
}

type TransactionChallengeIsNotAManageDataOperation struct {
	Message string `json:"message"`
}

func (e *TransactionChallengeIsNotAManageDataOperation) Error() string {
	return e.Message
}

func NewTransactionChallengeIsNotAManageDataOperation(
	message string,
) *TransactionChallengeIsNotAManageDataOperation {
	return &TransactionChallengeIsNotAManageDataOperation{message}
}

type TransactionOperationSourceAccountIsEmpty struct {
	Message string `json:"message"`
}

func (e *TransactionOperationSourceAccountIsEmpty) Error() string {
	return e.Message
}

func NewTransactionOperationSourceAccountIsEmpty(
	message string,
) *TransactionOperationSourceAccountIsEmpty {
	return &TransactionOperationSourceAccountIsEmpty{message}
}

type TransactionOperationsIsNil struct {
	Message string `json:"message"`
}

func (e *TransactionOperationsIsNil) Error() string {
	return e.Message
}

func NewTransactionOperationsIsNil(
	message string,
) *TransactionOperationsIsNil {
	return &TransactionOperationsIsNil{message}
}

type CannotParseClientPublicKey struct {
	Message string `json:"message"`
}

func (e *CannotParseClientPublicKey) Error() string {
	return e.Message
}

func NewCannotParseClientPublicKey(
	message string,
) *CannotParseClientPublicKey {
	return &CannotParseClientPublicKey{message}
}

type TransactionIsNotSignedByAnchor struct {
	Message string `json:"message"`
}

func (e *TransactionIsNotSignedByAnchor) Error() string {
	return e.Message
}

func NewTransactionIsNotSignedByAnchor(
	message string,
) *TransactionIsNotSignedByAnchor {
	return &TransactionIsNotSignedByAnchor{message}
}

type TransactionIsNotSignedByClient struct {
	Message string `json:"message"`
}

func (e *TransactionIsNotSignedByClient) Error() string {
	return e.Message
}

func NewTransactionIsNotSignedByClient(
	message string,
) *TransactionIsNotSignedByClient {
	return &TransactionIsNotSignedByClient{message}
}

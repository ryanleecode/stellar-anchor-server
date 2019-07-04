package authentication

type TransactionSourceAccountDoesntMatchAnchorPublicKey struct {
	message string
}

func (e *TransactionSourceAccountDoesntMatchAnchorPublicKey) Error() string {
	return e.message
}

func NewTransactionSourceAccountDoesntMatchAnchorPublicKey(
	message string,
) *TransactionSourceAccountDoesntMatchAnchorPublicKey {
	return &TransactionSourceAccountDoesntMatchAnchorPublicKey{message}
}

type TransactionIsMissingTimeBounds struct {
	message string
}

func (e *TransactionIsMissingTimeBounds) Error() string {
	return e.message
}

func NewTransactionIsMissingTimeBounds(
	message string,
) *TransactionIsMissingTimeBounds {
	return &TransactionIsMissingTimeBounds{message}
}

type TransactionChallengeExpired struct {
	message string
}

func (e *TransactionChallengeExpired) Error() string {
	return e.message
}

func NewTransactionChallengeExpired(
	message string,
) *TransactionChallengeExpired {
	return &TransactionChallengeExpired{message}
}

type TransactionChallengeDoesNotHaveOnlyOneOperation struct {
	message string
}

func NewTransactionChallengeDoesNotHaveOnlyOneOperation(
	message string,
) *TransactionChallengeDoesNotHaveOnlyOneOperation {
	return &TransactionChallengeDoesNotHaveOnlyOneOperation{message}
}

func (e *TransactionChallengeDoesNotHaveOnlyOneOperation) Error() string {
	return e.message
}

type TransactionChallengeIsNotAManageDataOperation struct {
	message string
}

func (e *TransactionChallengeIsNotAManageDataOperation) Error() string {
	return e.message
}

func NewTransactionChallengeIsNotAManageDataOperation(
	message string,
) *TransactionChallengeIsNotAManageDataOperation {
	return &TransactionChallengeIsNotAManageDataOperation{message}
}

type TransactionOperationSourceAccountIsEmpty struct {
	message string
}

func (e *TransactionOperationSourceAccountIsEmpty) Error() string {
	return e.message
}

func NewTransactionOperationSourceAccountIsEmpty(
	message string,
) *TransactionOperationSourceAccountIsEmpty {
	return &TransactionOperationSourceAccountIsEmpty{message}
}

type TransactionOperationsIsNil struct {
	message string
}

func (e *TransactionOperationsIsNil) Error() string {
	return e.message
}

func NewTransactionOperationsIsNil(
	message string,
) *TransactionOperationsIsNil {
	return &TransactionOperationsIsNil{message}
}

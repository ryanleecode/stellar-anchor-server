package authorization

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

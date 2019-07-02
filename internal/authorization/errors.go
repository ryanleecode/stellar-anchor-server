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

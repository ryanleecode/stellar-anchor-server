package logic

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/txnbuild"
)

type StellarAccount struct {
	address string
	issuer  Issuer
}

func NewStellarAccount(address string, issuer Issuer) *StellarAccount {
	return &StellarAccount{
		address: address,
		issuer:  issuer,
	}
}

func (s StellarAccount) Issue(amount int64, ethTx *EthereumTransaction) (issueTxHash string, err error) {
	var bytes [32]byte
	copy(bytes[:], ethTx.hash.Bytes())
	return s.issuer.IssueWithMemo(s.address, amount, txnbuild.MemoHash(bytes))
}

type DepositAccount struct {
	ethAddress     string
	stellarAddress string
}

func NewDepositAccount(ethAddress string, stellarAddress string) *DepositAccount {
	return &DepositAccount{ethAddress: ethAddress, stellarAddress: stellarAddress}
}

func (a DepositAccount) DepositInstructions() string {
	return a.ethAddress
}

type DepositAccountStorage interface {
	FindByStellar(addr string) (*DepositAccount, error)
	FindByEth(addr string) (*DepositAccount, error)
	New(stellarAcctAddr string) (*DepositAccount, error)
}

type AccountService struct {
	storage DepositAccountStorage
	issuer  Issuer
}

func NewAccountService(storage DepositAccountStorage, issuer Issuer) *AccountService {
	return &AccountService{storage: storage, issuer: issuer}
}

func (s AccountService) GetDepositAccount(stellarAcctAddr string) (*DepositAccount, error) {
	acct, err := s.storage.FindByStellar(stellarAcctAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to lookup deposit account for %s", stellarAcctAddr)
	}
	if acct != nil {
		return acct, nil
	}

	acct, err = s.storage.New(stellarAcctAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create new deposit account for %s", stellarAcctAddr)
	}

	return acct, nil
}

func (s AccountService) FindAccountFrom(ethAddr string) (*StellarAccount, error) {
	acct, err := s.storage.FindByEth(ethAddr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to lookup deposit account for %s", ethAddr)
	}
	if acct == nil {
		return nil, nil
	}

	return NewStellarAccount(acct.stellarAddress, s.issuer), nil
}

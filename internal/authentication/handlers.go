package authentication

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"stellar-fi-anchor/internal/stellar-sdk"
	"strings"
)

type GetAuthResponse struct {
	Transaction string `json:"transaction"`
}

type AuthenticationService interface {
	BuildSignEncodeChallengeTransactionForAccount(id string) (string, error)
	ValidateClientSignedChallengeTransaction(
		txe *xdr.TransactionEnvelope) []error
	Authenticate(txe *xdr.TransactionEnvelope) (string, error)
}

func NewGetAuthHandler(authService AuthenticationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accountID := r.URL.Query().Get("account")
		if accountID == "" {
			errorPayload := map[string]interface{}{
				"error": "account is a required query parameter",
			}
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
			return
		}

		transaction, err := authService.BuildSignEncodeChallengeTransactionForAccount(accountID)
		if err != nil {
			origErr := errors.Cause(err)
			switch origErr.(type) {
			case *stellarsdk.InvalidAccountID:
				errorPayload := map[string]interface{}{
					"error": "account id is invalid",
				}
				w.WriteHeader(http.StatusBadRequest)
				err := json.NewEncoder(w).Encode(&errorPayload)
				if err != nil {
					panic(err)
				}
				return
			default:
				panic(err)
			}
		}

		dataPayload := GetAuthResponse{Transaction: transaction}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&dataPayload)
		if err != nil {
			panic(err)
		}
	}
}

type transactionAuth struct {
	Transaction string `json:"transaction"`
}

type tokenPayload struct {
	Token string `json:"token"`
}

func NewPostAuthHandler(authService AuthenticationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rules := govalidator.MapData{
			"transaction": []string{"required"},
		}
		messages := govalidator.MapData{
			"transaction": []string{"required"},
		}

		body := transactionAuth{}
		opts := govalidator.Options{
			Request:         r,
			Data:            &body,
			Rules:           rules,
			Messages:        messages,
			RequiredDefault: true,
		}
		v := govalidator.New(opts)
		e := v.ValidateJSON()
		if len(e) > 0 {
			errorPayload := map[string]interface{}{
				"error": map[string]interface{}{
					"message": "request validation error",
					"errors":  e,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
			return
		}

		rawr := strings.NewReader(body.Transaction)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		var txe xdr.TransactionEnvelope
		_, err := xdr.Unmarshal(b64r, &txe)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errorPayload := map[string]interface{}{
				"error": "the transaction cannot be decoded or parsed",
			}
			err := json.NewEncoder(w).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
			return
		}

		validationErrs := authService.ValidateClientSignedChallengeTransaction(&txe)
		for _, e := range validationErrs {
			switch e.(type) {
			case *TransactionSourceAccountDoesntMatchAnchorPublicKey,
				*TransactionIsMissingTimeBounds,
				*TransactionChallengeExpired,
				*TransactionChallengeIsNotAManageDataOperation,
				*TransactionChallengeDoesNotHaveOnlyOneOperation,
				*TransactionOperationSourceAccountIsEmpty,
				*TransactionOperationsIsNil,
				*TransactionIsNotSignedByAnchor,
				*TransactionIsNotSignedByClient:

				continue
			default:
				panic(err)
			}
		}

		if len(validationErrs) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			errorPayload := map[string]interface{}{
				"error": map[string]interface{}{
					"message": "request validation error",
					"errors":  validationErrs,
				},
			}
			err := json.NewEncoder(w).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
			return
		}

		jwtToken, err := authService.Authenticate(&txe)

		dataPayload := tokenPayload{Token: jwtToken}
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(&dataPayload)
		if err != nil {
			panic(err)
		}

	}
}

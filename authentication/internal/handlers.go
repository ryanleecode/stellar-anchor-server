package internal

import (
	"encoding/base64"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
	"github.com/thedevsaddam/govalidator"
	"net/http"
	"strings"
)

type Error struct {
	Message string `json:"message"`
	Errors interface{} `json:"errors"`
}

type ErrorPayload struct {
	Error Error `json:"error"`
}


type GetAuthResponse struct {
	Transaction string `json:"transaction"`
}

type AuthenticationService interface {
	BuildSignEncodeChallengeTransactionForAccount(id string) (string, error)
	ValidateClientSignedChallengeTransaction(
		txe *xdr.TransactionEnvelope) ([]error, error)
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
			case *InvalidAccountID:
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

type TransactionAuth struct {
	Transaction string `json:"transaction"`
}

type TokenPayload struct {
	Token string `json:"token"`
}

type PostAuthResponseHandler func(auth TransactionAuth, w http.ResponseWriter, r *http.Request) (*TokenPayload, error)
type PostAuthRequestValidator func(w http.ResponseWriter, r *http.Request) (*TokenPayload, error)


type Encoder interface {
	Encode(interface{}) error
}

func PostAuthResponseWriter(h PostAuthRequestValidator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var encoder Encoder = json.NewEncoder(w)
		token, err := h(w, r)
		if err != nil {
			switch err.(type) {
			case *ValidationError:
				w.WriteHeader(http.StatusBadRequest)
				writeErr := encoder.Encode(&ErrorPayload{
					Error: Error{
						Message: err.Error(),
						Errors:  err.(*ValidationError).errors,
					},
				})
				if writeErr != nil {
					panic(err)
				}
			default:
				panic(err)
			}

			return
		}

		w.WriteHeader(http.StatusOK)
		err = encoder.Encode(&token)
		if err != nil {
			panic(err)
		}
	}
}

func NewPostAuthRequestValidator(h PostAuthResponseHandler) PostAuthRequestValidator {
	return func(w http.ResponseWriter, r *http.Request) (*TokenPayload, error) {
		rules := govalidator.MapData{
			"transaction": []string{"required"},
		}
		messages := govalidator.MapData{
			"transaction": []string{"required"},
		}
		body := TransactionAuth{}
		opts := govalidator.Options{
			Request:         r,
			Data:            &body,
			Rules:           rules,
			Messages:        messages,
			RequiredDefault: true,
		}
		v := govalidator.New(opts)
		if r.Header.Get("content-type") == "application/json" {
			e := v.ValidateJSON()
			if len(e) > 0 {
				return nil, &ValidationError{ errors: e }
			}
		}
		if r.Header.Get("content-type") == "application/x-www-form-urlencoded" {
			e := v.Validate()
			if len(e) > 0 {
				return nil, &ValidationError{ errors: e }
			}
			body.Transaction = r.Form.Get("transaction")
		}

		return h(body, w, r)
	}
}

func NewPostAuthHandler(authService AuthenticationService) PostAuthResponseHandler {
	return func(auth TransactionAuth, w http.ResponseWriter, r *http.Request) (*TokenPayload, error) {
		rawr := strings.NewReader(auth.Transaction)
		b64r := base64.NewDecoder(base64.StdEncoding, rawr)
		var txe xdr.TransactionEnvelope
		_, err := xdr.Unmarshal(b64r, &txe)
		if err != nil {
			return nil, &BadRequestError{ message:"the transaction cannot be decoded or parsed" }
		}

		validationErrs, err := authService.ValidateClientSignedChallengeTransaction(&txe)
		if err != nil {
			panic(err)
		}

		if len(validationErrs) > 0 {
			return nil, &ValidationError{ errors: validationErrs }
		}

		jwtToken, err := authService.Authenticate(&txe)

		return &TokenPayload{Token: jwtToken}, nil

	}
}

package internal

import (
	"encoding/json"
	"net/http"

	"github.com/drdgvhbh/stellar-anchor-server/ethereum/internal/logic"

	"github.com/gorilla/schema"
	_ "github.com/gorilla/schema"
	"github.com/thedevsaddam/govalidator"
)

type GetDepositQueryParams struct {
	AssetCode string `schema:"asset_code"`
	Account   string `schema:"account"`
}

type Account = interface {
	DepositInstructions() string
}

type AccountService interface {
	GetDepositAccount(stellarAcctAddr string) (*logic.DepositAccount, error)
}

func NewGetDepositHandler(accountService AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rules := govalidator.MapData{
			"asset_code":    []string{"required"},
			"account":       []string{"required"},
			"memo_type":     []string{},
			"memo":          []string{},
			"email_address": []string{},
			"type":          []string{},
		}
		messages := govalidator.MapData{
			"asset_code": []string{"asset_code is required"},
			"account":    []string{"account is required"},
		}
		opts := govalidator.Options{
			Request:  r,
			Rules:    rules,
			Messages: messages,
		}
		validationErrs := govalidator.New(opts).Validate()

		if len(validationErrs) > 0 {
			errorPayload := map[string]interface{}{
				"error": map[string]interface{}{
					"message": "request validation error",
					"errors":  validationErrs,
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(&errorPayload)
			if err != nil {
				panic(err)
			}
			return
		}
		queryParams := GetDepositQueryParams{}
		decoder := schema.NewDecoder()
		err := decoder.Decode(&queryParams, r.URL.Query())
		if err != nil {
			panic(err)
		}

		account, err := accountService.GetDepositAccount(queryParams.Account)
		if err != nil {
			panic(err)
		}

		payload := map[string]interface{}{
			"how": account.DepositInstructions(),
		}
		err = json.NewEncoder(w).Encode(&payload)
		if err != nil {
			panic(err)
		}

	}
}

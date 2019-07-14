package deposits

import (
	"encoding/json"
	"github.com/thedevsaddam/govalidator"
	"net/http"
)

type GetDepositQueryParams struct {
	AssetCode string `json:"asset_code,omitempty"`
	Account   string `json:"account,omitempty"`
}

func NewGetDepositHandler() http.HandlerFunc {
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

		queryParams := GetDepositQueryParams{}
		opts := govalidator.Options{
			Request:  r,
			Rules:    rules,
			Messages: messages,
			Data:     &queryParams,
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
	}
}

package model

import (
	"errors"
	"time"
)

type OptionType string

const (
	Put  OptionType = "Put"
	Call OptionType = "Call"
)

type Position string

const (
	Long  Position = "long"
	Short Position = "short"
)

type OptionsContract struct {
	Type           OptionType `json:"type"`
	LongShort      Position   `json:"long_short"`
	StrikePrice    float64    `json:"strike_price"`
	Bid            float64    `json:"bid"`
	Ask            float64    `json:"ask"`
	ExpirationDate time.Time  `json:"expiration_date"`
}

func IsOptionsContractValid(contract OptionsContract) error {
	// Check for the type being correctly set
	if contract.Type != Call && contract.Type != Put {
		return errors.New("invalid option type. Call or Put")
	}
	// Check that the contract position is correct
	if contract.LongShort != Long && contract.LongShort != Short {
		return errors.New("invalid position type. long or short")
	}
	// The strike has to be greater than 0
	if contract.StrikePrice <= 0 {
		return errors.New("strike price must be greater than zero")
	}
	// The bid cant be negative
	if contract.Bid < 0 {
		return errors.New("bid must be non-negative")
	}
	// The ask cant be negative
	if contract.Ask < 0 {
		return errors.New("ask must be non-negative")
	}
	// The contract cant be expired
	if contract.ExpirationDate.Before(time.Now()) {
		return errors.New("expiration date must be in the future")
	}
	return nil
}

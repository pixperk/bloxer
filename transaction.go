package main

type Transaction struct {
	FromAddress string
	ToAddress   string
	Amount      float64
}

func NewTransaction(from, to string, amount float64) Transaction {
	return Transaction{
		FromAddress: from,
		ToAddress:   to,
		Amount:      amount,
	}
}

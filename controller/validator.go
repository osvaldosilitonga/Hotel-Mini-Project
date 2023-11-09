package controller

func BalanceCheck(balance, amount int) bool {
	if balance < amount {
		return false
	}
	return true
}

package transaction

// Payment represents a transfer of coin from one user to another
type Payment struct {
	From   string
	To     string
	Amount float64
}

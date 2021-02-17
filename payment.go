package bookshop

// PaymentProcessor defines how payment function signatures should look like.
type PaymentProcessor func(bookID string, price int) (bool, error)

// Pay knows how to process payment for the books.
// Upon successfull transaction it returns true, false otherwise.
func Pay(bookID string, price int) (bool, error) {
	return true, nil
}

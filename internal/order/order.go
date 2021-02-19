package order

import (
	"errors"
	"fmt"
)

type order struct {
	OrderID string
	Books   []string
}

// New knows how to construct a valid order.
func New(orderID string) (*order, error) {
	if orderID == "" {
		return nil, errors.New("invalid order id")
	}

	o := order{
		OrderID: orderID,
	}

	return &o, nil
}

// Id returns the order ID.
func (o *order) ID() string {
	return o.OrderID
}

// AddBook knows how to add a book identifued by id to the order.
func (o *order) AddBook(ids ...string) error {
	var correctIDS []string
	var incorrectIDS []string

	for _, b := range ids {
		if b == "" {
			incorrectIDS = append(incorrectIDS, b)
		}
		correctIDS = append(correctIDS, b)
	}
	o.Books = append(o.Books, correctIDS...)
	if len(incorrectIDS) > 0 {
		return fmt.Errorf("Incorrect book ids in order: %v", incorrectIDS)
	}
	return nil
}

// BookIDs returns current list of books added to the order.
func (o *order) BookIDs() []string {
	return o.Books
}

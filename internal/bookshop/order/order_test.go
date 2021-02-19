package order_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/bookshop/internal/bookshop/order"
)

func TestNew(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		orderID     string
		want        string
		expectedErr bool
	}{
		{name: "Correct order id", orderID: "12282", want: "12282", expectedErr: false},
		{name: "Empty order id", orderID: "", want: "", expectedErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			o, err := order.New(tc.orderID)

			if (err != nil) != tc.expectedErr {
				t.Fatalf("%s, order.New(%s) got error %v", tc.name, tc.orderID, err)
			}

			if !tc.expectedErr && (o.OrderID != tc.want) {
				t.Errorf("%s, order.New(%s) = %s, want: %s", tc.name, tc.orderID, o.OrderID, tc.want)
			}

		})
	}
}

func TestNewOrder(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		orderID     string
		want        string
		expectedErr bool
	}{
		{name: "Correct order id", orderID: "12282", want: "12282", expectedErr: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			o, err := order.New(tc.orderID)
			if err != nil {
				t.Fatal(err)
			}

			got := o.ID()

			if !tc.expectedErr && (got != tc.want) {
				t.Errorf("%s, order.ID(%s) = %s, want: %s", tc.name, tc.orderID, o.OrderID, tc.want)
			}

		})
	}
}

func TestAddBook(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		orderID     string
		books       []string
		want        []string
		expectedErr bool
	}{
		{name: "Correct book", orderID: "123", books: []string{"123"}, want: []string{"123"}, expectedErr: false},
		{name: "Correct multiple books", orderID: "234", books: []string{"123", "456", "789"}, want: []string{"123", "456", "789"}, expectedErr: false},
		{name: "Incorrect book ID", orderID: "345", books: []string{""}, want: []string{}, expectedErr: true},
		{name: "Incorrect book ID among correct", orderID: "456", books: []string{"123", "456", "", "789"}, want: []string{"123", "456", "789"}, expectedErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			o, err := order.New(tc.orderID)
			if err != nil {
				t.Fatal(err)
			}

			err = o.AddBook(tc.books...)

			if (err != nil) != tc.expectedErr {
				t.Fatalf("%s, AddBook(%v) got error: %v", tc.name, tc.books, err)
			}

			if !tc.expectedErr && (!cmp.Equal(o.Books, tc.want)) {
				t.Errorf("%s, order.AddBook(%v) got:\n%s", tc.name, tc.books, cmp.Diff(o.Books, tc.want))
			}

		})
	}
}

func TestOrderBookIDs(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		orderID     string
		books       []string
		want        []string
		expectedErr bool
	}{
		{name: "Single book order", orderID: "123", books: []string{"123"}, want: []string{"123"}, expectedErr: false},
		{name: "Multiple books order", orderID: "234", books: []string{"123", "456", "789"}, want: []string{"123", "456", "789"}, expectedErr: false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var err error

			o, err := order.New(tc.orderID)
			if err != nil {
				t.Fatal(err)
			}

			err = o.AddBook(tc.books...)
			if err != nil {
				t.Fatal(err)
			}

			got := o.BookIDs()

			if (err != nil) != tc.expectedErr {
				t.Fatalf("%s, order.BookIDs(%v) got error: %v", tc.name, tc.books, err)
			}

			if !tc.expectedErr && (!cmp.Equal(got, tc.want)) {
				t.Errorf("%s, order.BookIDs(%v) got:\n%s", tc.name, tc.books, cmp.Diff(got, tc.want))
			}
		})
	}
}

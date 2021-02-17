package bookshop_test

import (
	"testing"

	"github.com/qba73/bookshop"
)

func TestPay(t *testing.T) {
	tt := []struct {
		name        string
		bookID      string
		price       int
		want        bool
		expectedErr bool
	}{
		{"Single book valid transaction", "1912bbf7-3f26-4196-b062-071b81b855e9", 1000, true, false},
	}

	for _, tc := range tt {
		got, err := bookshop.Pay(tc.bookID, tc.price)

		if (err != nil) != tc.expectedErr {
			t.Fatalf("expected error")
		}

		if got != tc.want {
			t.Errorf("%s Pay() = %v, want %v", tc.name, got, tc.want)
		}

	}

}

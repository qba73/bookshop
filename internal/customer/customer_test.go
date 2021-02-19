package bookshop_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/qba73/bookshop"
)

var customers = map[string]bookshop.Customer{
	"Customer1": {
		Title:   "Mrs",
		Name:    "Monika White",
		Address: "23 Avenue, Dublin, Ireland",
	},
}

func TestCustomerMailingLabel(t *testing.T) {

	customer := bookshop.Customer{
		Title:   "Mr",
		Name:    "James Brown",
		Address: "43 Temple Gardens, Dublin9, Ireland",
	}

	wantLabel := `Mr James Brown
43 Temple Gardens, Dublin9, Ireland`

	tt := []struct {
		name string
		c    bookshop.Customer
		want string
	}{
		{name: "Full mailing label", c: customer, want: wantLabel},
	}

	for _, tc := range tt {

		got := customer.MailingLabel()

		if !cmp.Equal(got, tc.want) {
			t.Errorf(cmp.Diff(got, tc.want))
		}
	}

}

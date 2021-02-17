package bookshop

import "fmt"

// Customer represent a bookshop customer.
type Customer struct {
	Title   string
	Name    string
	Address string
}

// MailingLabel knows how to construct and present
// customer data required to print mailing labels.
func (c Customer) MailingLabel() string {
	label := fmt.Sprintf("%s %s\n%s", c.Title, c.Name, c.Address)
	return label
}

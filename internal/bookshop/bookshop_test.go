package bookshop_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/qba73/bookshop"
	"github.com/qba73/bookshop/payment"
)

var testBooks = map[string]bookshop.Book{
	"Book1": {
		ID:             "1912bbf7-3f26-4196-b062-071b81b855e9",
		Edition:        1,
		Title:          "Bolek i Lolek",
		Authors:        []string{"Bolek"},
		Description:    "description",
		ReleaseYear:    1997,
		SeriesNumber:   1,
		PriceCents:     2000,
		PickOfTheMonth: true,
	},
	"Book2": {
		ID:             "",
		Edition:        1,
		Title:          "Bolek i Lolek i Matolek",
		Authors:        []string{"Bolek", "Gizmo"},
		Description:    "description",
		ReleaseYear:    1997,
		SeriesNumber:   1,
		PriceCents:     2000,
		PickOfTheMonth: true,
	},
	"Book3": {
		ID:             "1912bbf7-3f26-4196-b062-071b81b855e9",
		Edition:        1,
		Title:          "Bolek i Lolek",
		Authors:        []string{"Bolek", "Papcio"},
		Description:    "description",
		ReleaseYear:    1997,
		SeriesNumber:   1,
		PriceCents:     -2,
		PickOfTheMonth: true,
	},
}

func TestBook(t *testing.T) {
	_ = bookshop.Book{
		ID:             "123",
		Title:          "Bolek i Lolek",
		Authors:        []string{"Gizmo", "Bolek"},
		Description:    "Bolek i Lolek adventures",
		SeriesNumber:   3,
		PriceCents:     30000,
		PickOfTheMonth: true,
	}
}

func TestBook_SalePrice(t *testing.T) {
	bk := bookshop.Book{
		PriceCents: 2000,
	}

	tt := []struct {
		name          string
		b             bookshop.Book
		want          int
		expectedError bool
	}{
		{name: "Calculate NetPrice", b: bk, want: 1600, expectedError: false},
	}

	for _, tc := range tt {
		tc.b.SetDiscountPercent(20)
		got := tc.b.SalePrice()

		if got != tc.want {
			t.Errorf("%s, book.NetPrice() = %d, want: %d", tc.name, got, tc.want)
		}

	}
}

func TestBook_SetBookPrice(t *testing.T) {
	tt := []struct {
		name        string
		b           bookshop.Book
		newPrice    int
		want        int
		expectedErr bool
	}{
		{name: "Change price", b: bookshop.Book{Title: "Fox", Authors: []string{"Gizmo"}, PriceCents: 5000}, newPrice: 2000, want: 2000, expectedErr: false},
		{name: "Change price", b: bookshop.Book{Title: "Fox", Authors: []string{"Gizmo"}, PriceCents: 5000}, newPrice: -1000, want: 0, expectedErr: true},
	}

	for _, tc := range tt {
		got, err := tc.b.SetPriceCents(tc.newPrice)
		if (err != nil) != tc.expectedErr {
			t.Fatalf("%s, SetPiceCents(%d) = %d, want: %d", tc.name, tc.newPrice, got, tc.want)
		}

		if !tc.expectedErr && (got != tc.want) {
			t.Errorf("%s, SetPiceCents(%d) = %d, want: %d", tc.name, tc.newPrice, got, tc.want)
		}

		// Verify if the book struct has been modified.
		if !tc.expectedErr && (tc.b.PriceCents != tc.newPrice) {
			t.Errorf("got: %d, want: %d", tc.b.PriceCents, tc.newPrice)
		}
	}
}

func TestBook_SetCategory(t *testing.T) {
	b := bookshop.Book{
		Title:   "Hello World",
		Authors: []string{"Bolek"},
	}

	tt := []struct {
		name        string
		book        bookshop.Book
		category    int
		expectedErr bool
	}{
		{name: "Set invalid category", book: b, category: 10, expectedErr: true},
		{name: "Set category", book: b, category: bookshop.CategoryTech, expectedErr: false},
		{name: "Set invalid category", book: b, category: bookshop.CategoryRomance, expectedErr: false},
		{name: "Set invalid category", book: b, category: 10, expectedErr: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := b.SetCategory(tc.category)

			if (err != nil) != tc.expectedErr {
				t.Errorf("%s SetCategory(%d) got: %v", tc.name, tc.category, err)
			}
		})
	}
}

func TestBook_SetDiscount(t *testing.T) {
	b := bookshop.Book{
		Title:      "Harry",
		Authors:    []string{"Gizmo"},
		PriceCents: 2000,
	}

	tt := []struct {
		name        string
		book        bookshop.Book
		discount    int
		expectedErr bool
	}{
		{name: "Coorect discount range", book: b, discount: 50, expectedErr: false},
		{name: "Discount eq 0", book: b, discount: 0, expectedErr: false},
		{name: "Discount eq 100", book: b, discount: 100, expectedErr: false},
		{name: "Discount lt 0 incorrect", book: b, discount: -10, expectedErr: true},
		{name: "Discount gt 100 incorrect", book: b, discount: 110, expectedErr: true},
	}

	for _, tc := range tt {

		err := b.SetDiscountPercent(tc.discount)

		if !tc.expectedErr && (err != nil) {
			t.Errorf("%s, SetDiscount(%q) got %v, want nil", tc.name, tc.discount, err)
		}
	}
}

func TestGetAllBooks(t *testing.T) {
	want := bookshop.Books
	got := bookshop.GetAllBooks()

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(bookshop.Book{})) {
		t.Errorf("GetAllBooks() diff:%s\n", cmp.Diff(want, got))
	}
}

func TestNewID(t *testing.T) {
	got := bookshop.NewID()
	u, err := uuid.Parse(got)

	if err != nil {
		t.Errorf("NewID() error : incorrect id")
	}

	want := u.String()

	if got != want {
		t.Errorf("NewID() = %s; want: %s", got, want)
	}
}

func TestGetBookDetails(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name        string
		bookID      string
		want        string
		expectedErr bool
	}{
		{"Single author", "1912bbf7-3f26-4196-b062-071b81b855e9", "Title: Bolek i Lolek, Author: Bolek, Year: 1997, ID: 1912bbf7-3f26-4196-b062-071b81b855e9", false},
		{"Multiple authors", "1923bbf9-3f36-4196-b062-171b81b855e9", "Title: Zosia Samosia, Authors: Papcio Chmiel, Zigmas Laurin, Year: 2011, ID: 1923bbf9-3f36-4196-b062-171b81b855e9", false},

		// Expected errors
		//{"Single author incorrect description", "1912bbf7-3f26-4196-b062-071b81b855e9", "Title: Bolek i Lolek, Authors: Bolek, Year: 1997, ID: 1912bbf7-3f26-4196-b062-071b81b855e9", true},
		//{"Not existing book", "9992bbf7-3f26-4196-b062-071b81b855e9", "Title: Bolek i Lolek, Author: Bolek, Year: 1997, ID: 1912bbf7-3f26-4196-b062-071b81b855e9", true},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bookshop.GetBookDetails(tc.bookID, bookshop.Books)

			if (err != nil) != tc.expectedErr {
				t.Fatalf("%s GetBookDetails() = expected error", tc.name)
			}

			if got != tc.want {
				t.Errorf("%s GetBookDetails() = \n%s\nwant\n%s", tc.name, got, tc.want)
			}

			/*
				if !cmp.Equal(got, tc.want) && !tc.expectErr {
					t.Errorf("%s\n%s", tc.name, cmp.Diff(tc.want, got))
				}
			*/
		})
	}
}

func TestGetAllByAuthor(t *testing.T) {

	tt := []struct {
		name        string
		author      string
		books       map[string]bookshop.Book
		want        []string
		expectedErr bool
	}{
		{"Single author", "Bolek", bookshop.Books, []string{"1912bbf7-3f26-4196-b062-071b81b855e9", "2922bbf7-3g26-4196-b062-071b81b855e9"}, false},
		{"Single author", "Gienek", bookshop.Books, []string{"1912abf7-3f26-4196-b062-011b81b255e9"}, false},
		{"Multiple authors", "Gizmo", bookshop.Books, []string{"1923bbf9-4f36-4196-b062-171b81b855e9"}, false},
	}

	for _, tc := range tt {
		got, err := bookshop.GetAllByAuthor(tc.author, bookshop.Books)

		if (err != nil) != tc.expectedErr {
			t.Fatalf("%s GetAllByAuthor(%s) should return error", tc.name, tc.author)
		}

		if !cmp.Equal(got, tc.want) {
			t.Errorf("%s GetAllByAuthor(%s) \n%s", tc.name, tc.author, cmp.Diff(tc.want, got))
		}
	}
}

func TestNetPrice(t *testing.T) {
	tt := []struct {
		name string
		book bookshop.Book
		want int
	}{
		{"Calculate net price", bookshop.Book{ID: "1912bbf7-3f26-4196-b062-071b81b855e9", Title: "Bolek i Lolek", PriceCents: 2000}, 1600},
	}

	for _, tc := range tt {

		got, _ := bookshop.NetPrice(tc.book, bookshop.Books)

		if got != tc.want {
			t.Errorf("%s; NetPrice() = %d; want %d", tc.name, got, tc.want)
		}
	}
}

func TestGetAllBookDetails(t *testing.T) {
	allbooks := `Title: Bolek i Lolek, Author: Bolek, Year: 1997, ID: 1912bbf7-3f26-4196-b062-071b81b855e9
Title: Koziolek Matolek, Author: Bolek, Year: 1967, ID: 2922bbf7-3g26-4196-b062-071b81b855e9
Title: Pan Samochodzik, Authors: Papcio Chmiel, Zigmas Laurin, Gizmo, Year: 2011, ID: 1923bbf9-4f36-4196-b062-171b81b855e9
Title: Tytus, Author: Gienek, Year: 2017, ID: 1912abf7-3f26-4196-b062-011b81b255e9
Title: Zosia Samosia, Authors: Papcio Chmiel, Zigmas Laurin, Year: 2011, ID: 1923bbf9-3f36-4196-b062-171b81b855e9
`

	tt := []struct {
		name  string
		books map[string]bookshop.Book
		want  string
	}{
		{"All books in catalog", bookshop.Books, allbooks},
	}

	for _, tc := range tt {
		got := bookshop.GetAllBookDetails(tc.books)

		if !cmp.Equal(got, tc.want) {
			t.Errorf("%s GetAllBookDetails() =\n%s", tc.name, cmp.Diff(tc.want, got))
		}
	}
}

func TestBuyBook(t *testing.T) {
	paymentOK := func(bookID string, price int) (bool, error) {
		return true, nil
	}

	paymentFailed := func(bookID string, price int) (bool, error) {
		return false, errors.New("Payment failed")
	}

	tt := []struct {
		name        string
		book        bookshop.Book
		payer       payment.Processor
		want        bool
		expectedErr bool
	}{
		{name: "Succesful payment", book: testBooks["Book1"], payer: paymentOK, want: true, expectedErr: false},
		{name: "Not successful payment", book: testBooks["Book1"], payer: paymentFailed, want: false, expectedErr: true},
		{name: "Missing book id", book: testBooks["Book2"], payer: paymentOK, want: false, expectedErr: true},
		{name: "Incorrect price", book: testBooks["Books3"], payer: paymentOK, want: false, expectedErr: true},
	}

	for _, tc := range tt {
		got, err := bookshop.BuyBook(tc.book.ID, tc.book.PriceCents, tc.payer)

		if (err != nil) != tc.expectedErr {
			t.Fatalf("%s, BuyBook() = %v, expected unsuccesfull payment: %v, err: %s", tc.name, got, tc.want, err)
		}

		if !tc.expectedErr && got != tc.want {
			t.Errorf("%s, BuyBook(%s, %d) = %v, want: %v", tc.name, tc.book.ID, tc.book.PriceCents, got, tc.want)
		}
	}
}

func TestCatalogGetAllBooks(t *testing.T) {
	want := []bookshop.Book{
		testBooks["Book1"],
		testBooks["Book2"],
	}

	c := bookshop.Catalog{
		Books: []bookshop.Book{
			testBooks["Book1"],
			testBooks["Book2"],
		},
	}

	got := c.GetAllBooks()

	if !cmp.Equal(got, want, cmpopts.IgnoreUnexported(bookshop.Book{})) {
		t.Errorf(cmp.Diff(got, want))
	}
}

func TestCatalogLen(t *testing.T) {
	books := []bookshop.Book{
		testBooks["Book1"],
		testBooks["Book2"],
	}

	c := bookshop.Catalog{
		Books: books,
	}
	want := 2
	got := c.Len()

	if got != want {
		t.Errorf("Len() = %d, want: %d", got, want)
	}
}

func TestCatalogGetAllTitles(t *testing.T) {
	books := []bookshop.Book{
		testBooks["Book1"],
		testBooks["Book2"],
	}

	want := []string{"Bolek i Lolek", "Bolek i Lolek i Matolek"}

	c := bookshop.Catalog{
		Books: books,
	}

	got := c.GetAllTitles()

	if !cmp.Equal(got, want) {
		t.Errorf(cmp.Diff(got, want))
	}
}

func TestCatalogAddBook(t *testing.T) {
	b1 := bookshop.Book{
		ID:             "1212bbg6-3f26-4196-b062-071b81b855e9",
		Edition:        1,
		Title:          "Jas i Malgosia",
		Authors:        []string{"Zigmas"},
		Description:    "description",
		ReleaseYear:    1993,
		SeriesNumber:   1,
		PriceCents:     2000,
		PickOfTheMonth: true,
	}

	var ct bookshop.Catalog

	ct.AddBook(b1)

	books := ct.GetAllBooks()

	if len(books) != 1 {
		t.Fatalf("Book not added to the catalog")
	}

	want := books[0]

	if !cmp.Equal(b1, want, cmpopts.IgnoreUnexported(bookshop.Book{})) {
		t.Errorf(cmp.Diff(b1, want))
	}
}

func TestGetUniqueAuthors(t *testing.T) {
	var books []bookshop.Book

	for _, v := range testBooks {
		books = append(books, v)
	}

	c := bookshop.Catalog{
		Books: books,
	}
	want := []string{"Bolek", "Gizmo", "Papcio"}
	got := c.GetUniqueAuthors()

	if !cmp.Equal(want, got, cmpopts.IgnoreUnexported(bookshop.Book{})) {
		t.Errorf(cmp.Diff(want, got))
	}
}

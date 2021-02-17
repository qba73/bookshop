package bookshop

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

const (
	CategoryAutobiography = iota
	CategoryTech
	CategoryRomance
	CategoryProgramming
)

// Book represent a single book in the bookshop.
type Book struct {
	ID             string
	Edition        int
	Title          string
	Authors        []string
	Description    string
	ReleaseYear    int
	SeriesNumber   int
	PriceCents     int
	PickOfTheMonth bool
	discount       int
	category       int
}

// String implements Stringer interface for the Book struct.
func (b *Book) String() string {
	switch len(b.Authors) > 1 {
	case true:
		authors := strings.Join(b.Authors, ", ")
		return fmt.Sprintf("Title: %v, Authors: %v, Year: %v, ID: %v", b.Title, authors, b.ReleaseYear, b.ID)
	default:
		authors := strings.Join(b.Authors, "")
		return fmt.Sprintf("Title: %v, Author: %v, Year: %v, ID: %v", b.Title, authors, b.ReleaseYear, b.ID)
	}
}

// SalePrice knows how to calculate price for the book with applied discount.
func (b *Book) SalePrice() int {
	return b.PriceCents - (b.PriceCents * b.discount / 100)
}

// SetPriceCents ...
func (b *Book) SetPriceCents(p int) (int, error) {
	if p < 0 {
		return 0, fmt.Errorf("Invalid book price: %d", p)
	}
	b.PriceCents = p
	return b.PriceCents, nil
}

// SetCategory ...
func (b *Book) SetCategory(c int) error {
	if !validCategory(c) {
		return fmt.Errorf("unknown category: %q", c)
	}
	b.category = c
	return nil
}

// SetDiscountPercent knows how to discount a book with
// given discount percentage. Valid values  0 < discount < 100.
// It returns error if the discount value is not in the allowed range.
func (b *Book) SetDiscountPercent(d int) error {
	if d < 0 || d > 100 {
		return fmt.Errorf("Invalid discount value: %d", d)
	}
	b.discount = d
	return nil
}

// Books represent our collection of books, a small home library.
var Books = map[string]Book{
	"1912bbf7-3f26-4196-b062-071b81b855e9": {
		ID:             "1912bbf7-3f26-4196-b062-071b81b855e9",
		Edition:        1,
		Title:          "Bolek i Lolek",
		Authors:        []string{"Bolek"},
		Description:    "description",
		ReleaseYear:    1997,
		SeriesNumber:   1,
		PriceCents:     2000,
		PickOfTheMonth: true,
		discount:       20,
		category:       1,
	},
	"1912abf7-3f26-4196-b062-011b81b255e9": {
		ID:             "1912abf7-3f26-4196-b062-011b81b255e9",
		Edition:        1,
		Title:          "Tytus",
		Authors:        []string{"Gienek"},
		Description:    "description",
		ReleaseYear:    2017,
		SeriesNumber:   2,
		PriceCents:     3000,
		PickOfTheMonth: false,
		discount:       10,
		category:       0,
	},
	"2922bbf7-3g26-4196-b062-071b81b855e9": {
		ID:             "2922bbf7-3g26-4196-b062-071b81b855e9",
		Edition:        3,
		Title:          "Koziolek Matolek",
		Authors:        []string{"Bolek"},
		Description:    "description",
		ReleaseYear:    1967,
		SeriesNumber:   2,
		PriceCents:     2500,
		PickOfTheMonth: false,
		discount:       8,
		category:       0,
	},
	"1923bbf9-3f36-4196-b062-171b81b855e9": {
		ID:             "1923bbf9-3f36-4196-b062-171b81b855e9",
		Edition:        1,
		Title:          "Zosia Samosia",
		Authors:        []string{"Papcio Chmiel", "Zigmas Laurin"},
		Description:    "description",
		ReleaseYear:    2011,
		SeriesNumber:   1,
		PriceCents:     1000,
		PickOfTheMonth: true,
		discount:       5,
		category:       2,
	},
	"1923bbf9-4f36-4196-b062-171b81b855e9": {
		ID:             "1923bbf9-4f36-4196-b062-171b81b855e9",
		Edition:        1,
		Title:          "Pan Samochodzik",
		Authors:        []string{"Papcio Chmiel", "Zigmas Laurin", "Gizmo"},
		Description:    "description",
		ReleaseYear:    2011,
		SeriesNumber:   1,
		PriceCents:     1000,
		PickOfTheMonth: true,
		discount:       5,
		category:       1,
	},
}

// Catalog represents book catalog in a bookstore.
type Catalog struct {
	Books []Book
}

// GetAllBooks nows how to return all books in the bookstore's catalog.
func (c *Catalog) GetAllBooks() []Book {
	return c.Books
}

// Len knows how to return a total number of books in the catalog.
func (c *Catalog) Len() int {
	return len(c.Books)
}

// GetAllTitles returns a slice of book titles in the catalog.
func (c *Catalog) GetAllTitles() []string {
	var titles []string
	for _, t := range c.Books {
		titles = append(titles, t.Title)
	}
	return titles
}

// GetUniqueAuthors knows how to scan catalog and
// return a slice of unque authors across the catalog.
func (c *Catalog) GetUniqueAuthors() []string {
	authors := make(map[string]int)

	for _, b := range c.Books {
		for _, v := range b.Authors {
			_, ok := authors[v]
			if !ok {
				authors[v] = 1
			}
		}
	}

	var uniqueAuthors []string

	for k := range authors {
		uniqueAuthors = append(uniqueAuthors, k)
	}
	sort.Strings(uniqueAuthors)

	return uniqueAuthors
}

// AddBook adds a book to the catalog.
func (c *Catalog) AddBook(b Book) {
	c.Books = append(c.Books, b)
}

// GetAllBooks ....com
func GetAllBooks() map[string]Book {
	return Books
}

// NewID generates a unique uuid string.
func NewID() string {
	return uuid.New().String()
}

// GetBookDetails ...
func GetBookDetails(bookID string, books map[string]Book) (string, error) {
	b, ok := books[bookID]
	if !ok {
		return "", fmt.Errorf("book id %s not found", bookID)
	}

	return fmt.Sprint(&b), nil
}

// GetAllByAuthor ...
func GetAllByAuthor(author string, books map[string]Book) ([]string, error) {
	var bookIds []string

	for k, v := range books {
		for _, a := range v.Authors {
			if a == author {
				bookIds = append(bookIds, k)
			}
		}
	}
	sort.Strings(bookIds)
	return bookIds, nil
}

// NetPrice calculates the price with discount applied.
func NetPrice(b Book, books map[string]Book) (int, error) {
	book, ok := books[b.ID]
	if !ok {
		return 0, fmt.Errorf("book %s not found", b.ID)
	}
	price := book.PriceCents - (book.PriceCents * book.discount / 100)
	return price, nil
}

// GetAllBookDetails returns a list of book catalog in one string.
// The list is sorted by title.
func GetAllBookDetails(books map[string]Book) string {
	var bookDetails []string
	var bks []Book

	for _, b := range books {
		bks = append(bks, b)
	}

	sort.Slice(bks, func(i, j int) bool {
		return bks[i].Title < bks[j].Title
	})

	for _, b := range bks {
		bookDetails = append(bookDetails, fmt.Sprintf("%s\n", &b))
	}

	return strings.Join(bookDetails, "")
}

// BuyBook knows how to
func BuyBook(bookID string, price int, processPayment func(bookID string, price int) (bool, error)) (bool, error) {
	if bookID == "" {
		return false, errors.New("Invalid bookID")
	}
	if price < 0 {
		return false, fmt.Errorf("Invalid book price: %d", price)
	}
	return processPayment(bookID, price)
}

func validCategory(c int) bool {
	validCategories := map[int]bool{
		CategoryAutobiography: true,
		CategoryRomance:       true,
		CategoryProgramming:   true,
		CategoryTech:          true,
	}
	return validCategories[c]
}

package main

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

type Book struct {
	Title           string
	AuthorFirstName string
	AuthorLastName  string
	Type            string
	Shelf           int
	Row             int
}

type Library struct {
	Books []Book
}

func (l *Library) LoadBooks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		l.Books = append(l.Books, Book{
			Title:           record[0],
			AuthorFirstName: record[1],
			AuthorLastName:  record[2],
			Type:            record[3],
		})
	}
	return nil
}

type Shelf struct {
	Rows [][]Book
}

type DisplayLibrary struct {
	Shelves      []Shelf
	RowsPerShelf int
	BooksPerRow  int
}

func (dl *DisplayLibrary) ArrangeBooks(books []Book) {
	shelf := Shelf{}
	row := []Book{}

	for _, book := range books {
		row = append(row, book)
		if len(row) == dl.BooksPerRow {
			shelf.Rows = append(shelf.Rows, row)
			row = []Book{}
		}

		if len(shelf.Rows) == dl.RowsPerShelf {
			dl.Shelves = append(dl.Shelves, shelf)
			shelf = Shelf{}
		}
	}

	if len(row) > 0 {
		shelf.Rows = append(shelf.Rows, row)
	}
	if len(shelf.Rows) > 0 {
		dl.Shelves = append(dl.Shelves, shelf)
	}
}
func (l *Library) SortBooks() {
	sort.SliceStable(l.Books, func(i, j int) bool {
		if l.Books[i].Type == l.Books[j].Type {
			if l.Books[i].AuthorFirstName == l.Books[j].AuthorFirstName {
				return l.Books[i].AuthorLastName < l.Books[j].AuthorLastName
			}
			return l.Books[i].AuthorFirstName < l.Books[j].AuthorFirstName
		}
		return l.Books[i].Type < l.Books[j].Type
	})
}
func (dl *DisplayLibrary) DisplayLibrary() {
	for i, shelf := range dl.Shelves {
		fmt.Printf("Shelf %d:\n", i+1)
		for j, row := range shelf.Rows {
			fmt.Printf("  Row %d: ", j+1)
			for _, book := range row {
				fmt.Printf("[%s by %s %s (%s)] ", book.Title, book.AuthorFirstName, book.AuthorLastName, book.Type)
			}
			fmt.Println()
		}
	}
}
func (l *Library) ShuffleBooks() {
	for i := range l.Books {
		j := rand.Intn(len(l.Books))
		l.Books[i], l.Books[j] = l.Books[j], l.Books[i]
	}
}
func (l *Library) SearchBooks(query string) []Book {
	var results []Book
	query = strings.ToLower(query)

	for _, book := range l.Books {
		if strings.Contains(strings.ToLower(book.Title), query) ||
			strings.Contains(strings.ToLower(book.AuthorFirstName), query) ||
			strings.Contains(strings.ToLower(book.AuthorLastName), query) {
			results = append(results, book)
		}
	}
	return results
}
func (l *Library) SearchMultiBooks(titles ...string) map[string]string {
	results := make(map[string]string)

	for _, title := range titles {
		found := false
		for _, book := range l.Books {
			if book.Title == title {
				results[title] = fmt.Sprintf("Found at shelf %d, Row %d", book.Shelf, book.Row)
				found = true
				break
			}
		}
		if !found {
			results[title] = "Not found in the Library"
		}
	}
	return results
}

func main() {
	library := Library{}
	err := library.LoadBooks("books.csv")
	if err != nil {
		fmt.Println("Error loading books:", err)
		return
	}

	library.SortBooks()
	//fmt.Println("Sorted Books:", library.Books)
	//library.ShuffleBooks()
	//fmt.Println("Shuffled Books:", library.Books)
	//searchResults := library.SearchBooks("Anna Karenina")
	//fmt.Println("Search Results:", searchResults)

	searchResults := library.SearchMultiBooks("Maram", "The Brothers Karamazov")
	for title, result := range searchResults {
		fmt.Printf("Book '%s': %s\n", title, result)
	}
	displayLibrary := DisplayLibrary{
		RowsPerShelf: 2,
		BooksPerRow:  3,
	}
	displayLibrary.ArrangeBooks(library.Books)

	displayLibrary.DisplayLibrary()
}

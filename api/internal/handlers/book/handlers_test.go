package book

import (
	"booksonhooks.ca/internal/logger"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"booksonhooks.ca/internal/domain"
	"github.com/julienschmidt/httprouter"
)

// --- Mocks ---

type mockBookRepo struct {
	getBooksFn    func() ([]domain.Book, error)
	getBookByIDFn func(id int64) (*domain.Book, error)
	updateBookFn  func(book *domain.Book) error
	deleteBookFn  func(id int64) error
}

func (m *mockBookRepo) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.getBooksFn()
}

func (m *mockBookRepo) GetBookByID(ctx context.Context, id int64) (*domain.Book, error) {
	return m.getBookByIDFn(id)
}

func (m *mockBookRepo) UpdateBook(ctx context.Context, book *domain.Book) error {
	if m.updateBookFn == nil {
		return nil
	}
	return m.updateBookFn(book)
}

func (m *mockBookRepo) DeleteBook(ctx context.Context, id int64) error {
	if m.deleteBookFn == nil {
		return nil
	}
	return m.deleteBookFn(id)
}

// --- Tests ---

func TestGetBooksHandler(t *testing.T) {
	mockRepo := &mockBookRepo{
		getBooksFn: func() ([]domain.Book, error) {
			return []domain.Book{{Title: "Book 1"}, {Title: "Book 2"}}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, &logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books", nil)
	rr := httptest.NewRecorder()

	h.GetBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetBookHandler(t *testing.T) {
	mockRepo := &mockBookRepo{
		getBooksFn: func() ([]domain.Book, error) { return []domain.Book{}, nil },
		getBookByIDFn: func(id int64) (*domain.Book, error) {
			return &domain.Book{ID: id, Title: "Book 1"}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, &logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books/12", nil)
	params := httprouter.Params{{Key: "id", Value: "12"}}
	req = req.WithContext(context.WithValue(req.Context(), httprouter.ParamsKey, params))
	rr := httptest.NewRecorder()

	h.GetBook(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp domain.Book
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != 12 {
		t.Fatalf("expected id=12, got %d", resp.ID)
	}
}

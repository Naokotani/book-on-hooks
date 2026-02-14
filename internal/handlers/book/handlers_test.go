package book

import (
	"booksonhooks.ca/internal/logger"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"booksonhooks.ca/internal/domain"
	"github.com/julienschmidt/httprouter"
)

// --- Mocks ---

type mockBookRepo struct {
	getBookByRowAndColFn func(row, col int) (*domain.Book, error)
	getBooksFn           func() ([]domain.Book, error)
	getBookByIDFn        func(id int64) (*domain.Book, error)
}

func (m *mockBookRepo) GetBookByRowAndCol(ctx context.Context, row, col int) (*domain.Book, error) {
	return m.getBookByRowAndColFn(row, col)
}

func (m *mockBookRepo) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.getBooksFn()
}

func (m *mockBookRepo) GetBookByID(ctx context.Context, id int64) (*domain.Book, error) {
	return m.getBookByIDFn(id)
}

// --- Tests ---

func TestCoverHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockRepo := &mockBookRepo{
			getBookByRowAndColFn: func(row, col int) (*domain.Book, error) {
				return &domain.Book{Title: "The Great Gatsby"}, nil
			},
		}

		logger := logger.NewLogger("info")

		h := New(nil, nil, &logger, mockRepo)

		req := httptest.NewRequest("GET", "/api/cover/1/1", nil)

		// Inject httprouter params
		params := httprouter.Params{
			{Key: "row", Value: "1"},
			{Key: "col", Value: "1"},
		}
		ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.Cover(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", rr.Code)
		}
	})

	t.Run("Invalid Row Parameter", func(t *testing.T) {
		h := New(nil, nil, nil, nil) // Repo shouldn't even be called

		req := httptest.NewRequest("GET", "/api/cover/abc/1", nil)
		params := httprouter.Params{
			{Key: "row", Value: "abc"},
			{Key: "col", Value: "1"},
		}
		ctx := context.WithValue(req.Context(), httprouter.ParamsKey, params)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		h.Cover(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", rr.Code)
		}
	})
}

func TestGetBooksHandler(t *testing.T) {
	mockRepo := &mockBookRepo{
		getBooksFn: func() ([]domain.Book, error) {
			return []domain.Book{{Title: "Book 1"}, {Title: "Book 2"}}, nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, nil, &logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/books", nil)
	rr := httptest.NewRecorder()

	h.GetBooks(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

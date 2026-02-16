package machine

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"booksonhooks.ca/internal/domain"
	"booksonhooks.ca/internal/logger"
)

type mockMachineRepo struct {
	getMachinesFn    func() ([]domain.Machine, error)
	getMachineByIDFn func(id int64) (*domain.Machine, error)
	insertMachineFn  func(machine *domain.Machine) (int64, error)
	loadMachineFn    func(machineID int64, books []domain.BookMachine) error
}

func (m *mockMachineRepo) GetMachines(ctx context.Context) ([]domain.Machine, error) {
	return m.getMachinesFn()
}

func (m *mockMachineRepo) GetMachineById(ctx context.Context, id int64) (*domain.Machine, error) {
	return m.getMachineByIDFn(id)
}

func (m *mockMachineRepo) InsertMachine(ctx context.Context, machine *domain.Machine) (int64, error) {
	return m.insertMachineFn(machine)
}

func (m *mockMachineRepo) LoadMachine(ctx context.Context, machineID int64, books []domain.BookMachine) error {
	return m.loadMachineFn(machineID, books)
}

func TestMachinesViewHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{{Location: "A1", Rows: 6, Columns: 5}}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id, Location: "A1", Rows: 6, Columns: 5}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			return 1, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, nil, &logger, mockRepo)

	req := httptest.NewRequest("GET", "/api/machines", nil)
	rr := httptest.NewRecorder()

	h.MachinesView(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestMachineCreateHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			if machine.Location != "HQ" || machine.Rows != 3 || machine.Columns != 4 {
				t.Fatalf("unexpected machine payload: %+v", machine)
			}
			return 42, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, nil, &logger, mockRepo)

	body := []byte(`{"location":"HQ","rows":3,"cols":4}`)
	req := httptest.NewRequest("POST", "/api/admin/machine/create", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.MachineCreate(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}
}

func TestLoadMachineHandler(t *testing.T) {
	mockRepo := &mockMachineRepo{
		getMachinesFn: func() ([]domain.Machine, error) {
			return []domain.Machine{}, nil
		},
		getMachineByIDFn: func(id int64) (*domain.Machine, error) {
			return &domain.Machine{ID: id}, nil
		},
		insertMachineFn: func(machine *domain.Machine) (int64, error) {
			return 1, nil
		},
		loadMachineFn: func(machineID int64, books []domain.BookMachine) error {
			if machineID != 7 {
				t.Fatalf("expected machineID=7, got %d", machineID)
			}
			if len(books) != 2 {
				t.Fatalf("expected 2 books, got %d", len(books))
			}
			if books[0].BookID != 11 || books[0].Row != 1 || books[0].Col != 2 {
				t.Fatalf("unexpected first book payload: %+v", books[0])
			}
			if books[1].BookID != 12 || books[1].Row != 3 || books[1].Col != 4 {
				t.Fatalf("unexpected second book payload: %+v", books[1])
			}
			return nil
		},
	}

	logger := logger.NewLogger("info")
	h := New(nil, nil, &logger, mockRepo)

	body := []byte(`{"machine_id":7,"books":[{"book_id":11,"row":1,"col":2},{"book_id":12,"row":3,"col":4}]}`)
	req := httptest.NewRequest("POST", "/api/admin/machine/load", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	h.LoadMachine(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp struct {
		MachineID int64 `json:"machine_id"`
		Count     int   `json:"count"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response json: %v", err)
	}

	if resp.MachineID != 7 || resp.Count != 2 {
		t.Fatalf("unexpected response payload: %+v", resp)
	}
}

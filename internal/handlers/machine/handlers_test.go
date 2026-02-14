package machine

import (
	"bytes"
	"context"
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

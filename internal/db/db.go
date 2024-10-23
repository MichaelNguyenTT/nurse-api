package db

import (
	"errors"
	"math/rand"
	"nms/internal/models"
	"sync"
	"time"
)

var (
	ErrNoData     = errors.New("data base is empty")
	ErrNotFound   = errors.New("no subject with ID found")
	ErrExistingID = errors.New("found existing ID in DB")
)

type SubjectDB struct {
	ID        int
	Name      string
	Category  string
	Priority  int
	Notes     string
	CreatedAt string
}

type SubjectCollection struct {
	storage map[int]SubjectDB
	mu      sync.RWMutex
}

func NewPDB() *SubjectCollection {
	return &SubjectCollection{
		storage: make(map[int]SubjectDB),
	}
}

func (p *SubjectCollection) All() ([]SubjectDB, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	storage := p.storage
	if len(storage) == 0 {
		return []SubjectDB{}, ErrNoData
	}

	var patients []SubjectDB

	for _, p := range storage {
		patients = append(patients, p)
	}

	return patients, nil
}

func (p *SubjectCollection) ByID(pID int) (SubjectDB, error) {
	found, ok := p.storage[pID]
	if !ok {
		return SubjectDB{}, ErrNotFound
	}

	return found, nil
}

func (p *SubjectCollection) Add(patient models.SubjectModel) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	formattedPatient := processSubject(patient)

	//FIX: hack bc we havent implemented having unique IDs...
	_, exists := p.storage[formattedPatient.ID]
	if exists {
		return ErrExistingID
	}

	p.storage[formattedPatient.ID] = *formattedPatient

	return nil
}

func (p *SubjectCollection) Update(pr models.SubjectModel) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	updatePatient := processSubject(pr)
	p.storage[updatePatient.ID] = *updatePatient

	return nil
}

func (p *SubjectCollection) Delete(id int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, exists := p.storage[id]
	if !exists {
		return ErrNotFound
	}

	delete(p.storage, id)
	return nil
}

// processor to db from request
func processSubject(p models.SubjectModel) *SubjectDB {
	return &SubjectDB{
		ID:        generateID(),
		Name:      p.Name,
		Category:  p.Category,
		Priority:  p.Priority,
		CreatedAt: generateDateCreated(),
	}
}

func generateDateCreated() string {
	now := time.Now()
	return now.Format("2006-01-02")
}

func generateID() int {
	id := rand.Intn(1000)
	return id
}

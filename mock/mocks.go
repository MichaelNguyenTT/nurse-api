package mock

import (
	"nms/internal/db"
	"nms/internal/models"
)

func LoadMockData() *db.SubjectCollection {
	m1 := models.SubjectModel{
		Name:     "Pharmacology",
		Category: "Core",
		Priority: 3,
		Notes:    "nurse stuff idk",
	}

	m2 := models.SubjectModel{
		Name:     "Brain",
		Category: "Neuroscience",
		Priority: 1,
		Notes:    "the brain is complex",
	}

	m3 := models.SubjectModel{
		Name:     "dummy",
		Category: "something",
		Priority: 1,
		Notes:    "some test notes",
	}

	db := db.NewPDB()
	db.Add(m1)
	db.Add(m2)
	db.Add(m3)

	return db
}

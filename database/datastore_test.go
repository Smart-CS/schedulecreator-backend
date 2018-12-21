package database_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/database"
	"github.com/smart-cs/scheduler-backend/models"
	"github.com/stretchr/testify/assert"
)

func setup() {
	database.LoadLocalDatabase("test-coursedb.json")
}

func TestGetSections(t *testing.T) {
	setup()
	assert := assert.New(t)
	ds := database.NewDatastore()

	assert.Len(ds.GetSections("CPSC 110", "1-2", models.Lecture), 8)
	assert.Len(ds.GetSections("CPSC 110", "1-2", models.Laboratory), 53)
	// TODO: This should be 0 after updating the database.
	assert.Len(ds.GetSections("CPSC 110", "1-2", models.Tutorial), 1)
	assert.Len(ds.GetSections("CPSC 110", "1-2", models.Lecture, models.Laboratory, models.Tutorial), 8+53+1)

	assert.Len(ds.GetSections("CPEN 221", "1", models.Lecture), 1)
	assert.Len(ds.GetSections("CPEN 221", "1", models.Laboratory), 5)
	assert.Len(ds.GetSections("CPEN 221", "1", models.Tutorial), 1)
	assert.Len(ds.GetSections("CPEN 221", "1", models.Lecture, models.Laboratory, models.Tutorial), 1+5+1)

	assert.Empty(ds.GetSections("CPEN 221", "2", models.Lecture))
	assert.Empty(ds.GetSections("CPEN 221", "2", models.Laboratory))
	assert.Empty(ds.GetSections("CPEN 221", "2", models.Tutorial))
	assert.Empty(ds.GetSections("CPEN 221", "2", models.Lecture, models.Laboratory, models.Tutorial))

	assert.Len(ds.GetSections("CPEN 221", "1-2", models.Lecture), 1)
	assert.Len(ds.GetSections("CPEN 221", "1-2", models.Laboratory), 5)
	assert.Len(ds.GetSections("CPEN 221", "1-2", models.Tutorial), 1)
	assert.Len(ds.GetSections("CPEN 221", "1-2", models.Lecture, models.Laboratory, models.Tutorial), 1+5+1)

	assert.Empty(ds.GetSections("bogus", "1-2", models.Lecture))
	assert.Empty(ds.GetSections("bogus", "1-2", models.Laboratory))
	assert.Empty(ds.GetSections("bogus", "1-2", models.Tutorial))
	assert.Empty(ds.GetSections("bogus", "1-2", models.Lecture, models.Laboratory, models.Tutorial))

	assert.Empty(ds.GetSections("CPEN 221", "bogus", models.Lecture))
	assert.Empty(ds.GetSections("CPEN 221", "bogus", models.Laboratory))
	assert.Empty(ds.GetSections("CPEN 221", "bogus", models.Tutorial))
	assert.Empty(ds.GetSections("CPEN 221", "bogus", models.Lecture, models.Laboratory, models.Tutorial))
}

func TestCourseExists(t *testing.T) {
	setup()
	assert := assert.New(t)
	ds := database.NewDatastore()

	assert.True(ds.CourseExists("CPSC 110"))
	assert.True(ds.CourseExists("MATH 100"))
	assert.True(ds.CourseExists("MATH 101"))
	assert.True(ds.CourseExists("CAPS 398"))

	assert.False(ds.CourseExists("cpsc 110"))
	assert.False(ds.CourseExists("bogus"))
}

func TestCourseHasSectionWithActivity(t *testing.T) {
	setup()
	assert := assert.New(t)
	ds := database.NewDatastore()

	assert.True(ds.CourseHasSectionWithActivity("CPSC 110", models.Lecture))
	assert.True(ds.CourseHasSectionWithActivity("CPSC 110", models.Laboratory))
	assert.True(ds.CourseHasSectionWithActivity("CPSC 110", models.Tutorial))

	assert.True(ds.CourseHasSectionWithActivity("CPSC 121", models.Lecture))
	assert.True(ds.CourseHasSectionWithActivity("CPSC 121", models.Laboratory))
	assert.True(ds.CourseHasSectionWithActivity("CPSC 121", models.Tutorial))

	assert.True(ds.CourseHasSectionWithActivity("CPSC 221", models.Lecture))
	assert.True(ds.CourseHasSectionWithActivity("CPSC 221", models.Laboratory))
	assert.False(ds.CourseHasSectionWithActivity("CPSC 221", models.Tutorial))

	assert.False(ds.CourseHasSectionWithActivity("CAPS 398", models.Lecture))
	assert.False(ds.CourseHasSectionWithActivity("CAPS 398", models.Laboratory))
	assert.False(ds.CourseHasSectionWithActivity("CAPS 398", models.Tutorial))

	assert.False(ds.CourseHasSectionWithActivity("bogus", models.Lecture))
	assert.False(ds.CourseHasSectionWithActivity("bogus", models.Laboratory))
	assert.False(ds.CourseHasSectionWithActivity("bogus", models.Tutorial))
}

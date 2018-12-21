package database_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/database"
	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() { database.LoadLocalDatabase("test-coursedb.json") })
	assert.NotZero(len(database.CourseDB()))
}

func TestValidCourses(t *testing.T) {
	assert := assert.New(t)
	database.LoadLocalDatabase("test-coursedb.json")
	validCourses := database.ValidCourses()
	assert.Contains(validCourses, "CPSC 121")
	assert.Contains(validCourses, "MATH 100")
	assert.Contains(validCourses, "MATH 101")
	assert.NotContains(validCourses, "cpsc 121")
}

func TestLoadLocalDatabase(t *testing.T) {
	assert := assert.New(t)
	assert.Panics(func() { database.LoadLocalDatabase("bad/path/to/database") })
}

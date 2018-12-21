package schedules_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/schedules"

	"github.com/smart-cs/scheduler-backend/database"
	"github.com/stretchr/testify/assert"
)

func setupAutocompleterTests() {
	database.LoadLocalDatabase("../database/test-coursedb.json")
}

func TestNewAutoCompleter(t *testing.T) {
	setupAutocompleterTests()
	assert := assert.New(t)

	ac := schedules.NewAutoCompleter()
	result := ac.CoursesWithPrefix("CPSC")
	assert.Contains(result, "CPSC 110")
	assert.Contains(result, "CPSC 121")
	assert.Contains(result, "CPSC 221")

	result = ac.CoursesWithPrefix("C")
	assert.NotEmpty(result)

	result = ac.CoursesWithPrefix("1")
	assert.Empty(result)

	t.Log("autocomplete should be case insensitive")
	const alphabetLower = "abcdefghijklmnopqrstuvwxyz"
	const alphabetUpper = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	for i := range alphabetLower {
		assert.ElementsMatchf(
			ac.CoursesWithPrefix(string(alphabetLower[i])),
			ac.CoursesWithPrefix(string(alphabetUpper[i])),
			"CourseWithPrefix yielded different results with %c and %c",
			alphabetLower[i],
			alphabetUpper[i],
		)
	}
}

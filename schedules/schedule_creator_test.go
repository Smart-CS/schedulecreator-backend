package schedules_test

import (
	"testing"

	"github.com/smart-cs/scheduler-backend/database"
	"github.com/smart-cs/scheduler-backend/schedules"
	"github.com/stretchr/testify/assert"
)

func setupScheduleCreatorTests() {
	database.LoadLocalDatabase("../database/test-coursedb.json")
}

type scheduleCreatorTestTable struct {
	courses         []string
	term            string
	expSchedulesLen int
	expCoursesLen   int
}

var defaultTestTables = []scheduleCreatorTestTable{
	// Regular cases.
	{[]string{"APSC 201"}, "1-2", 7, 1},
	{[]string{"ASIA 100"}, "1-2", 2, 1},
	{[]string{"BIOL 111"}, "1-2", 2, 1},
	{[]string{"MATH 220"}, "1-2", 9, 1},
	{[]string{"MATH 253"}, "1-2", 6, 1},
	{[]string{"MATH 335"}, "1-2", 2, 1},
	{[]string{"MATH 220", "MATH 253"}, "1-2", 54, 2},
	{[]string{"MATH 220", "MATH 335"}, "1-2", 18, 2},
	{[]string{"MATH 001", "MATH 101", "BIOC 202", "BIOC 203", "BIOC 304"}, "1", 0, 0},
	{[]string{"MATH 001", "MATH 101", "BIOC 202", "BIOC 203", "BIOC 304"}, "2", 0, 0},
	// Special cases.
	{[]string{"APSC 210"}, "1-2", 0, 0}, // Co-op placement.
	{[]string{"non-existent-course 101"}, "1-2", 0, 0},
	{[]string{"BIOL 111", "non-existent-course 101"}, "1-2", 2, 1},
	{[]string{"non-existent-course 101", "BIOL 111"}, "1-2", 2, 1},
}

func assertTables(assert *assert.Assertions, testTables []scheduleCreatorTestTable, selectLabsAndTutorials bool) {
	sc := schedules.NewScheduleCreator()
	for _, tt := range testTables {
		options := schedules.ScheduleSelectOptions{
			Term: tt.term,
			SelectLabsAndTutorials: selectLabsAndTutorials,
		}
		schedules := sc.Create(tt.courses, options)
		assert.Equalf(
			tt.expSchedulesLen, len(schedules),
			"creating schedules from %v should return %d schedules, but got %d",
			tt.courses, tt.expSchedulesLen, len(schedules),
		)
		for _, schedule := range schedules {
			if !assert.Equalf(
				tt.expCoursesLen, len(schedule.Courses),
				"schedule %v should contain %d courses, but got %d",
				schedule, tt.expCoursesLen, len(schedule.Courses)) {
				break
			}
		}
	}
}

func TestScheduleCreator_CreateDefault(t *testing.T) {
	setupScheduleCreatorTests()
	testTables := append(defaultTestTables, []scheduleCreatorTestTable{
		{[]string{"APBI 260"}, "1-2", 1, 1},
		{[]string{"CPEN 221"}, "1-2", 1, 1},
		{[]string{"CPSC 110"}, "1-2", 8, 1},
		{[]string{"CPSC 210"}, "1-2", 6, 1},
		{[]string{"CPSC 221"}, "1-2", 5, 1},
		{[]string{"CPEN 221"}, "1", 1, 1},
		{[]string{"CPEN 221"}, "2", 0, 0},
		{[]string{"APBI 260", "ASIA 100"}, "1-2", 1, 2},
		{[]string{"CPSC 221", "CPSC 121"}, "1-2", 29, 2},
		{[]string{"MATH 001", "MATH 101", "BIOC 202", "BIOC 203", "BIOC 304"}, "1-2", 40, 5},
	}...)
	assertTables(assert.New(t), testTables, false)
}

func TestScheduleCreator_CreateWithLabsAndTuts(t *testing.T) {
	setupScheduleCreatorTests()
	testTables := append(defaultTestTables, []scheduleCreatorTestTable{
		{[]string{"APBI 260"}, "1-2", 1, 2},
		{[]string{"CPEN 221"}, "1-2", 5, 3},
		{[]string{"CPSC 110"}, "1-2", 81, 3},
		{[]string{"CPSC 210"}, "1-2", 99, 2},
		{[]string{"CPSC 221"}, "1-2", 72, 2},
		{[]string{"CPSC 221", "CPSC 121"}, "1-2", 64345, 5},
		{[]string{"APBI 260", "ASIA 100"}, "1-2", 0, 0},
		{[]string{"MATH 001", "MATH 101", "BIOC 202", "BIOC 203", "BIOC 304"}, "1-2", 102, 6},
	}...)
	assertTables(assert.New(t), testTables, true)
}

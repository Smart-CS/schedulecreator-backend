package database

import (
	"bufio"
	"encoding/json"
	"os"
)

const defaultDatabasePath = "database/coursedb.json"

var courseDB *CourseDatabase

// Section is a Section of a UBC course.
type Section struct {
	Activity  []string `json:"activity"`
	Days      []string `json:"days"`
	EndTime   []string `json:"end_time"`
	Interval  string   `json:"interval"`
	StartTime []string `json:"start_time"`
	Status    string   `json:"status"`
	Term      []string `json:"term"`
}

// CourseDatabase is the schema for our courses database.
// Schema: DEPARTMENT_NAME -> COURSE_NAME -> COURSE_SECTION_NAME -> Section
// e.g. DB["CPSC"]["CPSC 121"]["CPSC 121 101"] to get the underlying Section.
type CourseDatabase map[string]map[string]map[string]interface{}

// ParseSection returns a section from an interface.
func ParseSection(section interface{}) Section {
	b, err := json.Marshal(section)
	if err != nil {
		panic(err)
	}
	var parsedSection Section
	err = json.Unmarshal(b, &parsedSection)
	if err != nil {
		panic(err)
	}
	return parsedSection
}

// ValidCourses returns the valid courses.
func ValidCourses() []string {
	var validCourses []string
	for _, sectionMap := range CourseDB() {
		for sectionName := range sectionMap {
			validCourses = append(validCourses, sectionName)
		}
	}
	return validCourses
}

// CourseDB returns the CourseDatabase.
func CourseDB() CourseDatabase {
	if courseDB == nil {
		LoadLocalDatabase(defaultDatabasePath)
	}
	return *courseDB
}

// LoadLocalDatabase loads the database from the given file path
func LoadLocalDatabase(dbPath string) {
	f, err := os.Open(dbPath)
	if err != nil {
		panic("can't initialize database")
	}

	var db CourseDatabase
	json.NewDecoder(bufio.NewReader(f)).Decode(&db)
	courseDB = &db
}

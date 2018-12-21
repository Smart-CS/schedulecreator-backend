package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/smart-cs/scheduler-backend/models"
)

// Datastore provides read operations from the datastore.
type Datastore interface {
	// GetSections returns sections of a course with one of the specified types, thats in terms.
	// Possible terms: 1, 2, 1-2.
	GetSections(courseName, term string, activityTypes ...models.ActivityType) []models.CourseSection

	// CourseExists returns if the course name exists in the datastore, case sensenitive.
	CourseExists(courseName string) bool

	// CourseHasSectionWithActivty returns true if the course has a section with the given activity type.
	CourseHasSectionWithActivity(courseName string, activity models.ActivityType) bool
}

// DefaultDatastore is the default implementation of Datastore.
type DefaultDatastore struct {
	db     CourseDatabase
	helper models.CourseHelper
}

// NewDatastore returns a Datastore leveraging an in-memory database.
func NewDatastore() Datastore {
	return &DefaultDatastore{
		db:     CourseDB(),
		helper: models.CourseHelper{},
	}
}

// GetSections returns sections of a course with one of the specified types, thats in terms.
func (ds *DefaultDatastore) GetSections(courseName, term string, activityTypes ...models.ActivityType) []models.CourseSection {
	if !ds.CourseExists(courseName) || (term != "1" && term != "2" && term != "1-2") {
		return []models.CourseSection{}
	}

	var sections []models.CourseSection
	dept := strings.Split(courseName, " ")[0]
	for sectionName, section := range ds.db[dept][courseName] {
		if !strings.HasPrefix(sectionName, courseName) {
			continue
		}
		s := ParseSection(section)
		if !ds.helper.IsIncluded(s.Activity[0], activityTypes) {
			continue
		}

		if (term == "1" || term == "2") && s.Term[0] != term {
			continue
		}

		sessions, err := ds.sessions(s)
		if err != nil {
			fmt.Printf("WARNING: failed validating fields for %q: %s\n", sectionName, err.Error())
		}

		section := models.CourseSection{
			Name:     sectionName,
			Sessions: sessions,
		}
		sections = append(sections, section)
	}
	return sections
}

// CourseExists returns if the course name is valid.
func (ds *DefaultDatastore) CourseExists(courseName string) bool {
	dept := strings.Split(courseName, " ")[0]
	_, present := ds.db[dept][courseName]
	return present
}

// CourseHasSectionWithActivity returns the courses if it has the ActivityType.
func (ds *DefaultDatastore) CourseHasSectionWithActivity(courseName string, activity models.ActivityType) bool {
	dept := strings.Split(courseName, " ")[0]
	for sectionName, section := range ds.db[dept][courseName] {
		if !strings.HasPrefix(sectionName, courseName) {
			continue
		}
		s := ParseSection(section)
		if len(s.Activity) == 0 {
			fmt.Printf("WARNING: CourseHasSectionWithActivity(%s, %s) section has no activity\n", courseName, activity)
			continue
		}

		if s.Activity[0] == activity.String() {
			return true
		}
	}
	return false
}

func (ds *DefaultDatastore) sessions(s Section) ([]models.ClassSession, error) {
	var sessions []models.ClassSession
	for i, dayStr := range s.Days {
		// dayStr looks like "Mon Wed Fri".
		for _, day := range strings.Split(dayStr, " ") {
			start, err := parseTime(s.StartTime[0])
			if err != nil {
				return nil, errors.Wrap(err, "no startTime")
			}
			end, err := parseTime(s.EndTime[0])
			if err != nil {
				return nil, errors.Wrap(err, "no endTime")
			}
			session := models.ClassSession{
				Activity: s.Activity[i],
				Term:     s.Term[i],
				Day:      day,
				Start:    start,
				End:      end,
			}
			sessions = append(sessions, session)
		}
	}
	return sessions, nil
}

// praseTime parses time in the format HH:MM to an int HHMM.
func parseTime(time string) (int, error) {
	parsed, err := strconv.Atoi(strings.Replace(time, ":", "", -1))
	return parsed, err
}

package server

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nickwu241/schedulecreator-backend/models"
)

// ScheduleCreator is the interface to create schedules.
type ScheduleCreator interface {
	Create(courses []string) []models.Schedule
}

// DefaultScheduleCreator implements ScheduleCreator.
type DefaultScheduleCreator struct {
	db CourseDatabase
}

// NewScheduleCreator constructs a new ScheduleCreator.
func NewScheduleCreator() ScheduleCreator {
	return &DefaultScheduleCreator{
		db: CourseDB(),
	}
}

// Create returns all non-conflicting schedules given a list of courses.
func (d *DefaultScheduleCreator) Create(courses []string) []models.Schedule {
	var schedules []models.Schedule
	for _, c := range courses {
		// Skip invalid courses.
		if !d.courseExists(c) {
			continue
		}
		lectureTypes := []models.ActivityType{models.Lecture, models.Seminar, models.Studio}
		schedules = d.addSections(schedules, d.createSections(c, lectureTypes))
		// schedules = d.addSections(schedules, d.createSections(c, []models.ActivityType{models.Laboratory}))
		// schedules = d.addSections(schedules, d.createSections(c, []models.ActivityType{models.Tutorial}))
	}
	return schedules
}

func (d *DefaultScheduleCreator) courseExists(course string) bool {
	dept := strings.Split(course, " ")[0]
	_, present := d.db[dept][course]
	return present
}

func (d *DefaultScheduleCreator) addSections(schedules []models.Schedule, sections []models.CourseSection) []models.Schedule {
	if len(schedules) == 0 {
		for _, section := range sections {
			sections := []models.CourseSection{section}
			schedules = append(schedules, models.Schedule{
				Courses: sections,
			})
		}
		return schedules
	}

	newSchedules := []models.Schedule{}
	for _, schedule := range schedules {
		for _, section := range sections {
			newSchedule, added := d.addSection(schedule, section)
			if added {
				newSchedules = append(newSchedules, newSchedule)
			}
		}
	}
	return newSchedules
}

// addSection returns the new schedule if it succeeds, old schedule if it fails
func (d *DefaultScheduleCreator) addSection(schedule models.Schedule, section models.CourseSection) (models.Schedule, bool) {
	newSchedule := schedule
	newSchedule.Courses = append(newSchedule.Courses, section)
	if d.conflictInSchedule(newSchedule) {
		return schedule, false
	}
	return newSchedule, true
}

func (d *DefaultScheduleCreator) createSections(course string, activityTypes []models.ActivityType) []models.CourseSection {
	// Course format i.e. CPSC 121
	var sections []models.CourseSection
	dept := strings.Split(course, " ")[0]
	// Go through all sections for this course.
	for sectionName, s := range d.db[dept][course] {
		if !d.isIncluded(s.Activity[0], activityTypes) {
			continue
		}
		// Create the sessions for each section.
		var sessions []models.ClassSession
		for i, dayStr := range s.Days {
			// dayStr looks like "Mon Wed Fri".
			for _, day := range strings.Split(dayStr, " ") {

				// TODO: refactor this logic out.
				start, err := strconv.Atoi(strings.Replace(s.StartTime[i], ":", "", -1))
				if err != nil {
					// TODO: some sections don't have a time, figure out what to do with these
					fmt.Printf("no startTime for %s: %v\n", sectionName, s)
				}
				end, err := strconv.Atoi(strings.Replace(s.EndTime[i], ":", "", -1))
				if err != nil {
					// TODO: same as above
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
		section := models.CourseSection{
			Name:     sectionName,
			Sessions: sessions,
		}
		sections = append(sections, section)
	}
	return sections
}

func (d *DefaultScheduleCreator) isIncluded(activity string, desiredTypes []models.ActivityType) bool {
	for _, a := range desiredTypes {
		if activity == a.String() {
			return true
		}
	}
	return false
}

func (d *DefaultScheduleCreator) conflictSession(s1 models.ClassSession, s2 models.ClassSession) bool {
	return s1.Term == s2.Term && s1.Day == s2.Day &&
		((s1.Start <= s2.Start && s2.Start < s1.End) ||
			(s1.Start < s2.End && s2.End <= s1.End))
}

func (d *DefaultScheduleCreator) conflictSection(s1 models.CourseSection, s2 models.CourseSection) bool {
	for _, ses1 := range s1.Sessions {
		for _, ses2 := range s2.Sessions {
			if d.conflictSession(ses1, ses2) {
				return true
			}
		}
	}
	return false
}

func (d *DefaultScheduleCreator) conflictInSchedule(schedule models.Schedule) bool {
	for _, c1 := range schedule.Courses {
		for _, c2 := range schedule.Courses {
			if c1.Name != c2.Name && d.conflictSection(c1, c2) {
				return true
			}
		}
	}
	return false
}

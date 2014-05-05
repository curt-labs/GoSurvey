package survey

import (
	"time"
)

var (
	conn          = `root:@tcp(127.0.0.1:3306)/CurtDev?parseTime=true&loc=America%2FChicago`
	getAllSurveys = `select id, name, description,
										date_added, date_modifed, userID, deleted
										from Survey order by date_modifed desc`
	getSurveyRevisions = ``
)

type Survey struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	DateAdded    time.Time        `json:"date_added"`
	DateModified time.Time        `json:"date_modified"`
	UserID       int              `json:"-"`
	Deleted      bool             `json:"-"`
	Revisions    []SurveyRevision `json:"revisions"`
	Questions    []Questions      `json:"questions"`
	Completion   SurveyStatus     `json:"-"`
}

type SurveyStatus struct {
	Completed     bool `json:"completed"`
	QuestionCount int  `json:"question_count"`
	CorrectCount  int  `json:"correct_count"`
}

type SurveyRevision struct {
}

// GetSurveys will return a list of surveys in
// the database or an error if empty.
func GetSurveys() ([]Survey, error) {

	return make([]Survey, 0), nil
}

// Get will update the values on the bound
// Survey or return an error.
func (s *Survey) Get() error {
	return nil
}

// revisions will assign revision
// history to the bound Survey.
func (s *Survey) revisions() error {
	return nil
}

// Add will commit the current Survey
// to the database.
func (s *Survey) Add() error {
	return nil
}

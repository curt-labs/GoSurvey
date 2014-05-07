package survey

import (
	"database/sql"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	getSurveyQuestions = `select id, question, date_modified,
												date_added, userID, deleted
												from SurveyQuestion
												where surveyID = ? && deleted = 0
												order by date_added`
	getQuestionRevisions = `select qr.ID as revisionID, qr.new_question,
													qr.old_question, qr.date, qr.changeType,
													u.id as userID, u.fname, u.lname, u.username
													from SurveyQuestion_Revisions as qr
													join admin.user as u on qr.userID = u.id
													where qr.questionID = ?
													order by date desc`
)

type Question struct {
	ID           int                `json:"id"`
	Question     string             `json:"question"`
	DateModified time.Time          `json:"date_modified"`
	DateAdded    time.Time          `json:"date_added"`
	UserID       int                `json:"-"`
	Deleted      bool               `json:"-"`
	Revisions    []QuestionRevision `json:"revisions"`
	Answers      []Answer           `json:"answers"`
}

type QuestionRevision struct {
	ID          int          `json:"id"`
	User        RevisionUser `json:"user"`
	NewQuestion string       `json:"new_question"`
	OldQuestion string       `json:"old_question"`
	Date        time.Time    `json:"date"`
	ChangeType  string       `json:"change_type"`
}

func (s *Survey) questions() error {
	s.Questions = make([]Question, 0)
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getSurveyQuestions)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(s.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var q Question
		err = res.Scan(&q.ID, &q.Question, &q.DateModified, &q.DateAdded, &q.UserID, &q.Deleted)
		if err == nil {
			q.revisions()
			s.Questions = append(s.Questions, q)
		}
	}

	return nil
}

func (q *Question) revisions() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getQuestionRevisions)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(q.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var qr QuestionRevision
		err = res.Scan(&qr.ID, &qr.NewQuestion, &qr.OldQuestion, &qr.Date,
			&qr.ChangeType, &qr.User.ID, &qr.User.FirstName, &qr.User.LastName,
			&qr.User.Username)
		if err == nil {
			q.Revisions = append(q.Revisions, qr)
		}
	}

	return nil
}

package survey

import (
	"database/sql"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var (
	getQuestionAnswers = `select id, answer, date_modified,
												date_added, userID, deleted
												from SurveyAnswer
												where questionID = ? && deleted = 0
												order by date_added`
	getAnswerRevisions = `select ar.ID as revisionID, qr.new_answer,
													ar.old_answer, ar.date, ar.changeType,
													u.id as userID, u.fname, u.lname, u.username
													from SurveyAnswer_Revisions as ar
													join admin.user as u on ar.userID = u.id
													where ar.answerID = ?
													order by date desc`
)

type Answer struct {
	ID           int              `json:"id"`
	Answer       string           `json:"answer"`
	DataType     string           `json:"data_type"`
	DateAdded    time.Time        `json:"date_added"`
	DateModified time.Time        `json:"date_modified"`
	UserID       int              `json:"-"`
	Deleted      bool             `json:"-"`
	Revisions    []AnswerRevision `json:"revisions"`
}

type AnswerRevision struct {
	ID         int          `json:"id"`
	User       RevisionUser `json:"user"`
	NewAnswer  string       `json:"new_answer"`
	OldAnswer  string       `json:"old_answer"`
	Date       time.Time    `json:"date"`
	ChangeType string       `json:"change_type"`
}

func (q *Question) answers() error {
	q.Answers = make([]Answer, 0)
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getQuestionAnswers)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(q.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var a Answer
		err = res.Scan(&a.ID, &a.Answer, &a.DateModified, &a.DateAdded, &a.UserID, &a.Deleted)
		if err == nil {
			a.revisions()
			q.Answers = append(q.Answers, a)
		}
	}

	return nil
}

func (a *Answer) revisions() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getAnswerRevisions)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Query(a.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var ar AnswerRevision
		err = res.Scan(&ar.ID, &ar.NewAnswer, &ar.OldAnswer, &ar.Date,
			&ar.ChangeType, &ar.User.ID, &ar.User.FirstName, &ar.User.LastName,
			&ar.User.Username)
		if err == nil {
			a.Revisions = append(a.Revisions, ar)
		}
	}
	return nil
}

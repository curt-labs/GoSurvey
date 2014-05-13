package survey

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoSurvey/helpers/database"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

// SQL Statements
var (
	getSurveyQuestions = `select id, question, date_modified,
												date_added, userID, deleted
												from SurveyQuestion
												where surveyID = ? && deleted = 0
												order by date_added`
	getQuestionRevisions = `select qr.ID as revisionID, IFNULL(qr.new_question, ""),
													IFNULL(qr.old_question, ""), qr.date, qr.changeType,
													u.id as userID, u.fname, u.lname, u.username
													from SurveyQuestion_Revisions as qr
													join admin.user as u on qr.userID = u.id
													where qr.questionID = ?
													order by date desc`
	insertQuestion = `insert into SurveyQuestion(question, date_added, userID, surveyID)
										values(?,NOW(),?,?)`
	updateQuestion = `update SurveyQuestion
										set question = ?, userID = ?, surveyID = ?
										where id = ?`
	deleteQuestion = `update SurveyQuestion
										set deleted = 0, userID = ?
										where id = ?`
)

// Question contains information for a question
// on a survey. It contains answers and revision history
// for both the question an each answer.
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

// QuestionRevision is a change record for a Question.
type QuestionRevision struct {
	ID          int          `json:"id"`
	User        RevisionUser `json:"user"`
	NewQuestion string       `json:"new_question"`
	OldQuestion string       `json:"old_question"`
	Date        time.Time    `json:"date"`
	ChangeType  string       `json:"change_type"`
}

// AddQuestion will commit a new question to a
// referenced Survey.
func (s *Survey) AddQuestion(q Question) (Question, error) {
	var err error

	if q.Question == "" {
		return q, errors.New("question cannot be blank")
	}
	if q.UserID == 0 {
		return q, errors.New("user reference not found")
	}

	if q.ID == 0 {
		err = q.insert(s.ID)
	} else {
		err = q.update(s.ID)
	}

	if err != nil {
		return q, err
	}

	for i, answer := range q.Answers {
		if err = q.AddAnswer(answer); err != nil {
			q.Answers = append(q.Answers[:i], q.Answers[i+1:]...)
			return q, err
		}
	}

	return q, err
}

// Delete will mark the referenced Question
// as deleted.
func (q *Question) Delete() error {

	if q.ID == 0 {
		return errors.New("cannot delete a question that doesn't exist")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(deleteQuestion)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(q.UserID, q.ID)
	return err
}

// insert will insert a new Question and bint it
// to the given Survey.
func (q *Question) insert(surveyID int) error {

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(insertQuestion)
	if err != nil {
		return err
	}

	res, err := stmt.Exec(q.Question, q.UserID, surveyID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	q.ID = int(id)

	return nil
}

// update will update the question, userID, and surveyID
// properties for the given Question.
func (q *Question) update(surveyID int) error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	stmt, err := db.Prepare(updateQuestion)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(q.Question, q.UserID, surveyID, q.ID)

	return err
}

// questions will retrieve all possible questions
// for the referenced Survey.
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
			q.answers()
			s.Questions = append(s.Questions, q)
		}
	}

	return nil
}

// revisions will retrieve all revision history
// for the referenced Question and push them onto
// the Revisions slice.
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

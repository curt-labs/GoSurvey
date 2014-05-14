package prize

import (
	"database/sql"
	"errors"
	"github.com/curt-labs/GoSurvey/helpers/database"
	"github.com/curt-labs/GoSurvey/models/survey"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)

type Prize struct {
	ID           int                `json:"id"`
	Title        string             `json:"name"`
	Part         string             `json:"-"`
	Description  string             `json:"description"`
	Image        url.URL            `json:"image"`
	Winner       *survey.SurveyUser `json:"winner"`
	DateAdded    time.Time          `json:"date_added"`
	DateModified time.Time          `json:"date_modified"`
	UserID       int                `json:"-"`
	Deleted      bool               `json:"-"`
	Revisions    []PrizeRevision    `json:"revisions"`
}

type PrizeRevision struct {
	ID         int                 `json:"id"`
	NewTitle   *string             `json:"new_title"`
	OldTitle   *string             `json:"old_title"`
	Date       time.Time           `json:"date"`
	ChangeType string              `json:"change_type"`
	User       survey.RevisionUser `json:"user"`
}

func All() ([]Prize, error) {
	var prizes []Prize
	var err error

	return prizes, err
}

func GetCurrent() (Prize, error) {
	var prize Prize
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return prize, err
	}
	defer db.Close()

	err = errors.New("not implemented")
	return prize, err
}

func (p *Prize) Winner() (survey.SurveyUser, error) {
	var user survey.SurveyUser
	var err error

	err = errors.New("not implemented")
	return user, err
}

func (p *Prize) Add() error {
	return errors.New("not implemented")
}

func (p *Prize) Delete() error {
	return errors.New("not implemented")
}

func (p *Prize) insert() error {
	return errors.New("not implemented")
}

func (p *Prize) update() error {
	return errors.New("not implemented")
}

func (p *Prize) revisions() error {
	return errors.New("not implemented")
}

package prize

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/curt-labs/GoSurvey/helpers/database"
	"github.com/curt-labs/GoSurvey/models/survey"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"net/url"
	"time"
)

var (
	allPrizes = `select
								sp.id, sp.title, sp.part, sp.description, sp.image,
								sp.date_added, sp.date_modified, sp.userID, sp.deleted, sp.current,
								IFNULL(su.id,0) as winnerID, IFNULL(su.fname,""), IFNULL(su.lname, ""), IFNULL(su.email, "")
								from SurveyPrize as sp
								left join SurveyUser as su on sp.winner = su.id
								where sp.deleted = 0
								order by date_modified desc
								limit ?, ?`
	getCurrentPrize = `select
											sp.id, sp.title, sp.part, sp.description, sp.image,
											sp.date_added, sp.date_modified, sp.userID, sp.deleted, sp.current
											from SurveyPrize as sp
											where sp.deleted = 0 && sp.current = 1
											limit 1`
	getPrize = `select
								sp.title, sp.part, sp.description, sp.image,
								sp.date_added, sp.date_modified, sp.userID, sp.deleted, sp.current,
								IFNULL(su.id,0) as winnerID, IFNULL(su.fname,""), IFNULL(su.lname, ""), IFNULL(su.email, "")
								from SurveyPrize as sp
								left join SurveyUser as su on sp.winner = su.id
								where sp.id = ?
								limit 1`
	getCurrentUsers = `select su.id, su.fname, su.lname, su.email
											from SurveyUser as su
											where su.date_added > ? && su.date_added < ?`
	setWinner         = `update SurveyPrize set winner = ? where id = ?`
	getPrizeRevisions = `select
												sr.id as revisionID, sr.changeType, sr.date,
												IFNULL(sr.new_title,""), IFNULL(sr.old_title, ""),
												IFNULL(sr.new_description, ""), IFNULL(sr.old_description, ""),
												IFNULL(sr.new_image, ""), IFNULL(sr.old_image, ""),
												u.ID as userID, u.fname, u.lname, u.username
												from SurveyPrize_Revisions as sr
												join admin.user as u on sr.userID = u.id
												where sr.prizeID = ?
												order by date desc`
	insertPrize = `insert into SurveyPrize(title, description, image, date_added, current, part, userID)
								values(?,?,?,NOW(),?,?,?)`
	updatePrize = `update SurveyPrize
									set title = ?, description = ?, image = ?,
									current = ?, part = ?, userID = ?
									where id = ?`
	deletePrize = `update SurveyPrize set deleted = 1, userID = ? where id = ?`
)

type Prize struct {
	ID           int               `json:"id"`
	Title        string            `json:"name"`
	Part         int               `json:"-"`
	Description  string            `json:"description"`
	Image        *url.URL          `json:"image"`
	Winner       survey.SurveyUser `json:"winner"`
	DateAdded    time.Time         `json:"date_added"`
	DateModified time.Time         `json:"date_modified"`
	UserID       int               `json:"-"`
	Deleted      bool              `json:"-"`
	Current      bool              `json:"current"`
	Revisions    []PrizeRevision   `json:"revisions"`
}

type PrizeRevision struct {
	ID             int                 `json:"id"`
	NewTitle       *string             `json:"new_title"`
	OldTitle       *string             `json:"old_title"`
	NewDescription *string             `json:"new_description"`
	OldDescription *string             `json:"old_description"`
	NewImage       *url.URL            `json:"new_image"`
	OldImage       *url.URL            `json:"old_image"`
	Date           time.Time           `json:"date"`
	ChangeType     string              `json:"change_type"`
	User           survey.RevisionUser `json:"user"`
}

func All(skip, take int) ([]Prize, error) {
	if take == 0 {
		take = 1000
	}

	if skip > 0 {
		skip = (skip - 1) * take
	}

	var prizes []Prize
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return prizes, err
	}
	defer db.Close()

	stmt, err := db.Prepare(allPrizes)
	if err != nil {
		return prizes, err
	}
	defer stmt.Close()

	res, err := stmt.Query(skip, take)
	if err != nil {
		return prizes, err
	}

	for res.Next() {
		var p Prize
		var urlStr string
		err = res.Scan(&p.ID, &p.Title, &p.Part, &p.Description,
			&urlStr, &p.DateAdded, &p.DateModified, &p.UserID, &p.Deleted, &p.Current,
			&p.Winner.ID, &p.Winner.FirstName, &p.Winner.LastName, &p.Winner.Email)
		if err == nil {
			if urlStr != "" {
				p.Image, _ = url.Parse(urlStr)
			}
			p.revisions()
			prizes = append(prizes, p)
		}
	}

	return prizes, err
}

func Current() (Prize, error) {
	var p Prize
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return p, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurrentPrize)
	if err != nil {
		return p, err
	}
	defer stmt.Close()

	var urlStr string
	err = stmt.QueryRow().Scan(&p.ID, &p.Title, &p.Part, &p.Description,
		&urlStr, &p.DateAdded, &p.DateModified, &p.UserID, &p.Deleted, &p.Current)
	if urlStr != "" {
		p.Image, _ = url.Parse(urlStr)
	}
	p.revisions()

	return p, err
}

func (p *Prize) Get() error {

	if p.ID == 0 {
		return errors.New("invalid prize reference")
	}

	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(getPrize)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var urlStr string
	err = stmt.QueryRow(p.ID).Scan(&p.Title, &p.Part, &p.Description,
		&urlStr, &p.DateAdded, &p.DateModified, &p.UserID, &p.Deleted, &p.Current,
		&p.Winner.ID, &p.Winner.FirstName, &p.Winner.LastName, &p.Winner.Email)
	if err != nil {
		return err
	}

	if urlStr != "" {
		p.Image, _ = url.Parse(urlStr)
	}

	p.revisions()

	return nil
}

func (p *Prize) PickWinner(start, end time.Time) (survey.SurveyUser, error) {
	var user survey.SurveyUser
	var err error

	err = p.Get()
	if err != nil {
		return user, err
	}

	users, err := users(start, end)
	if err != nil {
		return user, err
	}

	user = users[rand.Intn(len(users)-1)]
	if user.ID == 0 {
		return user, errors.New("no surveys completed in the given date range")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return user, err
	}
	defer db.Close()

	stmt, err := db.Prepare(setWinner)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, p.ID)
	return user, err
}

func users(start, end time.Time) ([]survey.SurveyUser, error) {
	var users []survey.SurveyUser
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return users, err
	}
	defer db.Close()

	stmt, err := db.Prepare(getCurrentUsers)
	if err != nil {
		return users, err
	}

	res, err := stmt.Query(start, end)
	if err != nil {
		return users, err
	}

	for res.Next() {
		var u survey.SurveyUser
		err = res.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email)
		if err == nil {
			users = append(users, u)
		}
	}

	return users, nil
}

func (p *Prize) Save() error {
	if p.Title == "" {
		return errors.New("prize title cannot be blank")
	}
	if p.Description == "" {
		return errors.New("prize description cannot be blank")
	}
	if p.Image == nil {
		return errors.New("prize must have an image")
	}
	if p.UserID == 0 {
		return errors.New("invalid user reference")
	}
	if p.Current {
		// Make sure no other surveys are marked as current
		currentPrize, _ := Current()
		if currentPrize.ID > 0 {
			return errors.New(fmt.Sprintf("%s is marked as the current prize.", currentPrize.Title))
		}
	}

	if p.ID == 0 {
		return p.insert()
	}

	err := p.update()
	return err
}

func (p *Prize) Delete() error {
	if p.ID == 0 {
		return errors.New("invalid prize record")
	}

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(deletePrize)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.UserID, p.ID)
	if rows, _ := res.RowsAffected(); rows == 0 {
		return errors.New("failed to delete prize")
	}

	return err
}

func (p *Prize) insert() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(insertPrize)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(p.Title, p.Description, p.Image.String(), p.Current, p.Part, p.UserID)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)

	return nil
}

func (p *Prize) update() error {
	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(updatePrize)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(p.Title, p.Description, p.Image.String(), p.Current, p.Part, p.UserID, p.ID)

	return err
}

func (p *Prize) revisions() error {
	var err error

	db, err := sql.Open("mysql", database.ConnectionString())
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare(getPrizeRevisions)
	if err != nil {
		return err
	}

	defer stmt.Close()

	res, err := stmt.Query(p.ID)
	if err != nil {
		return err
	}

	for res.Next() {
		var r PrizeRevision
		err = res.Scan(&r.ID, &r.ChangeType, &r.Date, &r.NewTitle, &r.OldTitle,
			&r.NewDescription, &r.OldDescription, &r.NewImage, &r.OldImage,
			&r.User.ID, &r.User.FirstName, &r.User.LastName,
			&r.User.Username)
		if err == nil {
			p.Revisions = append(p.Revisions, r)
		}
	}

	return nil
}

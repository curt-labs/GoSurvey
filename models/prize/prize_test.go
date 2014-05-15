package prize

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/url"
	"testing"
	"time"
)

func TestGetters(t *testing.T) {
	Convey("Testing all Prize getters", t, func() {
		Convey("Testing All()", func() {
			prizes, err := All(0, 5)
			So(err, ShouldEqual, nil)
			So(len(prizes), ShouldBeLessThan, 6)
		})

		Convey("Test PrizeCount", func() {
			total := PrizeCount()
			So(count, ShouldNotEqual, nil)
		})

		Convey("Test Current()", func() {
			p, err := Current()
			if err != nil {
				So(err.Error(), ShouldContainSubstring, "no rows")
				So(p.ID, ShouldEqual, 0)
			} else {
				So(err, ShouldEqual, nil)
				So(p.ID, ShouldBeGreaterThan, 0)
			}
		})
		Convey("Test Get()", func() {
			var p Prize
			err := p.Get()
			So(err, ShouldNotEqual, nil)

			p.ID = 1
			err = p.Get()
			if err != nil {
				So(err.Error(), ShouldContainSubstring, "no rows")
				So(p.Title, ShouldEqual, "")
			} else {
				So(err, ShouldEqual, nil)
				So(p.Title, ShouldNotEqual, "")
			}

			Convey("Try and pick a winner", func() {
				user, err := p.PickWinner(time.Now().AddDate(0, -1, 0), time.Now())
				if err != nil {
					So(err.Error(), ShouldContainSubstring, "no rows")
					So(user.ID, ShouldEqual, 0)
				} else {
					So(err, ShouldEqual, nil)
					So(user.ID, ShouldBeGreaterThan, 0)
				}
			})

			Convey("Test Delete", func() {
				err := p.Delete()
				So(err, ShouldNotEqual, nil)
			})
		})
	})
}

func TestInsert(t *testing.T) {
	Convey("Testing Insert/Update/Delete and Getters with data", t, func() {
		var p Prize
		Convey("Testing Insert", func() {
			err := p.Save()
			So(err, ShouldNotEqual, nil)

			p.Title = "Test Prize"
			err = p.Save()
			So(err, ShouldNotEqual, nil)

			p.Title = "Test Prize"
			p.Description = "Test Description"
			err = p.Save()
			So(err, ShouldNotEqual, nil)

			p.Title = "Test Prize"
			p.Description = "Test Description"
			p.Image, _ = url.Parse("http://google.com")
			err = p.Save()
			So(err, ShouldNotEqual, nil)

			p.Title = "Test Prize"
			p.Description = "Test Description"
			p.Image, _ = url.Parse("http://google.com")
			p.UserID = 1
			err = p.Save()
			So(err, ShouldEqual, nil)

			p.Title = ""
			err = p.Save()
			So(err, ShouldNotEqual, nil)

			p.Title = "Updated Test Prize"
			p.Current = true
			err = p.Save()
			if err != nil {
				So(err, ShouldContainSubstring, "is marked")
			} else {
				So(err, ShouldEqual, nil)
			}

			p, err = Current()
			So(err, ShouldEqual, nil)
			So(p.ID, ShouldBeGreaterThan, 0)

			user, err := p.PickWinner(time.Now().AddDate(0, -1, 0), time.Now())
			if err != nil {
				So(err.Error(), ShouldContainSubstring, "no rows")
				So(user.ID, ShouldEqual, 0)
			} else {
				So(err, ShouldEqual, nil)
				So(user.ID, ShouldBeGreaterThan, 0)
			}

			_, err = All(1, 0)
			So(err, ShouldEqual, nil)

			_, err = All(0, 5)
			So(err, ShouldEqual, nil)

			err = p.Delete()
			So(err, ShouldEqual, nil)

		})
	})
}

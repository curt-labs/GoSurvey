package warranties

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestWarranties(t *testing.T) {
	Convey("Testing Warranties", t, func() {
		Convey("Test All()", func() {
			ws, err := All()
			So(err, ShouldEqual, nil)
			So(ws, ShouldNotEqual, nil)
		})

		var w Warranty
		Convey("Test Get", func() {
			err := w.Get()
			So(err, ShouldNotEqual, nil)
			So(w.ID, ShouldEqual, 0)
		})

		Convey("Test Add with no data", func() {
			err := w.Add()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Add with first name", func() {
			w.FirstName = "Test"
			err := w.Add()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Add with first/last name", func() {
			w.FirstName = "Test"
			w.LastName = "User"
			w.Part = 0
			err := w.Add()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Add with first/last and email", func() {
			w.FirstName = "Test"
			w.LastName = "User"
			w.Email = "test@example.com"
			w.Part = 1
			err := w.Add()
			So(err, ShouldNotEqual, nil)
		})
		Convey("Test Add with valid data", func() {
			w.FirstName = "Test"
			w.LastName = "User"
			w.Email = "test@example.com"
			w.Part = 11000
			err := w.Add()
			So(err, ShouldEqual, nil)
		})

		Convey("Test Get with number", func() {
			err := w.Get()
			So(err, ShouldEqual, nil)
		})

		Convey("Test Delete", func() {
			err := w.Delete()
			So(err, ShouldEqual, nil)
		})

		Convey("Test Delete with bad part number", func() {
			err := w.Delete()
			So(err, ShouldNotEqual, nil)
		})
	})
}

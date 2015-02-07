package exit

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(t *testing.T) {

	Convey("Exit", t, func() {
		exiter := Default()
		exiter = New(func(int) {})
		exiter.Exit(3)
		So(exiter.Status(), ShouldEqual, 3)
	})
}

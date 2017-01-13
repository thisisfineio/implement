package parse

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
)

var (
	gopath = os.Getenv("GOPATH")
	path = strings.Join([]string{gopath, "src", "github.com", "thisisfineio", "implement", "test_data"}, string(os.PathSeparator))
)

func TestInspect(t *testing.T) {
	Convey("We can parse a file and inspect it for its interfaces", t, func(){
		f, data, err := File(path + string(os.PathSeparator) + "data.go")
		So(err, ShouldBeNil)
		_, err = Inspect(f, data)
		So(err, ShouldBeNil)

	})
}

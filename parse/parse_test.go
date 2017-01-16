package parse

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
	"github.com/thisisfineio/implement"
	"fmt"
)

var (
	gopath = os.Getenv("GOPATH")
	path = strings.Join([]string{gopath, "src", "github.com", "thisisfineio", "implement", "test_data"}, string(os.PathSeparator))
)


var (
	m = make(map[string]string)
)

func init(){
	m["TestInterface"] = "Implementation"
}

func TestInspect(t *testing.T) {
	Convey("We can parse a file and inspect it for its interfaces", t, func(){
		f, data, err := File(path + string(os.PathSeparator) + "data.go")
		So(err, ShouldBeNil)
		signatures := Inspect(f, data)
		So(err, ShouldBeNil)

		interfaces := implement.GetInterfaces(signatures, m)

		for _, i := range interfaces {
			fmt.Println(i.String())
		}
	})
}

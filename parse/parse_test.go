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

func TestInspect(t *testing.T) {
	Convey("We can parse a file and inspect it for its interfaces", t, func(){
		f, data, err := File(path + string(os.PathSeparator) + "data.go")
		So(err, ShouldBeNil)
		signatures, err := Inspect(f, data)
		So(err, ShouldBeNil)
		m := make(map[string]string)
		m["TestInterface"] = "Implementation"
		m["IgnoredInterface"] = ""
		interfaces := make([]*implement.Interface, 0)
		for k, v := range signatures {
			i := &implement.Interface{Name: k, ImplementedName: m[k], Functions: v}
			interfaces = append(interfaces, i)
		}

		for _, i := range interfaces {
			fmt.Println(i.String())
		}
	})
}

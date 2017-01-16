package implement

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strings"
	"fmt"
	"io/ioutil"
)

const(
	localStruct = "*Struct"
	packageStruct = "*bytes.Buffer"
	primitive = "string"
	errType = "error"
)


var (
	gopath = os.Getenv("GOPATH")
	path = strings.Join([]string{gopath, "src", "github.com", "thisisfineio", "implement", "test_data"}, string(os.PathSeparator))
)

const expected = `
type Implementation struct{}

func (i *Implementation) Func() {

}
func (i *Implementation) ParamFunc(s *Struct) {

}
func (i *Implementation) ReturnFunc() string {
	return ""
}
func (i *Implementation) MultiIntParam(i int, j int) {

}
func (i *Implementation) TupleFunc() (string, error) {
	return "", nil
}
func (i *Implementation) MultiParam(s *Struct, reader io.Reader) {

}
func (i *Implementation) Combo(s *Struct, r io.Reader) (string, error) {
	return "", nil
}
func (i *Implementation) FunctionParam(f func(int) (interface{}, error)) {

}
func (i *Implementation) FunctionParamAndReturn(f func(int) (interface{}, error)) (*Struct, error) {
	return nil, nil
}
func (i *Implementation) AnonymousFuncs(f func(func(int, int) func()) func(int)) func(int) {
	return nil
}
func (i *Implementation) PackageStruct(b *bytes.Buffer) *bytes.Reader {
	return nil
}
func (i *Implementation) ValueFunc(s Struct) Struct {
	return Struct{}
}
`


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

		interfaces := GetInterfaces(signatures, m)

		for _, i := range interfaces {
			fmt.Println(i.String())
		}
	})
}

func TestLowerFirstLetter(t *testing.T) {
	Convey("We can get the first letter of a variable", t, func(){
		l := LowerFirstLetterOfVar(localStruct)
		So(l, ShouldEqual, "s")

		p := LowerFirstLetterOfVar(packageStruct)
		So(p, ShouldEqual, "b")

		p = LowerFirstLetterOfVar(primitive)
		So(p, ShouldEqual, "s")

		e := LowerFirstLetterOfVar(errType)
		So(e, ShouldEqual, "e")
	})
}


func TestInterface_Save(t *testing.T) {
	Convey("We can test saving an interface as a file", t, func(){
		f, data, err := File(path + string(os.PathSeparator) + "data.go")
		So(err, ShouldBeNil)
		signatures := Inspect(f, data)
		So(err, ShouldBeNil)

		interfaces := GetInterfaces(signatures, m)

		for _, i := range interfaces {
			fmt.Println(i.Name)
			d, err := ioutil.TempDir("/tmp", "")
			So(err, ShouldBeNil)
			err = i.Save(d)
			So(err, ShouldBeNil)
			data, err := ioutil.ReadFile(d + string(os.PathSeparator) + i.Name + ".go")
			So(err, ShouldBeNil)
			index := strings.Index(string(data), "\n")
			data =  data[index:]
			So(string(data), ShouldEqual, expected)
			err = os.RemoveAll(d)
			So(err, ShouldBeNil)
		}
	})
}
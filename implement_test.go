package implement

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const(
	localStruct = "*Struct"
	packageStruct = "*bytes.Buffer"
	primitive = "string"
	errType = "error"
)

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
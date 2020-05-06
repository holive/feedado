package urlutil

import (
	"net/url"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestClone(t *testing.T) {
	Convey("Testing clone", t, func() {
		original, err := url.Parse("http://www.americanas.com.br/")
		So(err, ShouldBeNil)

		clone := Clone(original)
		clone.Path = path.Join(original.Path, "/hotsite/blacfriday")

		So(original.String(), ShouldResemble, "http://www.americanas.com.br/")
		So(clone.String(), ShouldResemble, "http://www.americanas.com.br/hotsite/blacfriday")
	})
}

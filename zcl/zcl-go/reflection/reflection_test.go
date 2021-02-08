package reflection

import (
	"testing"

	. "gopkg.in/check.v1"
)

func TestReflection(t *testing.T) { TestingT(t) }

type ReflectionSuite struct{}

var _ = Suite(&ReflectionSuite{})

func (s *ReflectionSuite) TestCopy(c *C) {
	type Struct struct {
		V1 uint8
		V2 string
	}

	copy1 := Copy(&Struct{1, "2"})

	c.Assert(copy1, DeepEquals, &Struct{})

	copy2 := Copy(Struct{1, "2"})

	c.Assert(copy2, DeepEquals, Struct{})
}

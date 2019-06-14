package argp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoArgument(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(nil)
	assrt.NotNil(err)
	assrt.Nil(args)
	args, err = Parse(strings.NewReader(""))
	assrt.Nil(err)
	assrt.Nil(args)
	assrt.Equal(len(args), 0)
	args, err = Parse(strings.NewReader(" "))
	assrt.Nil(err)
	assrt.Nil(args)
	assrt.Equal(len(args), 0)
}
func TestArgEncapSingle(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"abc"`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`abc`, args[0])
}
func TestArgEncapSingleWithEmbeddedDblQuoteAndEscape(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"ab\\c\""`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`ab\c"`, args[0])
}
func TestArgEncapTwo(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"abc" "def"`))
	assrt.Nil(err)
	assrt.Equal(len(args), 2)
	assrt.Equal(`abc`, args[0])
	assrt.Equal(`def`, args[1])
}
func TestArgEncapSingleQuoteNoReplaceEscapeSingle(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`'a\"b\\c'`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`a\"b\\c`, args[0])
}
func TestArgEncapSingleQuoteSingle(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`'abc'`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`abc`, args[0])
}
func TestArgNotEncap(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`abc`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`abc`, args[0])
}
func TestArgNotEncapWithEmbeddedDblQuoteEscapeSingleQuote(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`\"a\\b\'c\l`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`"a\b'cl`, args[0])
}
func TestArgNotEncapTwo(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`abc def`))
	assrt.Nil(err)
	assrt.Equal(len(args), 2)
	assrt.Equal(`abc`, args[0])
	assrt.Equal(`def`, args[1])
}
func TestArgEncapSingleDoubleOne(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"abc"'def'`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`abcdef`, args[0])
}
func TestArgNotEncapStartWithWhitespaceTwo(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`abc def`))
	assrt.Nil(err)
	assrt.Equal(len(args), 2)
	assrt.Equal(`abc`, args[0])
	assrt.Equal(`def`, args[1])
}
func TestArgNotEncapAndEncapSingle(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"abc"def`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`abcdef`, args[0])
}
func TestArgNotEncapSingleFail(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`"`))
	assrt.NotNil(err)
	assrt.Nil(args)
	args, err = Parse(strings.NewReader(`'`))
	assrt.NotNil(err)
	assrt.Nil(args)
}
func TestArgNotEncapEscapeSingle(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`\\`))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`\`, args[0])
}
func TestArgNotEncapTrailingWhitespace(t *testing.T) {
	assrt := assert.New(t)
	args, err := Parse(strings.NewReader(`a                                                                                                       `))
	assrt.Nil(err)
	assrt.Equal(len(args), 1)
	assrt.Equal(`a`, args[0])
}

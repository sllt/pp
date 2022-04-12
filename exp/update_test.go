package exp

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type updateExpressionTestSuite struct {
	suite.Suite
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withInvalidValue() {
	_, err := NewUpdateExpressions(true)
	uets.EqualError(err, "pp: unsupported update interface type bool")
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withRecords() {
	ie, err := NewUpdateExpressions(Record{"c": "a", "b": "d"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "b").Set("d"),
		NewIdentifierExpression("", "", "c").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withMap() {
	ie, err := NewUpdateExpressions(map[string]interface{}{"c": "a", "b": "d"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "b").Set("d"),
		NewIdentifierExpression("", "", "c").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructs() {
	type testRecord struct {
		C string `db:"c"`
		B string `db:"b"`
	}
	ie, err := NewUpdateExpressions(testRecord{C: "a", B: "d"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "b").Set("d"),
		NewIdentifierExpression("", "", "c").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructsWithoutTags() {
	type testRecord struct {
		FieldA int64
		FieldB bool
		FieldC string
	}
	ie, err := NewUpdateExpressions(testRecord{FieldA: 1, FieldB: true, FieldC: "a"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "fielda").Set(int64(1)),
		NewIdentifierExpression("", "", "fieldb").Set(true),
		NewIdentifierExpression("", "", "fieldc").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructsIgnoredDbTag() {
	type testRecord struct {
		FieldA int64 `db:"-"`
		FieldB bool
		FieldC string
	}
	ie, err := NewUpdateExpressions(testRecord{FieldA: 1, FieldB: true, FieldC: "a"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "fieldb").Set(true),
		NewIdentifierExpression("", "", "fieldc").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructsWithPpSkipUpdate() {
	type testRecord struct {
		FieldA int64
		FieldB bool   `pp:"skipupdate"`
		FieldC string `pp:"skipinsert"`
	}
	ie, err := NewUpdateExpressions(testRecord{FieldA: 1, FieldB: true, FieldC: "a"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "fielda").Set(int64(1)),
		NewIdentifierExpression("", "", "fieldc").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructPointers() {
	type testRecord struct {
		C string `db:"c"`
		B string `db:"b"`
	}
	ie, err := NewUpdateExpressions(&testRecord{C: "a", B: "d"})
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "b").Set("d"),
		NewIdentifierExpression("", "", "c").Set("a"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructsWithEmbeddedStructs() {
	type Phone struct {
		Primary string `db:"primary_phone"`
		Home    string `db:"home_phone"`
	}
	type item struct {
		Phone
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	ie, err := NewUpdateExpressions(
		item{Address: "111 Test Addr", Name: "Test1", Phone: Phone{Home: "123123", Primary: "456456"}},
	)
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "address").Set("111 Test Addr"),
		NewIdentifierExpression("", "", "home_phone").Set("123123"),
		NewIdentifierExpression("", "", "name").Set("Test1"),
		NewIdentifierExpression("", "", "primary_phone").Set("456456"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withStructsWithEmbeddedStructPointers() {
	type Phone struct {
		Primary string `db:"primary_phone"`
		Home    string `db:"home_phone"`
	}
	type item struct {
		*Phone
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	ie, err := NewUpdateExpressions(
		item{Address: "111 Test Addr", Name: "Test1", Phone: &Phone{Home: "123123", Primary: "456456"}},
	)
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "address").Set("111 Test Addr"),
		NewIdentifierExpression("", "", "home_phone").Set("123123"),
		NewIdentifierExpression("", "", "name").Set("Test1"),
		NewIdentifierExpression("", "", "primary_phone").Set("456456"),
	}
	uets.Equal(eie, ie)
}

func (uets *updateExpressionTestSuite) TestNewUpdateExpressions_withNilEmbeddedStructPointers() {
	type Phone struct {
		Primary string `db:"primary_phone"`
		Home    string `db:"home_phone"`
	}
	type item struct {
		*Phone
		Address string `db:"address"`
		Name    string `db:"name"`
	}
	ie, err := NewUpdateExpressions(
		item{Address: "111 Test Addr", Name: "Test1"},
	)
	uets.NoError(err)
	eie := []UpdateExpression{
		NewIdentifierExpression("", "", "address").Set("111 Test Addr"),
		NewIdentifierExpression("", "", "name").Set("Test1"),
	}
	uets.Equal(eie, ie)
}

func TestUpdateExpressionSuite(t *testing.T) {
	suite.Run(t, new(updateExpressionTestSuite))
}

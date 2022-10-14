package moccasin_test

import (
	"fmt"
	"testing"

	"github.com/rainforestpay/moccasin"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	moccasin.Embed
}

func (t *TestStruct) String() string {
	if t.MMocked(true) {
		return t.MGet(0).(string) //nolint:forcetypeassert
	}
	return "default"
}

func (t TestStruct) Crazy() (string, int, bool) {
	if t.MMocked(true) {
		var r0 string
		if t.MGet(0) != nil {
			r0, _ = t.MGet(0).(string)
		}

		var r1 int
		if t.MGet(1) != nil {
			r1, _ = t.MGet(1).(int)
		}

		var r2 bool
		if t.MGet(2) != nil {
			r2, _ = t.MGet(2).(bool)
		}

		return r0, r1, r2
	}

	return "default", 0, false
}

func ExampleMockResponse_MReturn() {
	ts := &TestStruct{}
	ts.MAttach("String").MReturn("override")

	fmt.Println(ts.String())
	// Output: override
}

func ExampleMockResponse_MAddReturn() {
	ts := &TestStruct{}
	ts.MAttach("String").MReturn("firstResponse").MAddReturn("secondResponse")

	fmt.Println(ts.String())
	fmt.Println(ts.String())
	fmt.Println(ts.String())
	// Output: firstResponse
	// secondResponse
	// default
}

func ExampleMockResponse_MTimes() {
	ts := &TestStruct{}
	ts.MAttach("String").MReturn("override").MTimes(2)

	fmt.Println(ts.String())
	fmt.Println(ts.String())
	fmt.Println(ts.String())
	// Output: override
	// override
	// default
}

func TestOneArgument(t *testing.T) {
	t.Run("Override", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("String").MReturn("override")
		assert.Equal(t, "override", ts.String())
	})

	t.Run("Default", func(t *testing.T) {
		ts := &TestStruct{}
		assert.Equal(t, "default", ts.String())
	})

	t.Run("MTimes", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("String").MReturn("override").MTimes(2)
		assert.Equal(t, "override", ts.String())
		assert.Equal(t, "override", ts.String())
		assert.Equal(t, "default", ts.String())
	})
}

func TestMultiArgument(t *testing.T) {
	t.Run("Override", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("Crazy").MReturn("override", 42, true)
		a, b, c := ts.Crazy()
		assert.Equal(t, "override", a)
		assert.Equal(t, 42, b)
		assert.Equal(t, true, c)
	})

	t.Run("Default", func(t *testing.T) {
		ts := &TestStruct{}
		a, b, c := ts.Crazy()
		assert.Equal(t, "default", a)
		assert.Equal(t, 0, b)
		assert.Equal(t, false, c)
	})

	t.Run("Mismatch", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("Crazy").MReturn("override")
		a, b, c := ts.Crazy()
		assert.Equal(t, "override", a)
		assert.Equal(t, 0, b)
		assert.Equal(t, false, c)
	})
}

func TestMultiAttach(t *testing.T) {
	t.Run("Override", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("String").MReturn("stringOverride")
		ts.MAttach("Crazy").MReturn("crazyOverride", 42, true)

		assert.Equal(t, "stringOverride", ts.String())

		a, b, c := ts.Crazy()
		assert.Equal(t, "crazyOverride", a)
		assert.Equal(t, 42, b)
		assert.Equal(t, true, c)
	})

	t.Run("Default", func(t *testing.T) {
		ts := &TestStruct{}

		assert.Equal(t, "default", ts.String())

		a, b, c := ts.Crazy()
		assert.Equal(t, "default", a)
		assert.Equal(t, 0, b)
		assert.Equal(t, false, c)
	})
}

func TestReturnQueue(t *testing.T) {
	t.Run("Multiple", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("String").MAddReturn("firstTime").MAddReturn("secondTime")

		assert.Equal(t, "firstTime", ts.String())
		assert.Equal(t, "secondTime", ts.String())
		assert.Equal(t, "default", ts.String())
	})
}

func TestRemove(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := &TestStruct{}
		ts.MAttach("String").MReturn("override")
		ts.MRemove("String")

		assert.Equal(t, "default", ts.String())
	})
}

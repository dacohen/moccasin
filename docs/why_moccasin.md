# Why moccasin?
Existing mocking frameworks for go are inflexible, and opinionated. Here are two pain points that moccasin solves:

## Overriding shouldn't be Hard
Here's a trivial go function stub that we want to mock:
```go
func (m *MyMock) String() (string, error) {
	// Mocking logic here
}
```

In our test code, we will likely define a default mock that used by most of the test cases:
```go
myMock.MAttach("String").MReturn("default", nil)
```

However, for just a few test cases, we want to mock the `String()` function returning an error. In other mocking 
frameworks, setting another value for a return adds it to a call stack, making it difficult to override this default 
mock. Moccasin takes a more explicit approach, and allows you to simply write the following:

```go
myMock.MAttach("String").MReturn("default", nil)

func SpecificTestCase(t *testing.T) {
    myMock.MAttach("String").MReturn("override", errors.New("weird error"))
	// ...
}

```





## To Mock, or not to Mock?
Take a simple mock that needs to maintain state throughout a test, in order to emulate the backed it's mocking:
```go
type MockBill struct {
	payState map[string]bool
}

func (m *MockBill) Pay(billId string) error {
	m.payState[billId] = true
	return nil
}

func (m *MockBill) IsPaid(billId string) (bool, error) {
	isPaid, ok := m.payState[billId]
	if !ok {
		return false, nil
    }
	return isPaid, nil
}
```

This mock works great, until we want to mock an error occurring inside one of the functions. Existing mocking 
frameworks make this difficult, by requiring that all mocked functions always have mocks defined, but in this 
example, the vast majority of test cases don't require this.

Moccasin makes this easy:

```go
type MockBill struct {
	moccasin.Embed
	
	payState map[string]bool
}

func (m *MockBill) Pay(billId string) error {
	if m.MMocked(true) {
        err, _ := m.MGet(0).(error)
        return err
	}
	
	m.payState[billId] = true
	return nil
}

func (m *MockBill) IsPaid(billId string) (bool, error) {
    if m.MMocked(true) {
        err, _ := m.MGet(0).(error)
        return false, err
    }
	
	isPaid, ok := m.payState[billId]
	if !ok {
		return false, nil
    }
	
	return isPaid, nil
}
```

Now the old mocking logic is unaffected unless we explicitly inject an error, the best of both worlds.


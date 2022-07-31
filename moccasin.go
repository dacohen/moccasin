// Package moccasin is a comfortable, un-opinionated mocking tool.
// It is designed to easily and explicitly attach mocks that can be optionally consumed by receivers.
package moccasin

// Embed manages the mocking process for a struct. To use it, simply embed it in a struct:
//	type MyStruct struct {
//		moccasin.Embed
//
//		Name string
//		Value int64
//	}
type Embed struct {
	activeMocks map[string]*MockResponse
}

// MAttach attaches a mock for the provided method on the parent struct.
// To specify the response, see the MockResponse struct below.
func (e *Embed) MAttach(methodName string) *MockResponse {
	if e.activeMocks == nil {
		e.activeMocks = make(map[string]*MockResponse)
	}

	mr := &MockResponse{
		popReturnQueue: false,
	}
	e.activeMocks[methodName] = mr

	return mr
}

// MRemove un-registers any mocks attached to the given method
func (e *Embed) MRemove(methodName string) {
	delete(e.activeMocks, methodName)
}

// MMocked determines if there is a mock available for the calling function. In general, the called parameter should be
// set to true. However, if you want to peek at the value without affecting state, it can be set to false.
func (e *Embed) MMocked(called bool) bool {
	fnName, err := getFnName()
	if err != nil {
		return false
	}

	if mock, ok := e.activeMocks[fnName]; ok {
		if called {
			if len(mock.returnVals) > 0 {
				if mock.popReturnQueue {
					// Pop the first return val off the queue
					mock.returnVals = mock.returnVals[1:len(mock.returnVals)]
				} else if len(mock.returnVals) > 1 {
					// Only start popping the queue if we had multiple frames to begin with
					mock.popReturnQueue = true
				}
			}
		}
		return len(mock.returnVals) > 0
	}
	return false
}

// MGet returns the mocked return value at index i, if available. It should be used in conjunction with MMocked.
func (e *Embed) MGet(i int) interface{} {
	fnName, err := getFnName()
	if err != nil {
		return nil
	}

	if mock, ok := e.activeMocks[fnName]; ok {
		if len(mock.returnVals) == 0 {
			return nil
		}
		currentRetVals := mock.returnVals[0]
		if len(currentRetVals) <= i {
			return nil
		}

		return currentRetVals[i]
	}

	return nil
}

// MockResponse holds the registered response to a mocked function.
type MockResponse struct {
	returnVals     [][]interface{}
	popReturnQueue bool
}

// MTimes sets the number of times the current return value (registered with MReturn) will be returned.
// It is equivalent to adding multiple returns to the queue with MAddReturn
func (m *MockResponse) MTimes(times int) *MockResponse {
	if len(m.returnVals) == 0 {
		return m
	}
	for i := 1; i < times; i++ {
		m.MAddReturn(m.returnVals[0]...)
	}
	return m
}

// MReturn specifies the return values to register for the MockResponse. It overrides any existing return set on this
// MockResponse. To add additional returns, use MAddReturn
func (m *MockResponse) MReturn(retVals ...interface{}) *MockResponse {
	m.returnVals = [][]interface{}{retVals}
	m.popReturnQueue = false

	return m
}

// MAddReturn adds a return value to the MockResponse. It can be called multiple times to add different return values
// for different calls to the mock.
func (m *MockResponse) MAddReturn(retVals ...interface{}) *MockResponse {
	m.returnVals = append(m.returnVals, retVals)
	m.popReturnQueue = false

	return m
}

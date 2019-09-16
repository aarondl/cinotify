package cinotify

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gorilla/mux"
)

type testHandler struct {
}

func (t testHandler) Handle(r *http.Request) fmt.Stringer {
	return testNotification{}
}

func (t testHandler) Route(r *mux.Route) {
	r.Path("/")
}

type testNotification struct {
}

func (t testNotification) String() string {
	return "t"
}

type testNotifier struct {
}

func (t *testNotifier) Notify(name string, note fmt.Stringer) {
}

func TestRegister(t *testing.T) {
	name := "test"
	handler := testHandler{}
	Register(name, handler)

	h := handlers[name]
	if h.name != name {
		t.Error("It should store the name.")
	}
	if h.realHandler != handler {
		t.Error("It should register the Handler to the name.")
	}
	if h.notifiers == nil {
		t.Error("It should initialize the notifiers list.")
	}
}

func TestTo(t *testing.T) {
	handlers = make(map[string]*handler)
	h1 := testHandler{}
	Register("test1", h1)
	h2 := testHandler{}
	Register("test2", h2)

	var notifier Notifier = &testNotifier{}
	To(notifier)
	if n := handlers["test1"].notifiers; len(n) != 1 {
		t.Error("Expected test1 to have exactly one notifier, got: ", len(n))
	}
	if n := handlers["test2"].notifiers; len(n) != 1 {
		t.Error("Expected test2 to have exactly one notifier, got: ", len(n))
	}
}

func TestWhenTo(t *testing.T) {
	handlers = make(map[string]*handler)
	h1 := testHandler{}
	Register("test1", h1)
	h2 := testHandler{}
	Register("test2", h2)

	var notifier Notifier = &testNotifier{}
	When("test2").To(notifier)
	if n := handlers["test1"].notifiers; len(n) != 0 {
		t.Error("Expected test1 to have exactly one notifier, got: ", len(n))
	}
	if n := handlers["test2"].notifiers; len(n) != 1 {
		t.Error("Expected test2 to have exactly one notifier, got: ", len(n))
	}
}

func TestDispatch(t *testing.T) {
	handlers = make(map[string]*handler)
	h1 := testHandler{}
	h2 := testHandler{}
	Register("test1", h1)
	Register("test2", h2)

	var n, s string
	var count = 0
	To(NotifyFunc(func(name string, notification fmt.Stringer) {
		n, s = name, notification.String()
		count++
	}))

	if len(n) > 0 || len(s) > 0 {
		t.Error("Test set up is strange.")
	}
	dispatch("test1", testNotification{})
	if n != "test1" {
		t.Error("Expected name to be test1 but was: ", n)
	}
	if s != (testNotification{}.String()) {
		t.Errorf("Expected s to be: %v, but got: %v", testNotification{}, s)
	}
	if 1 != count {
		t.Errorf("Expected count to be: 1, but got: %v", count)
	}
}

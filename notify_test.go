package cinotify

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	. "testing"
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

func TestRegister(t *T) {
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
	if h.notifyfuncs == nil {
		t.Error("It should initialize the notifyfuncs list.")
	}
}

func TestHandler_Route(t *T) {
	var handler Handler = testHandler{}

	router := mux.NewRouter()

	response := httptest.NewRecorder()
	body := bytes.NewBufferString("body")
	request, err := http.NewRequest("POST", "/", body)

	if err != nil {
		t.Fatal("Could not create mock request.")
	}

	router.ServeHTTP(response, request)
	if response.Code != http.StatusNotFound {
		t.Error("There should be no routes registered!")
	}

	route := router.NewRoute()
	handler.Route(route)
	route.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
	})

	response = httptest.NewRecorder()
	body = bytes.NewBufferString("body")
	request, err = http.NewRequest("POST", "/", body)

	if err != nil {
		t.Fatal("Could not create mock request.")
	}

	router.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Error("There is no route registered!")
	}
}

func TestHandler_Handle(t *T) {
	var handler Handler = testHandler{}
	var note fmt.Stringer = handler.Handle(nil)

	if note == nil {
		t.Error("Expected Handle to return a stringer.")
	}
}

func TestNotification_String(t *T) {
	var note fmt.Stringer = testNotification{}
	expect := testNotification{}.String()

	if s := note.String(); s != expect {
		t.Errorf("Expected: %s, got: %s", expect, s)
	}
}

func TestTo(t *T) {
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

	ToFunc(notifier.Notify)
	if n := handlers["test1"].notifyfuncs; len(n) != 1 {
		t.Error("Expected test1 to have exactly one notifyfunc, got: ", len(n))
	}
	if n := handlers["test2"].notifyfuncs; len(n) != 1 {
		t.Error("Expected test2 to have exactly one notifyfunc, got: ", len(n))
	}
}

func TestWhenTo(t *T) {
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

	When("test1").ToFunc(notifier.Notify)
	if n := handlers["test1"].notifyfuncs; len(n) != 1 {
		t.Error("Expected test1 to have exactly one notifyfunc, got: ", len(n))
	}
	if n := handlers["test2"].notifyfuncs; len(n) != 0 {
		t.Error("Expected test2 to have exactly one notifyfunc, got: ", len(n))
	}
}

func TestCreateRouter(t *T) {
	handlers = make(map[string]*handler)
	h1 := testHandler{}
	Register("test1", h1)

	router := createRouter()
	if router.Get("test1").GetHandler() == nil {
		t.Error("It should have hooked up the given route to a http handler.")
	}
}

func TestDispatch(t *T) {
	handlers = make(map[string]*handler)
	h1 := testHandler{}
	Register("test1", h1)

	var n, s string
	ToFunc(func(name string, notification fmt.Stringer) {
		n, s = name, notification.String()
	})

	router := createRouter()
	if router.Get("test1").GetHandler() == nil {
		t.Error("It should have hooked up the given route to a http handler.")
	}

	if len(n) > 0 || len(s) > 0 {
		t.Error("Test set up is strange.")
	}
	dispatch(testNotification{})
	if n != "test1" {
		t.Error("Expected name to be test1 but was: ", n)
	}
	if s != (testNotification{}.String()) {
		t.Errorf("Expected s to be: %v, but got: %v", testNotification{}, s)
	}
}

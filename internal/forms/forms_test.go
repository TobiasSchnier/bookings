package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")

	if form.Valid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows invalid when required fields are there")
	}
}
func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// invalid
	has := form.Has("a")
	if has {
		t.Error("form shows has field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("form shows invalid when required fields are there")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	// invalid
	form.MinLength("a", 4)
	if form.Valid() {
		t.Error("form shows valid when request has no a")
	}

	isError := form.Errors.Get("a")
	if isError == "" {
		t.Error("should have an error but doesnt")

	}

	postedData := url.Values{}
	postedData.Add("a", "abcd")

	form = New(postedData)
	form.MinLength("a", 100)
	if form.Valid() {
		t.Error("form shows invalid when length is correct")
	}

	postedData = url.Values{}
	postedData.Add("b", "abcd")

	form = New(postedData)
	form.MinLength("b", 1)
	if !form.Valid() {
		t.Error("form shows valid when length is not correct")
	}
	isError = form.Errors.Get("b")
	if isError != "" {
		t.Error("should have gotten no error but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)

	// invalid
	form.IsEmail("a")
	if form.Valid() {
		t.Error("form shows valid when request is not an email")
	}

	postedData = url.Values{}
	postedData.Add("a", "a@gmail.com")
	form = New(postedData)
	form.IsEmail("a")
	if !form.Valid() {
		t.Error("form shows invalid when mail is ok")
	}
	postedData = url.Values{}
	postedData.Add("b", "ab.com")
	form = New(postedData)
	form.IsEmail("b")
	if form.Valid() {
		t.Error("form shows valid when mail is not ok")
	}
}

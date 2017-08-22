package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func describe(desc string, expected, actual interface{}) string {
	return fmt.Sprintf("%s\nExpected: %v\nbut got: %v", desc, expected, actual)
}

func Test_givenTextAround_whenFindLeft_thenWritesLeftText(t *testing.T) {

	expected := "{before"

	r := bytes.NewReader([]byte(expected + "{{after"))
	var w bytes.Buffer
	p := NewParser(r, &w)

	err := p.findLeftCurlies()
	if err != nil {
		t.Fatal(err)
	}

	actual := string(w.Bytes())
	if actual != expected {
		t.Fatal(describe("Sink has wrong content", expected, actual))
	}
}

func Test_givenTextAround_whenCutRight_thenReturnsTextBefore(t *testing.T) {

	expectedCut := "}{{before"

	r := bytes.NewReader([]byte(expectedCut + "}}after"))
	var w bytes.Buffer
	p := NewParser(r, &w)

	actualCut, err := p.cutTillRightCurlies()
	if err != nil {
		t.Fatal(err)
	}

	actualSink := string(w.Bytes())
	if actualSink != "" {
		t.Fatal(describe("Sink should be empty", "", actualSink))
	}

	if actualCut != expectedCut {
		t.Fatal(describe("Wrong text was cut", expectedCut, actualCut))
	}
}

func Test_givenTextAroundHandlebarNode_whenParseNext_thenReceiveHandlebarText(t *testing.T) {

	expectedHandlebar := "inside"
	expectedSink := "before"

	r := bytes.NewReader([]byte(expectedSink + "{{" + expectedHandlebar + "}}after"))
	var w bytes.Buffer
	p := NewParser(r, &w)

	actualHandlebar, _, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	if actualHandlebar != expectedHandlebar {
		t.Fatal(describe("Wrong handlebar text", expectedHandlebar, actualHandlebar))
	}

	actualSink := string(w.Bytes())
	if actualSink != expectedSink {
		t.Fatal(describe("Sink has wrong content", expectedSink, actualSink))
	}
}

func Test_givenMultipleHandlebarNodes_whenParsed_allAreFound(t *testing.T) {

	r := bytes.NewReader([]byte("asdf{{cap1}}safd\nfdas{{cap2}}fdfda"))
	var w bytes.Buffer
	p := NewParser(r, &w)

	actual := []string{}
	for {
		s, cont, err := p.Next()
		if err != nil {
			t.Fatal(err)
		}
		if !cont {
			break
		}
		actual = append(actual, s)
	}

	expected := []string{"cap1", "cap2"}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatal(describe("Handlebar nodes are wrong", expected, actual))
	}
}

func Test_givenTextWithoutHandlebars_whenParsed_thenReturnsFalse(t *testing.T) {

	expectedSink := "someTextWithoutHandlebars"

	r := bytes.NewReader([]byte(expectedSink))
	var w bytes.Buffer
	p := NewParser(r, &w)

	hb, found, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	if found {
		t.Fatal("Found handlebar node", hb)
	}

	actualSink := string(w.Bytes())
	if actualSink != expectedSink {
		t.Fatal(describe("Sink has wrong content", expectedSink, actualSink))
	}
}

func Test_givenTextWithSpuriousLeftCurlies_whenParsed_thenReturnsFalse(t *testing.T) {

	expectedSink := "someText{{WithoutHandlebars"

	r := bytes.NewReader([]byte(expectedSink))
	var w bytes.Buffer
	p := NewParser(r, &w)

	hb, found, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	if found {
		t.Fatal("Found handlebar node", hb)
	}

	actualSink := string(w.Bytes())
	if actualSink != expectedSink {
		t.Fatal(describe("Sink has wrong content", expectedSink, actualSink))
	}
}

func Test_givenNewlineFalsePositive_whenParsed_thenIgnored(t *testing.T) {

	expected := "before{{inside\n}}after"

	r := bytes.NewReader([]byte(expected))
	var w bytes.Buffer
	p := NewParser(r, &w)

	falsePositive, found, err := p.Next()
	if err != nil {
		t.Fatal(err)
	}

	if found {
		t.Fatal("Found handlebar node:", falsePositive)
	}

	actual := string(w.Bytes())
	if actual != expected {
		t.Fatal(describe("Sink has wrong content", expected, actual))
	}
}

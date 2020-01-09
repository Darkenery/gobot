package test

import (
	"github.com/darkenery/gobot/bot/util"
	"testing"
)

func TestToLowerCase(t *testing.T) {
	text := "I AM TEST STRING"
	result := util.ToLowerCase(text)

	if result != "i am test string" {
		t.Error("Expected 'i am test string', got ", result)
	}
}

func TestRemoveWhitespace(t *testing.T) {
	text := `I AM
TEST
STRING`
	result := util.RemoveWhitespace(text)

	if result != "I AM TEST STRING" {
		t.Error("Expected 'I AM TEST STRING', got ", result)
	}
}

func TestRemoveNonWordSymbols(t *testing.T) {
	text := `I. AM, T3ST; STRING! ХАХА!`
	result := util.RemoveNonWordSymbols(text)

	if result != "I AM T3ST STRING ХАХА" {
		t.Error("Expected 'I AM T3ST STRING ХАХА', got ", result)
	}
}

func TestUcFirst(t *testing.T) {
	text := "я am test string"
	result := util.UcFirst(text)

	if result != "Я am test string" {
		t.Error("Expected 'Я am test string', got ", result)
	}
}

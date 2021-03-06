package vartypes

import (
	"errors"
	"strconv"

	"github.com/xiaq/persistent/hashmap"
	"github.com/xiaq/persistent/vector"
)

var (
	errShouldBeList   = errors.New("should be list")
	errShouldBeMap    = errors.New("should be map")
	errShouldBeBool   = errors.New("should be bool")
	errShouldBeNumber = errors.New("should be number")
)

func ShouldBeList(v interface{}) error {
	if _, ok := v.(vector.Vector); !ok {
		return errShouldBeList
	}
	return nil
}

func ShouldBeMap(v interface{}) error {
	if _, ok := v.(hashmap.Map); !ok {
		return errShouldBeMap
	}
	return nil
}

func ShouldBeBool(v interface{}) error {
	if _, ok := v.(bool); !ok {
		return errShouldBeBool
	}
	return nil
}

func ShouldBeNumber(v interface{}) error {
	if _, ok := v.(string); !ok {
		return errShouldBeNumber
	}
	_, err := strconv.ParseFloat(string(v.(string)), 64)
	return err
}

package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// Validate checks every exported field of the struct v, or of the
// struct v points to, against the rules of its check tag, and goes
// into nested structs; it joins every problem it finds.
func Validate(v any) error {
	s := reflect.Indirect(reflect.ValueOf(v))
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("validate: want a struct, got %T", v)
	}
	var problems []error
	for f, fv := range s.Fields() {
		if !f.IsExported() {
			continue
		}
		name := s.Type().Name() + "." + f.Name
		rules := f.Tag.Get("check")
		for rule := range strings.SplitSeq(rules, ",") {
			if rule != "" {
				err := checkRule(name, fv, rule)
				problems = append(problems, err)
			}
		}
		if fv.Kind() == reflect.Struct {
			problems = append(problems, Validate(fv.Interface()))
		}
	}
	return errors.Join(problems...)
}

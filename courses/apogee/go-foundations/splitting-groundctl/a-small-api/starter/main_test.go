package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"slices"
	"testing"

	"apogee/groundctl/telemetry"
)

func TestParse(t *testing.T) {
	line := "thistle 4182 2027-05-11T04:13:07Z batt_v=7.91 temp_c=-4.5 mode=NOMINAL"
	f, err := telemetry.Parse(line)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v, want nil", line, err)
	}
	if f.SatID != "thistle" {
		t.Errorf("Parse(%q).SatID = %q, want %q", line, f.SatID, "thistle")
	}
	if f.Seq != 4182 {
		t.Errorf("Parse(%q).Seq = %v, want 4182", line, f.Seq)
	}
	if math.Abs(f.BattV-7.91) > 0.01 {
		t.Errorf("Parse(%q).BattV = %.2f, want 7.91", line, f.BattV)
	}
}

func TestParseError(t *testing.T) {
	for _, line := range []string{
		"thistle 41x2 2027-05-11T04:13:07Z batt_v=7.91",
		"thistle 4183 2027-05-11T04:13:37Z batt_v=7.9l",
	} {
		if _, err := telemetry.Parse(line); err == nil {
			t.Errorf("Parse(%q) error = nil, want an error", line)
		}
	}
}

// TestExportsOnlyTheAPI checks that telemetry exports what other teams need, and nothing else.
func TestExportsOnlyTheAPI(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "telemetry/telemetry.go", nil, parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("reading telemetry/telemetry.go: %v", err)
	}
	var got []string
	for _, decl := range file.Decls {
		for _, name := range declared(decl) {
			if ast.IsExported(name) {
				got = append(got, name)
			}
		}
	}
	slices.Sort(got)
	if want := []string{"Frame", "Parse"}; !slices.Equal(got, want) {
		t.Errorf("telemetry exports %v, want %v", got, want)
	}
}

// declared returns the names a top-level declaration introduces, methods aside.
func declared(decl ast.Decl) []string {
	var names []string
	switch d := decl.(type) {
	case *ast.FuncDecl:
		if d.Recv == nil {
			names = append(names, d.Name.Name)
		}
	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				names = append(names, s.Name.Name)
			case *ast.ValueSpec:
				for _, n := range s.Names {
					names = append(names, n.Name)
				}
			}
		}
	}
	return names
}

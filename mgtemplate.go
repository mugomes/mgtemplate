// Copyright (C) 2026 Murilo Gomes Julio
// SPDX-License-Identifier: MIT

// Site: https://mugomes.github.io

package mgtemplate

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"
)

type MGTemplate struct {
	source     string
	context    map[string]any
	blocks     map[string]string
	blockAccum map[string]string
}

func ReadFile(path string) (*MGTemplate, error) {
	sHTML, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &MGTemplate{
		source:     string(sHTML),
		context:    map[string]any{},
		blocks:     map[string]string{},
		blockAccum: map[string]string{},
	}, nil
}

func (e *MGTemplate) IncludeFile(varname string, path string) error {
	sHTML, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	e.source = strings.ReplaceAll(e.source, "{{"+varname+"}}", string(sHTML))

	return nil
}

func (e *MGTemplate) Var(name string, value any) {
	e.context[name] = value
}

func (e *MGTemplate) VarExists(name string) bool {
	if strings.Contains(e.source, "{{"+name+"}}") {
		return true
	}

	return false
}

func (e *MGTemplate) Section(name string) {
	body, ok := e.blocks[name]
	if !ok {
		open := "[[" + name + "]]"
		close := "[[/" + name + "]]"

		a := strings.Index(e.source, open)
		b := strings.Index(e.source, close)
		if a == -1 || b == -1 || b < a {
			return
		}

		body = e.source[a+len(open) : b]
		e.blocks[name] = body

		token := "{{__" + name + "__}}"
		e.source = e.source[:a] + token + e.source[b+len(close):]
	}

	e.blockAccum[name] += e.interpolate(body)
}

func (e *MGTemplate) Render() string {
	out := e.source
	for k, v := range e.blockAccum {
		out = strings.ReplaceAll(out, "{{__"+k+"__}}", v)
	}
	return e.interpolate(out)
}

func (e *MGTemplate) interpolate(input string) string {
	var out strings.Builder
	for {
		start := strings.Index(input, "{{")
		if start == -1 {
			out.WriteString(input)
			break
		}

		out.WriteString(input[:start])
		end := strings.Index(input[start:], "}}")
		if end == -1 {
			out.WriteString(input)
			break
		}

		expr := input[start+2 : start+end]
		out.WriteString(e.eval(expr))

		input = input[start+end+2:]
	}
	return out.String()
}

func (e *MGTemplate) eval(expr string) string {
	parts := strings.Split(expr, "|")
	val := e.resolve(parts[0])

	for _, m := range parts[1:] {
		val = transform(val, m)
	}
	return val
}

func (e *MGTemplate) resolve(path string) string {
	segments := strings.Split(path, ".")
	current, ok := e.context[segments[0]]
	if !ok {
		return ""
	}

	v := reflect.ValueOf(current)
	for _, seg := range segments[1:] {
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}
		if v.Kind() != reflect.Struct {
			return ""
		}

		field := v.FieldByNameFunc(func(n string) bool {
			return normalize(n) == normalize(seg)
		})
		if !field.IsValid() {
			return ""
		}
		v = field
	}
	return fmt.Sprint(v.Interface())
}

func transform(v, op string) string {
	switch strings.ToLower(op) {
	case "upper":
		return strings.ToUpper(v)
	case "lower":
		return strings.ToLower(v)
	case "trim":
		return strings.TrimSpace(v)
	default:
		return v
	}
}

func normalize(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

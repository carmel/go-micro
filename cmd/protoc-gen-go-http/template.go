package main

import (
	"bytes"
	_ "embed"
	"strings"
	"text/template"
)

//go:embed httpTemplate.tpl
var httpTemplate string

type serviceDesc struct {
	MethodSets  map[string]*methodDesc
	ServiceType string
	ServiceName string
	Metadata    string
	Methods     []*methodDesc
}

type methodDesc struct {
	Name         string
	OriginalName string
	Request      string
	Reply        string
	Comment      string
	Path         string
	Method       string
	Body         string
	ResponseBody string
	Num          int
	HasVars      bool
	HasBody      bool
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}

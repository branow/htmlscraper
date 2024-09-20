package scrape

import (
	"fmt"
	"strings"
)

type ScrapeErr struct {
	Cause error
}

func (e ScrapeErr) Error() string {
	msg := e.Cause.Error()
	msgs := strings.Split(msg, "\n")
	for i, m := range msgs {
		msgs[i] = "scrape: " + m
	}
	return strings.Join(msgs, "\n")
}

type AttributeNotFoundErr struct {
	Attr string
}

func (e AttributeNotFoundErr) Error() string {
	return fmt.Sprintf("attribute \"%s\" not found", e.Attr)
}

type ExtractTagErr struct {
	ExtractTag string
}

func (e ExtractTagErr) Error() string {
	return fmt.Sprintf("invalid extract tag \"%s\"", e.ExtractTag)
}

type ScrapingErr struct {
	Selector string
	Cause    error
}

func (e ScrapingErr) Error() string {
	if e.Selector == "" {
		return e.Cause.Error()
	}
	msg := e.Cause.Error()
	msgs := strings.Split(msg, "\n")
	for i, m := range msgs {
		msgs[i] = e.Selector + " " + m
	}
	return strings.Join(msgs, "\n")
}

type NoNodesFoundErr struct{}

func (e NoNodesFoundErr) Error() string {
	return "no nodes found"
}

type KindErr struct {
	Var     any
	KindExp any
	KindAct any
}

func (e KindErr) Error() string {
	return fmt.Sprintf("%v must be a %v, but it is a %v", e.Var, e.KindExp, e.KindAct)
}

type NilErr struct {
	Var string
}

func (e NilErr) Error() string {
	return fmt.Sprintf("%s is nil", e.Var)
}

package service

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

type ParserService struct {
}

func NewParserService() *ParserService {
	return &ParserService{}
}

// ParseRef parses the raw text inside one ${...} expression
func (ps *ParserService) ParseRef(inner string) (TemplateRef, error) {
	p, err := newRefParser(inner)
	if err != nil {
		return TemplateRef{}, err
	}
	return p.parse(inner)
}

func (ps *ParserService) Parse(template string) (ParseResult, error) {
	res := ParseResult{Original: template}
	remaining := template
	cursor := 0

	for {
		start := strings.Index(remaining, "${")
		if start == -1 {
			res.Literals = append(res.Literals, remaining)
			break
		}

		end := strings.Index(remaining[start:], "}")
		if end == -1 {
			return res, fmt.Errorf("uncl")
		}
		end += start

		res.Literals = append(res.Literals, remaining[:start])

		inner := remaining[start+2 : end]
		ref, err := ps.ParseRef(inner)
		if err != nil {
			return res, fmt.Errorf("parse error in ${%s}: %w", inner, err)
		}

		absStart := cursor + start
		absEnd := cursor + end + 1
		res.Spans = append(res.Spans, Span{
			Start: absStart,
			End:   absEnd,
			Ref:   ref,
		})

		cursor += end + 1
		remaining = remaining[end+1:]
	}
	return res, nil
}

type StepKind string

const (
	StepKey   StepKind = "key"   // .fieldName
	StepIndex StepKind = "index" // [n]
)

type Step struct {
	Kind  StepKind
	Key   string // valid when Kind == StepKey
	Index int    // valid when Kind == StepIndex
}

func (s Step) String() string {
	if s.Kind == StepIndex {
		return fmt.Sprintf("[%d]", s.Index)
	}
	return "." + s.Key
}

type TemplateRef struct {
	Raw     string // exact text inside ${}
	Service string // first identifier
	Path    []Step // remaining accessor chain
}

func (r TemplateRef) String() string {
	var sb strings.Builder
	sb.WriteString("${")
	sb.WriteString(r.Service)
	for _, step := range r.Path {
		sb.WriteString(step.String())
	}
	sb.WriteString("}")
	return sb.String()
}

type Span struct {
	Start int // byte offset of '$'
	End   int // byte offset one past '}'
	Ref   TemplateRef
}

type ParseResult struct {
	Original string
	Literals []string
	Spans    []Span
}

// ---------------------------------------------------------------------------
// Lexer
// ---------------------------------------------------------------------------

type tokenType int

const (
	tokIdent tokenType = iota
	tokDot
	tokLBracket
	tokRBracket
	tokInt
	tokEOF
)

type token struct {
	kind  tokenType
	value string
}

type lexer struct {
	input []rune
	pos   int
}

func (l *lexer) peek() (rune, bool) {
	if l.pos >= len(l.input) {
		return 0, false
	}
	return l.input[l.pos], true
}

func (l *lexer) consume() rune {
	ch := l.input[l.pos]
	l.pos++
	return ch
}

func (l *lexer) next() (token, error) {
	ch, ok := l.peek()
	if !ok {
		return token{kind: tokEOF}, nil
	}
	switch {
	case ch == '.':
		l.consume()
		return token{tokDot, "."}, nil
	case ch == '[':
		l.consume()
		return token{tokLBracket, "["}, nil
	case ch == ']':
		l.consume()
		return token{tokRBracket, "]"}, nil
	case unicode.IsDigit(ch):
		var b strings.Builder
		for c, ok := l.peek(); ok && unicode.IsDigit(c); c, ok = l.peek() {
			b.WriteRune(l.consume())
		}
		return token{tokInt, b.String()}, nil
	default:
		var b strings.Builder
		for c, ok := l.peek(); ok && c != '.' && c != '[' && c != ']'; c, ok = l.peek() {
			b.WriteRune(l.consume())
		}
		if b.Len() == 0 {
			return token{}, fmt.Errorf("unexpected character: %q at position %d", ch, l.pos)
		}
		return token{tokIdent, b.String()}, nil
	}
}

// ---------------------------------------------------------------------------
// Parser
// ---------------------------------------------------------------------------
//
// Grammar (EBNF):
//   ref      = ident { accessor } EOF
//   accessor = '.' ident
//            | '[' int ']'

type refParser struct {
	lex *lexer
	cur token
}

func newRefParser(inner string) (*refParser, error) {
	p := &refParser{lex: &lexer{input: []rune(inner)}}
	return p, p.advance()
}

func (p *refParser) advance() error {
	tok, err := p.lex.next()
	if err != nil {
		return err
	}
	p.cur = tok
	return nil
}

func (p *refParser) expect(k tokenType, desc string) error {
	if p.cur.kind != k {
		return fmt.Errorf("expected %s, got %q", desc, p.cur.value)
	}
	return nil
}

func (p *refParser) parse(raw string) (TemplateRef, error) {
	ref := TemplateRef{Raw: raw}

	if err := p.expect(tokIdent, "service name"); err != nil {
		return ref, err
	}
	ref.Service = p.cur.value
	if err := p.advance(); err != nil {
		return ref, err
	}

	for p.cur.kind != tokEOF {
		switch p.cur.kind {
		case tokDot:
			if err := p.advance(); err != nil {
				return ref, err
			}
			if err := p.expect(tokIdent, "field name after '.'"); err != nil {
				return ref, err
			}
			ref.Path = append(ref.Path, Step{Kind: StepKey, Key: p.cur.value})
			if err := p.advance(); err != nil {
				return ref, err
			}

		case tokLBracket:
			if err := p.advance(); err != nil {
				return ref, err
			}
			if err := p.expect(tokInt, "integer index"); err != nil {
				return ref, err
			}
			idx, err := strconv.Atoi(p.cur.value)
			if err != nil {
				return ref, fmt.Errorf("invalid index %q: %w", p.cur.value, err)
			}
			ref.Path = append(ref.Path, Step{Kind: StepIndex, Index: idx})
			if err := p.advance(); err != nil {
				return ref, err
			}
			if err := p.expect(tokRBracket, "']'"); err != nil {
				return ref, err
			}
			if err := p.advance(); err != nil {
				return ref, err
			}

		default:
			return ref, fmt.Errorf("unexpected token: %q", p.cur.value)
		}
	}
	return ref, nil
}

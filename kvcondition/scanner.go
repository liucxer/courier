package kvcondition

import (
	"bytes"
	"errors"
	"text/scanner"
)

var (
	SyntaxErrorMissingRightBracket = errors.New("kv condition missing right bracket")
	SyntaxErrorInvalidCondition    = errors.New("kv condition invalid condition")
)

var keywords = map[rune]bool{
	'!': true,
	'=': true,
	'(': true,
	')': true,
	'*': true,
	'^': true,
	'$': true,
	'&': true,
	'|': true,
}

func newNodeScanner(b []byte) *nodeScanner {
	s := &scanner.Scanner{}
	s.Init(bytes.NewReader(b))

	return &nodeScanner{
		Scanner: s,
	}
}

type nodeScanner struct {
	*scanner.Scanner
}

func (s *nodeScanner) ScanNode() (Node, error) {
	node := Node(nil)
	rs := &ruleScanner{}

	completeNode := func(nextNode Node) {
		if c, ok := node.(*Condition); ok {
			c.Right = nextNode
		} else {
			node = nextNode
		}
	}

	completeRule := func() error {
		r, err := rs.ToRule()
		if err != nil {
			return err
		}
		if r != nil {
			completeNode(r)
			rs = &ruleScanner{}
		}
		return nil
	}

	for {
		tok := s.Peek()

		if tok == ')' {
			if err := completeRule(); err != nil {
				return nil, err
			}
			break
		}

		tok = s.Next()

		// quote skip
		if tok == '"' {
			for {
				tok = s.Next()

				if tok == scanner.EOF {
					break
				}

				if tok == '"' {
					tok = s.Next()
					break
				}

				if tok == '\\' {
					tok = s.Next()
				}

				rs.WriteRune(tok)
			}
		}

		if tok == scanner.EOF {
			if err := completeRule(); err != nil {
				return nil, err
			}
			break
		}

		if keywords[tok] {
			switch tok {
			case '(', '&', '|':
				if err := completeRule(); err != nil {
					return nil, err
				}

				switch tok {
				case '(':
					nextNode, err := s.ScanNode()
					if err != nil {
						return nil, err
					}
					completeNode(nextNode)

					if s.Peek() != ')' {
						return nil, SyntaxErrorMissingRightBracket
					}

					tok = s.Next()
				case '&':
					node = And(node, nil)
				case '|':
					node = Or(node, nil)
				}
			default:
				// collect operator
				rs.WriteOperator(tok)
			}
			continue
		}

		// collect key or value
		rs.WriteRune(tok)
	}

	if expr, ok := node.(*Condition); ok {
		if expr.Right == nil {
			return nil, SyntaxErrorInvalidCondition
		}
	}

	return node, nil
}

type ruleScanner struct {
	key      bytes.Buffer
	value    bytes.Buffer
	operator bytes.Buffer
}

func (s *ruleScanner) WriteRune(r rune) {
	if s.operator.Len() != 0 {
		s.value.WriteRune(r)
	} else {
		s.key.WriteRune(r)
	}
}

func (s *ruleScanner) WriteOperator(r rune) {
	s.operator.WriteRune(r)
}

func (s *ruleScanner) ToRule() (*Rule, error) {
	k := bytes.TrimSpace(s.key.Bytes())

	if len(k) == 0 {
		return nil, nil
	}

	operator, err := ParseOperator(s.operator.String())
	if err != nil {
		return nil, err
	}

	r := &Rule{
		Operator: operator,
		Key:      k,
		Value:    bytes.TrimSpace(s.value.Bytes()),
	}

	return r, nil
}

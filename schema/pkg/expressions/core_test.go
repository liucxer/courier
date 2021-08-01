package expressions

import (
	"context"
	"fmt"
	"testing"

	"github.com/onsi/gomega"
)

type fixture struct {
	args []interface{}
	ret  interface{}
}

type args = []interface{}

var (
	cases = []struct {
		summary  string
		expr     Expression
		fixtures []fixture
	}{
		{
			summary: "simple expression",
			expr:    Expression{"eq", 1.0},
			fixtures: []fixture{
				{args: args{1}, ret: true},
				{args: args{2}, ret: false},
			},
		},
		{
			summary: "match",
			expr:    Expression{"match", `[a-z]+`},
			fixtures: []fixture{
				{args: args{"abc"}, ret: true},
				{args: args{"ABC"}, ret: false},
			},
		},
		{
			summary: "in",
			expr:    Expression{"in", Expression{"literal", []interface{}{"a", "b", "c"}}},
			fixtures: []fixture{
				{args: args{"a"}, ret: true},
				{args: args{"d"}, ret: false},
			},
		},
		{
			summary: "in string",
			expr:    Expression{"in", "abc"},
			fixtures: []fixture{
				{args: args{"bc"}, ret: true},
				{args: args{"ab"}, ret: true},
				{args: args{"ac"}, ret: false},
			},
		},
		{
			summary: "with maths",
			expr:    Expression{"eq", 0, Expression{"mod", 2}},
			fixtures: []fixture{
				{args: args{2}, ret: true},
				{args: args{3}, ret: false},
			},
		},
		{
			summary: "all",
			expr:    Expression{"all", Expression{"gt", 1.0}, Expression{"lte", 3.0}},
			fixtures: []fixture{
				{args: args{2}, ret: true},
				{args: args{2.11111}, ret: true},
				{args: args{2.99999}, ret: true},
				{args: args{3}, ret: true},
				{args: args{4}, ret: false},
				{args: args{-1}, ret: false},
				{args: args{1}, ret: false},
			},
		},
		{
			summary: "any",
			expr:    Expression{"any", Expression{"eq", 5.0}, Expression{"lte", 3}},
			fixtures: []fixture{
				{args: args{5}, ret: true},
				{args: args{3}, ret: true},
				{args: args{1}, ret: true},
				{args: args{7}, ret: false},
			},
		},
		{
			summary: "case",
			expr: Expression{"case",
				// if
				Expression{"lt", 2}, "1",
				// else if
				Expression{"lt", 10}, "2",
				// else
				"3",
			},
			fixtures: []fixture{
				{args: args{1}, ret: "1"},
				{args: args{5}, ret: "2"},
				{args: args{11}, ret: "3"},
			},
		},
	}
)

func TestFactory(t *testing.T) {
	for i := range cases {
		c := cases[i]

		for i := range c.fixtures {
			f := c.fixtures[i]

			t.Run(fmt.Sprintf("|%s|=|%s|", StringifyExpression(append(Expression{StringifyExpression(c.expr)}, f.args...)), StringifyValue(f.ret)), func(t *testing.T) {
				exec, err := DefaultFactory.From(c.expr)
				gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

				ret, err := exec(context.Background(), f.args...)
				gomega.NewWithT(t).Expect(err).To(gomega.BeNil())
				gomega.NewWithT(t).Expect(ret).To(gomega.Equal(f.ret))
			})
		}
	}
}

func BenchmarkFactory(b *testing.B) {
	for i := range cases {
		c := cases[i]

		f := c.fixtures[0]

		exec, _ := DefaultFactory.From(c.expr)

		b.Run(c.summary, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = exec(context.Background(), f.args...)
			}
		})
	}
}

package main

import (
	"github.com/liucxer/courier/codegen"
)

func main() {
	basicTypes := []codegen.BuiltInType{
		"bool",

		"int",
		"int8",
		"int16",
		"int32",
		"int64",

		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"uintptr",

		"float32",
		"float64",
		"complex64",
		"complex128",

		"string",
		"byte",
		"rune",
	}

	{
		file := codegen.NewFile("ptr", "ptr.go")

		for _, basicType := range basicTypes {
			file.WriteBlock(
				codegen.Func(codegen.Var(basicType, "v")).
					Return(codegen.Var(codegen.Star(basicType))).
					Named(codegen.UpperCamelCase(string(basicType))).
					Do(codegen.Return(codegen.Unary(codegen.Id("v")))),
			)
			file.WriteRune('\n')
		}

		file.WriteFile()
	}

	{
		file := codegen.NewFile("ptr", "ptr_test.go")

		for _, basicType := range basicTypes {
			name := codegen.UpperCamelCase(string(basicType))

			funcType := codegen.Func(codegen.Var(codegen.Star(codegen.Type(file.Use("testing", "T"))), "t")).
				Named("Test" + name)

			asEqual := func(values ...codegen.Snippet) {
				blocks := make([]codegen.Snippet, len(values))

				for i := range blocks {
					blocks[i] = codegen.Expr(
						"?(?).Expect(?).To(?(?))",
						codegen.Id(file.Use("github.com/onsi/gomega", "NewWithT")),
						codegen.Id("t"),
						codegen.CallWith(basicType, values[i]),
						codegen.Id(file.Use("github.com/onsi/gomega", "Equal")),
						codegen.Expr("*?", codegen.Call(name, values[i])),
					)
				}

				file.WriteBlock(funcType.Do(blocks...))
			}

			switch basicType {
			case "string":
				asEqual(codegen.Val("string"))
			case "rune":
				asEqual(codegen.Val('r'))
			case "byte":
				asEqual(codegen.Expr("?[0]", codegen.Val([]byte("bytes"))))
			case "bool":
				asEqual(codegen.Val(true), codegen.Val(false))
			default:
				asEqual(codegen.Val(1))
			}

			file.WriteRune('\n')
		}

		file.WriteFile()
	}
}

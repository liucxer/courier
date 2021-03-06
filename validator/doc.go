/*
Each Validator have at least two process methods, one for 'Parsing' and one for 'Validating'.

In 'Parsing' stage, we call validator as 'ValidatorCreator', and it will be register to some 'ValidatorMgr' for caching.


Parsing

There are common 'Rule DSL' below:

	// simple
	@name

	// with parameters
	@name<param1>
	@name<param1,param2>

	// with ranges
	@name[from, to)
	@name[length]

	// with values
	@name{VALUE1,VALUE2,VALUE3}
	@name{%v}

	// with regexp
	@name/\d+/

	// optional and default value
	@name?
	@name = value
	@name = 'some string value'

	// composes
	@map<@string[1,10],@string{A,B,C}>
	@map<@string[1,10],@string/\d+/>[0,10]

Then the parsed rule will be transform to special validators:

@string: https://godoc.org/github.com/liucxer/courier/validator#StringValidator

@uint: https://godoc.org/github.com/liucxer/courier/validator#UintValidator

@int: https://godoc.org/github.com/liucxer/courier/validator#IntValidator

@float: https://godoc.org/github.com/liucxer/courier/validator#FloatValidator

@struct: https://godoc.org/github.com/liucxer/courier/validator#StructValidator

@map: https://godoc.org/github.com/liucxer/courier/validator#MapValidator

@slice: https://godoc.org/github.com/liucxer/courier/validator#SliceValidator


Validating

We can create validator by 'Rule DSL', and also can configure them by validator struct field as conditions.
Then call the method `Validate(v interface{}) error` to do value validations.
*/
package validator

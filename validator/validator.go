package validator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/liucxer/courier/reflectx/typesutil"
	"github.com/liucxer/courier/validator/rules"
)

func MustParseRuleStringWithType(ruleStr string, typ typesutil.Type) *Rule {
	r, err := ParseRuleWithType([]byte(ruleStr), typ)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseRuleWithType(ruleBytes []byte, typ typesutil.Type) (*Rule, error) {
	r := &rules.Rule{}

	if len(ruleBytes) != 0 {
		parsedRule, err := rules.ParseRule(ruleBytes)
		if err != nil {
			return nil, err
		}
		r = parsedRule
	}

	return &Rule{
		Type: typ,
		Rule: r,
	}, nil
}

type Rule struct {
	*rules.Rule

	ErrMsg []byte
	Type   typesutil.Type
}

func (r *Rule) String() string {
	return typesutil.FullTypeName(r.Type) + string(r.Rule.Bytes())
}

type RuleProcessor func(rule *Rule)

// mgr for compiling validator
type ValidatorMgr interface {
	// compile rule string to validator
	Compile(context.Context, []byte, typesutil.Type, RuleProcessor) (Validator, error)
}

var ValidatorMgrDefault = NewValidatorFactory()

type contextKeyValidatorMgr int

func ContextWithValidatorMgr(c context.Context, validatorMgr ValidatorMgr) context.Context {
	return context.WithValue(c, contextKeyValidatorMgr(1), validatorMgr)
}

func ValidatorMgrFromContext(c context.Context) ValidatorMgr {
	return c.Value(contextKeyValidatorMgr(1)).(ValidatorMgr)
}

type ValidatorCreator interface {
	// name and aliases of validator
	// we will register validator to validator set by these names
	Names() []string
	// create new instance
	New(context.Context, *Rule) (Validator, error)
}

type Validator interface {
	// validate value
	Validate(v interface{}) error
	// stringify validator rule
	String() string
}

func NewValidatorFactory() *ValidatorFactory {
	return &ValidatorFactory{
		validatorSet: map[string]ValidatorCreator{},
	}
}

type ValidatorFactory struct {
	validatorSet map[string]ValidatorCreator
}

func (f *ValidatorFactory) Register(validators ...ValidatorCreator) {
	for i := range validators {
		validator := validators[i]
		for _, name := range validator.Names() {
			f.validatorSet[name] = validator
		}
	}
}

func (f *ValidatorFactory) MustCompile(ctx context.Context, rule []byte, typ typesutil.Type, ruleProcessor RuleProcessor) Validator {
	v, err := f.Compile(ctx, rule, typ, ruleProcessor)
	if err != nil {
		panic(err)
	}
	return v
}

func (f *ValidatorFactory) Compile(ctx context.Context, ruleBytes []byte, typ typesutil.Type, ruleProcessor RuleProcessor) (Validator, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if len(ruleBytes) == 0 {
		if _, ok := typesutil.EncodingTextMarshalerTypeReplacer(typ); !ok {
			switch typesutil.Deref(typ).Kind() {
			case reflect.Struct:
				ruleBytes = []byte("@struct")
			case reflect.Slice:
				ruleBytes = []byte("@slice")
			case reflect.Map:
				ruleBytes = []byte("@map")
			}
		}
	}

	rule, err := ParseRuleWithType(ruleBytes, typ)
	if err != nil {
		return nil, err
	}

	if ruleProcessor != nil {
		ruleProcessor(rule)
	}

	validatorCreator, ok := f.validatorSet[rule.Name]
	if len(ruleBytes) != 0 && !ok {
		return nil, fmt.Errorf("%s not match any validator", rule.Name)
	}

	return NewValidatorLoader(validatorCreator).New(ContextWithValidatorMgr(ctx, f), rule)
}

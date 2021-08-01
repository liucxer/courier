package kvcondition_test

import (
	"encoding/json"
	"testing"

	"github.com/liucxer/courier/kvcondition"
	. "github.com/onsi/gomega"
)

func BenchmarkParseKVCondition(b *testing.B) {
	rule := []byte(`env & tag = ONLINE & tag = "some label" & ip & ( ip != 1.1.1.1 | ip ^= "8.8" | tag *= "test&" ) | ip $= 4.4`)

	for i := 0; i < b.N; i++ {
		kvcondition.ParseKVCondition(rule)
	}
}

func TestKVCondition_IsZero(t *testing.T) {
	ql, _ := kvcondition.ParseKVCondition([]byte(``))

	NewWithT(t).Expect(ql.IsZero()).To(BeTrue())
}

func TestKVCondition(t *testing.T) {
	ql, _ := kvcondition.ParseKVCondition([]byte(`ip != "1.1.1.1"`))

	type Data struct {
		QL kvcondition.KVCondition `json:"ql"`
	}

	data, err := json.Marshal(&Data{
		QL: *ql,
	})
	NewWithT(t).Expect(err).To(BeNil())

	d := Data{}
	er := json.Unmarshal(data, &d)

	NewWithT(t).Expect(er).To(BeNil())
	NewWithT(t).Expect(ql.String()).To(Equal(d.QL.String()))
}

func TestParseKVCondition(t *testing.T) {
	rule := []byte(`env & tag = ONLINE & tag = "some label" & ip & ( ip != 1.1.1.1 | ip ^= "8.8" | tag *= "test\&" ) | ip $= 4.4`)

	kvCondition := &kvcondition.KVCondition{}
	err := kvCondition.UnmarshalText(rule)
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(kvCondition.Node.String()).To(Equal(`( ( ( ( ( env & tag = "ONLINE" ) & tag = "some label" ) & ip ) & ( ( ip != "1.1.1.1" | ip ^= "8.8" ) | tag *= "test&" ) ) | ip $= "4.4" )`))

	kvc, err := kvcondition.ParseKVCondition([]byte(kvCondition.Node.String()))
	NewWithT(t).Expect(err).To(BeNil())
	NewWithT(t).Expect(kvc.String()).To(Equal(kvCondition.String()))

	rules := make([]*kvcondition.Rule, 0)

	kvc.Range(func(label *kvcondition.Rule) {
		rules = append(rules, label)
	})

	NewWithT(t).Expect(rules).To(Equal([]*kvcondition.Rule{
		kvcondition.OperatorExists.Of("env", ""),
		kvcondition.OperatorEqual.Of("tag", "ONLINE"),
		kvcondition.OperatorEqual.Of("tag", "some label"),
		kvcondition.OperatorExists.Of("ip", ""),
		kvcondition.OperatorNotEqual.Of("ip", "1.1.1.1"),
		kvcondition.OperatorStartsWith.Of("ip", "8.8"),
		kvcondition.OperatorContains.Of("tag", "test&"),
		kvcondition.OperatorEndsWith.Of("ip", "4.4"),
	}))
}

func TestParseKVConditionFailed(t *testing.T) {
	rule := "tag = ONLINE & tag = \"some label\" & ( ip = 1.1.1.1 | ip = \"8.8.8.8\" | tag = test & ip = 4.4.4.4"

	_, err := kvcondition.ParseKVCondition([]byte(rule))
	NewWithT(t).Expect(err).NotTo(BeNil())
}

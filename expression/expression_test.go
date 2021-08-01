package expression

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	ex = Ex("allOf",
		Ex("has", "key"),
		Ex("not", Ex("toBoolean", Ex("get", "should"))),
		Ex("match", "[0-9]+", Ex("get", "key")),
		Ex("eq", Ex("toNumber", Ex("get", "key")), 10),
		Ex("in", Ex("toString", Ex("get", "key")), "10", "11"),
		Ex("lt", Ex("toNumber", Ex("get", "key")), 20),
		Ex("neq", Ex("toNumber", Ex("get", "key")), 20),
		Ex("lte", Ex("toNumber", Ex("get", "key")), 10),
		Ex("gt", Ex("toNumber", Ex("get", "key")), 1),
		Ex("gte", Ex("toNumber", Ex("get", "key")), 10),
	)
	expr, _ = json.MarshalIndent(ex, "", "  ")
)

func TestExpression(t *testing.T) {
	t.Run("#ExDo", func(t *testing.T) {
		args := make([]interface{}, 0)
		err := json.Unmarshal(expr, &args)
		require.NoError(t, err)

		exDo, err := New(args)
		require.NoError(t, err)

		v, err := exDo(map[string]interface{}{
			"key":    "10",
			"should": "false",
		})
		require.NoError(t, err)
		require.True(t, v == true)
	})

}

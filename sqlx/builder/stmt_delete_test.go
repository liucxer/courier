package builder_test

import (
	"testing"

	. "github.com/liucxer/courier/sqlx/builder"
	. "github.com/liucxer/courier/sqlx/builder/buidertestingutils"
	"github.com/onsi/gomega"
)

func TestStmtDelete(t *testing.T) {
	table := T("T")

	t.Run("delete", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Delete().From(table,
				Where(Col("F_a").Eq(1)),
				Comment("Comment"),
			),
		).To(BeExpr(`
DELETE FROM T
WHERE f_a = ?
/* Comment */
`, 1))
	})
}

package builder_test

import (
	"testing"

	. "github.com/liucxer/courier/sqlx/builder"
	. "github.com/liucxer/courier/sqlx/builder/buidertestingutils"
	"github.com/onsi/gomega"
)

func TestStmtUpdate(t *testing.T) {
	table := T("T")

	t.Run("update", func(t *testing.T) {
		gomega.NewWithT(t).Expect(
			Update(table).
				Set(
					Col("F_a").ValueBy(1),
					Col("F_b").ValueBy(2),
				).
				Where(
					Col("F_a").Eq(1),
					Comment("Comment"),
				),
		).To(BeExpr(`
UPDATE T SET f_a = ?, f_b = ?
WHERE f_a = ?
/* Comment */`, 1, 2, 1))
	})
}

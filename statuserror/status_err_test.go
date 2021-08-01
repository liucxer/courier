package statuserror_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/pkg/errors"

	"github.com/liucxer/courier/statuserror"
	examples "github.com/liucxer/courier/statuserror/__examples__"
	. "github.com/onsi/gomega"
)

func ExampleStatusErr() {
	fmt.Println(examples.Unauthorized)
	fmt.Println(statuserror.FromErr(nil))
	fmt.Println(statuserror.FromErr(fmt.Errorf("unknown")))
	// Output:
	//[]@StatusErr[Unauthorized][401999001][Unauthorized]!
	//<nil>
	//[]@StatusErr[UnknownError][500000000][unknown error] unknown
	//
}

func TestStatusErr(t *testing.T) {
	unknownErr := statuserror.Wrap(errors.New(""), http.StatusInternalServerError, "UnknownError")

	t.Logf("%+v", unknownErr)

	summary := unknownErr.Summary()

	NewWithT(t).Expect(summary).To(Equal("@StatusErr[UnknownError][500000000][UnknownError]"))

	statusErr, err := statuserror.ParseStatusErrSummary(summary)
	NewWithT(t).Expect(err).To(BeNil())

	NewWithT(t).Expect(statusErr.Summary()).To(Equal(unknownErr.Summary()))

	NewWithT(t).Expect(examples.Unauthorized.StatusErr().Summary()).To(Equal("@StatusErr[Unauthorized][401999001][Unauthorized]!"))
	NewWithT(t).Expect(examples.InternalServerError.StatusErr().Summary()).To(Equal("@StatusErr[InternalServerError][500999001][InternalServerError]"))
	NewWithT(t).Expect(examples.Unauthorized.StatusCode()).To(Equal(401))
	NewWithT(t).Expect(examples.Unauthorized.StatusErr().StatusCode()).To(Equal(401))

	NewWithT(t).Expect(errors.Is(examples.Unauthorized, examples.Unauthorized)).To(BeTrue())
	NewWithT(t).Expect(errors.Is(examples.Unauthorized.StatusErr(), examples.Unauthorized)).To(BeTrue())
	NewWithT(t).Expect(errors.Is(examples.Unauthorized.StatusErr(), examples.Unauthorized.StatusErr())).To(BeTrue())
}

func TestStatusErrBuilders(t *testing.T) {
	t.Log(examples.Unauthorized.StatusErr().WithMsg("msg overwrite"))
	t.Log(examples.Unauthorized.StatusErr().WithDesc("desc overwrite"))
	t.Log(examples.Unauthorized.StatusErr().DisableErrTalk().EnableErrTalk())
	t.Log(examples.Unauthorized.StatusErr().WithID("111"))
	t.Log(examples.Unauthorized.StatusErr().AppendSource("service-abc"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorField("header", "Authorization", "missing"))
	t.Log(examples.Unauthorized.StatusErr().AppendErrorFields(
		statuserror.NewErrorField("query", "key", "missing"),
		statuserror.NewErrorField("header", "Authorization", "missing"),
	))
}

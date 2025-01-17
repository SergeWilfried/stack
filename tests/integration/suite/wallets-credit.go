package suite

import (
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	. "github.com/formancehq/stack/tests/integration/internal"
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"math/big"
)

var _ = WithModules([]*Module{modules.Auth, modules.Ledger, modules.Wallets}, func() {

	When("creating a wallet", func() {
		var (
			response *operations.CreateWalletResponse
			err      error
		)
		BeforeEach(func() {
			response, err = Client().Wallets.CreateWallet(
				TestContext(),
				&shared.CreateWalletRequest{
					Name:     uuid.NewString(),
					Metadata: map[string]string{},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))
		})
		Then("crediting it", func() {
			BeforeEach(func() {
				_, err := Client().Wallets.CreditWallet(TestContext(), operations.CreditWalletRequest{
					CreditWalletRequest: &shared.CreditWalletRequest{
						Amount: shared.Monetary{
							Amount: big.NewInt(1000),
							Asset:  "USD/2",
						},
						Sources:  []shared.Subject{},
						Metadata: map[string]string{},
					},
					ID: response.CreateWalletResponse.Data.ID,
				})
				Expect(err).To(Succeed())
			})
			It("should be ok", func() {})
		})
	})
})

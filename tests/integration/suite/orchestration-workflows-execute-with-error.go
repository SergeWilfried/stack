package suite

import (
	"github.com/formancehq/stack/tests/integration/internal/modules"
	"math/big"
	"net/http"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	. "github.com/formancehq/stack/tests/integration/internal"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pborman/uuid"
)

var _ = WithModules([]*Module{modules.Auth, modules.Orchestration, modules.Ledger}, func() {
	BeforeEach(func() {
		createLedgerResponse, err := Client().Ledger.V2CreateLedger(TestContext(), operations.V2CreateLedgerRequest{
			Ledger: "default",
		})
		Expect(err).To(BeNil())
		Expect(createLedgerResponse.StatusCode).To(Equal(http.StatusNoContent))
	})
	When("creating a new workflow which will fail with insufficient fund error", func() {
		var (
			createWorkflowResponse *shared.V2CreateWorkflowResponse
		)
		BeforeEach(func() {
			response, err := Client().Orchestration.V2CreateWorkflow(
				TestContext(),
				&shared.V2CreateWorkflowRequest{
					Name: ptr(uuid.New()),
					Stages: []map[string]interface{}{
						{
							"send": map[string]any{
								"source": map[string]any{
									"account": map[string]any{
										"id":     "empty:account",
										"ledger": "default",
									},
								},
								"destination": map[string]any{
									"account": map[string]any{
										"id":     "bank",
										"ledger": "default",
									},
								},
								"amount": map[string]any{
									"amount": 100,
									"asset":  "EUR/2",
								},
							},
						},
					},
				},
			)
			Expect(err).ToNot(HaveOccurred())
			Expect(response.StatusCode).To(Equal(201))

			createWorkflowResponse = response.V2CreateWorkflowResponse
		})
		Then("executing it", func() {
			var runWorkflowResponse *shared.V2RunWorkflowResponse
			BeforeEach(func() {
				response, err := Client().Orchestration.V2RunWorkflow(
					TestContext(),
					operations.V2RunWorkflowRequest{
						RequestBody: map[string]string{},
						WorkflowID:  createWorkflowResponse.Data.ID,
					},
				)
				Expect(err).ToNot(HaveOccurred())
				Expect(response.StatusCode).To(Equal(201))

				runWorkflowResponse = response.V2RunWorkflowResponse
			})
			Then("waiting for first stage retried at least once", func() {
				var getWorkflowInstanceHistoryStageResponse *shared.V2GetWorkflowInstanceHistoryStageResponse
				BeforeEach(func() {
					Eventually(func(g Gomega) int64 {

						response, err := Client().Orchestration.V2GetInstanceStageHistory(
							TestContext(),
							operations.V2GetInstanceStageHistoryRequest{
								InstanceID: runWorkflowResponse.Data.ID,
								Number:     0,
							},
						)
						if err != nil {
							return 0
						}
						if response.StatusCode != 200 {
							return 0
						}

						getWorkflowInstanceHistoryStageResponse = response.V2GetWorkflowInstanceHistoryStageResponse
						g.Expect(getWorkflowInstanceHistoryStageResponse.Data).To(HaveLen(1))
						return getWorkflowInstanceHistoryStageResponse.Data[0].Attempt
					}).Should(BeNumerically(">", 2))
				})
				It("should be retried with insufficient fund error ", func() {
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt).NotTo(BeZero())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution).NotTo(BeNil())
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0].Attempt).To(BeNumerically(">", 2))
					Expect(getWorkflowInstanceHistoryStageResponse.Data[0]).To(Equal(shared.V2WorkflowInstanceHistoryStage{
						Name: "CreateTransaction",
						Input: shared.V2WorkflowInstanceHistoryStageInput{
							CreateTransaction: &shared.V2ActivityCreateTransaction{
								Ledger: ptr("default"),
								Data: &shared.V2PostTransaction{
									Postings: []shared.V2Posting{{
										Amount:      big.NewInt(100),
										Asset:       "EUR/2",
										Destination: "bank",
										Source:      "empty:account",
									}},
									Metadata: metadata.Metadata{},
								},
							},
						},
						LastFailure:   ptr("running numscript: script execution failed: no more fund to withdraw"),
						Attempt:       getWorkflowInstanceHistoryStageResponse.Data[0].Attempt,
						NextExecution: getWorkflowInstanceHistoryStageResponse.Data[0].NextExecution,
						StartedAt:     getWorkflowInstanceHistoryStageResponse.Data[0].StartedAt,
					}))
				})
			})
		})
	})
})

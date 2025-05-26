package jobs_test

import (
	"context"
	"testing"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/custom_hostnames"
	"github.com/riverqueue/river"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/theopenlane/core/pkg/openlaneclient"

	cfmocks "github.com/theopenlane/riverboat/internal/cloudflare/mocks"
	olmocks "github.com/theopenlane/riverboat/internal/olclient/mocks"
	"github.com/theopenlane/riverboat/pkg/jobs"
)

func (suite *TestSuite) TestDeleteCustomDomainWorker() {
	t := suite.T()
	testCases := []struct {
		name                string
		callCustomHostnames bool
		input               jobs.DeleteCustomDomainArgs
	}{
		{
			name: "happy path, delete custom domain",
			input: jobs.DeleteCustomDomainArgs{
				CustomDomainID: "customdomainid123",
			},
		},
		{
			name:                "happy path, delete cloudflare custom hostname",
			callCustomHostnames: true,
			input: jobs.DeleteCustomDomainArgs{
				CloudflareCustomHostnameID: "cloudflarecustomhostnameid123",
				CloudflareZoneID:           "cloudflarezoneid123",
			},
		},
		{
			name: "happy path, delete dns verification",
			input: jobs.DeleteCustomDomainArgs{
				DNSVerificationID: "dnsverificationid123",
			},
		},
		{
			name:                "happy path, delete all",
			callCustomHostnames: true,
			input: jobs.DeleteCustomDomainArgs{
				DNSVerificationID:          "dnsverificationid123",
				CloudflareCustomHostnameID: "cloudflarecustomhostnameid123",
				CloudflareZoneID:           "cloudflarezoneid123",
				CustomDomainID:             "customdomainid123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfMock := cfmocks.NewMockClient(t)
			cfHostnamesMock := cfmocks.NewMockCustomHostnamesService(t)

			if tc.callCustomHostnames {
				cfMock.EXPECT().CustomHostnames().Return(cfHostnamesMock)
			}

			olMock := olmocks.NewMockOpenlaneGraphClient(t)

			if tc.input.CustomDomainID != "" {
				olMock.EXPECT().DeleteCustomDomain(mock.Anything, tc.input.CustomDomainID).Return(&openlaneclient.DeleteCustomDomain{
					DeleteCustomDomain: openlaneclient.DeleteCustomDomain_DeleteCustomDomain{},
				}, nil)
			}

			if tc.input.CloudflareCustomHostnameID != "" {
				cfHostnamesMock.EXPECT().Delete(mock.Anything, tc.input.CloudflareCustomHostnameID, custom_hostnames.CustomHostnameDeleteParams{
					ZoneID: cloudflare.F(tc.input.CloudflareZoneID),
				}).Return(&custom_hostnames.CustomHostnameDeleteResponse{}, nil)
			}

			if tc.input.DNSVerificationID != "" {
				olMock.EXPECT().DeleteDNSVerification(mock.Anything, tc.input.DNSVerificationID).Return(&openlaneclient.DeleteDNSVerification{
					DeleteDNSVerification: openlaneclient.DeleteDNSVerification_DeleteDNSVerification{},
				}, nil)
			}

			worker := &jobs.DeleteCustomDomainWorker{
				Config: jobs.CustomDomainConfig{
					CloudflareAPIKey: "test",
				},
			}
			worker.WithCloudflareClient(cfMock)
			worker.WithOpenlaneClient(olMock)

			ctx := context.Background()

			err := worker.Work(ctx, &river.Job[jobs.DeleteCustomDomainArgs]{Args: tc.input})

			require.NoError(t, err)
		})
	}
}

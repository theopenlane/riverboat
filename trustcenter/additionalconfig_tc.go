//go:build trustcenter

package trustcenter

import (
	"github.com/theopenlane/corejobs"

	"github.com/theopenlane/riverboat/pkg/jobs"
)

// Workers holds the configuration for additional trust center specific workers
type Workers struct {
	// OpenlaneConfig configuration for openlane jobs, this is shared across multiple workers
	// if a worker needs specific configuration, it can be set in the worker's config
	OpenlaneConfig jobs.OpenlaneConfig `koanf:"openlaneconfig" json:"openlaneconfig"`

	// CreateCustomDomainWorker configuration for creating custom domains
	CreateCustomDomainWorker corejobs.CreateCustomDomainWorker `koanf:"createcustomdomainworker" json:"createcustomdomainworker"`

	// ValidateCustomDomainWorker configuration for validating custom domains
	ValidateCustomDomainWorker corejobs.ValidateCustomDomainWorker `koanf:"validatecustomdomainworker" json:"validatecustomdomainworker"`

	// DeleteCustomDomainWorker configuration for deleting custom domains
	DeleteCustomDomainWorker corejobs.DeleteCustomDomainWorker `koanf:"deletecustomdomainworker" json:"deletecustomdomainworker"`

	// WatermarkDocWorker configuration for watermarking documents
	WatermarkDocWorker corejobs.WatermarkDocWorker `koanf:"watermarkdocworker" json:"watermarkdocworker"`

	// CreatePirschDomainWorker configuration for creating Pirsch domains
	CreatePirschDomainWorker corejobs.CreatePirschDomainWorker `koanf:"createpirschdomainworker" json:"createpirschdomainworker"`

	// DeletePirschDomainWorker configuration for deleting Pirsch domains
	DeletePirschDomainWorker corejobs.DeletePirschDomainWorker `koanf:"deletepirschdomainworker" json:"deletepirschdomainworker"`

	// UpdatePirschDomainWorker configuration for updating Pirsch domains
	UpdatePirschDomainWorker corejobs.UpdatePirschDomainWorker `koanf:"updatepirschdomainworker" json:"updatepirschdomainworker"`

	// PreviewDomainWorkers configuration for preview domain workers
	// CreatePreviewDomainWorker configuration for creating preview domains
	CreatePreviewDomainWorker corejobs.CreatePreviewDomainWorker `koanf:"createpreviewdomainworker" json:"createpreviewdomainworker"`
	// DeletePreviewDomainWorker configuration for deleting preview domains
	DeletePreviewDomainWorker corejobs.DeletePreviewDomainWorker `koanf:"deletepreviewdomainworker" json:"deletepreviewdomainworker"`
	// ValidatePreviewDomainWorker configuration for validating preview domains
	ValidatePreviewDomainWorker corejobs.ValidatePreviewDomainWorker `koanf:"validatepreviewdomainworker" json:"validatepreviewdomainworker"`

	// AttestNDARequestWorker configuration for attesting NDA requests
	AttestNDARequestWorker corejobs.AttestNDARequestWorker `koanf:"attestndarequestworker" json:"attestndarequestworker"`

	// CreatePreviewDomainAcmeWorker for storing the acme content from cloudflare so we can serve it when request
	CreatePreviewDomainAcmeWorker corejobs.CreatePreviewDomainAcmeWorker `koanf:"createpreviewdomainacmeworker" json:"createpreviewdomainacmeworker"`

	// CreateDomainScanWorker for pinging cloudflare to start the scanning
	CreateDomainScanWorker corejobs.CreateDomainScanWorker `koanf:"createdomainscanworker" json:"createdomainscanworker"`

	// RetrieveDomainScanWorker for fetching the scan data from cloudflare and creating a notification
	RetrieveDomainScanWorker corejobs.RetrieveDomainScanWorker `koanf:"retrievedomainscanworker" json:"retrievedomainscanworker"`

	// OrganizationDeletionReminderWorker configuration for sending organization deletion reminders
	OrganizationDeletionReminderWorker corejobs.OrganizationPaymentReminderWorker `koanf:"organizationdeletionreminderworker" json:"organizationdeletionreminderworker"`

	// OrganizationDeletionWorker configuration for deleting organizations
	OrganizationDeletionWorker corejobs.OrganizationDeleteWorker `koanf:"organizationdeletionworker" json:"organizationdeletionworker"`

	// add more trust center specific workers here
}

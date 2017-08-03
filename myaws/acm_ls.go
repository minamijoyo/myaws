package myaws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/pkg/errors"
)

// ACMLsOptions customize the behavior of the Ls command.
type ACMLsOptions struct {
	Verbose bool
	Pending bool
}

// ACMLs describes Amazon Certificate resources.
func (client *Client) ACMLs(options ACMLsOptions) error {
	certificates, err := client.listACMCertificates()
	if err != nil {
		return err
	}
	for _, cert := range certificates {
		if options.Pending && !hasAnyPendingACMCertificate(cert) {
			continue
		}

		if options.Verbose {
			fmt.Fprintln(client.stdout, cert)
		} else {
			fmt.Fprintln(client.stdout, client.FormatTime(cert.NotAfter), *cert.DomainName)
		}
	}

	return nil
}

func (client *Client) listACMCertificates() ([]*acm.CertificateDetail, error) {
	var certificates []*acm.CertificateDetail
	params := &acm.ListCertificatesInput{}

	response, err := client.ACM.ListCertificates(params)
	if err != nil {
		return certificates, errors.Wrap(err, "ListCertificates failed:")
	}

	for {
		for _, cert := range response.CertificateSummaryList {
			describeParams := &acm.DescribeCertificateInput{
				CertificateArn: cert.CertificateArn,
			}
			r, err := client.ACM.DescribeCertificate(describeParams)
			if err != nil {
				return certificates, errors.Wrap(err, "DescribeCertificate failed:")
			}
			certificates = append(certificates, r.Certificate)
		}

		if response.NextToken == nil {
			break
		}

		params.SetNextToken(*response.NextToken)
		response, err = client.ACM.ListCertificates(params)
		if err != nil {
			return certificates, errors.Wrap(err, "ListCertificates with next token failed:")
		}
	}

	return certificates, nil
}

// hasAnyPendingACMCertificate returns true if there is at least one "pending verification" in either domain.
func hasAnyPendingACMCertificate(certificate *acm.CertificateDetail) bool {
	if certificate.RenewalSummary == nil {
		return false
	}

	renewalSummary := certificate.RenewalSummary
	validationOptions := renewalSummary.DomainValidationOptions

	for _, validation := range validationOptions {
		if *validation.ValidationStatus == "PENDING_VALIDATION" {
			return true
		}
	}

	return false
}

/*
Copyright 2021 The cert-manager Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	cmapi "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
//+kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=`.status.conditions[?(@.type == "Ready")].status`,description="CertificateRequestPolicy is ready for evaluation"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Timestamp CertificateRequestPolicy was created"
//+kubebuilder:resource:categories=cert-manager,shortName=crp,scope=Cluster
//+kubebuilder:subresource:status

// CertificateRequestPolicy is an object for describing a "policy profile" that
// makes decisions on whether applicable CertificateRequests should be approved
// or denied.
type CertificateRequestPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CertificateRequestPolicySpec   `json:"spec,omitempty"`
	Status CertificateRequestPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// CertificateRequestPolicyList is a list of CertificateRequestPolicies.
type CertificateRequestPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CertificateRequestPolicy `json:"items"`
}

// CertificateRequestPolicySpec defines the desired state of
// CertificateRequestPolicy.
type CertificateRequestPolicySpec struct {
	// Allowed is the set of attributes that are "allowed" by this policy. A
	// CertificateRequest will only be considered permissible for this policy if
	// the CertificateRequest has the same or less as what is allowed.  Empty or
	// `nil` allowed fields mean CertificateRequests are not allowed to have that
	// field present to be permissible.
	// +optional
	Allowed *CertificateRequestPolicyAllowed `json:"allowed,omitempty"`

	// Constraints is the set of attributes that _must_ be satisfied by the
	// CertificateRequest for the request to be permissible by the policy. Empty
	// or `nil` constraint fields mean CertificateRequests satisfy that field
	// with any value of their corresponding attribute.
	// +optional
	Constraints *CertificateRequestPolicyConstraints `json:"constraints,omitempty"`

	// Plugins define a set of plugins and their configuration that should be
	// executed when this policy is evaluated against a CertificateRequest. A
	// plugin must already be built within policy-approver for it to be
	// available.
	// +optional
	Plugins map[string]CertificateRequestPolicyPluginData `json:"plugins,omitempty"`

	// Selector is used for selecting over which CertificateRequests this
	// CertificateRequestPolicy is appropriate for and so will used for its
	// evaluation.
	Selector CertificateRequestPolicySelector `json:"selector"`
}

// CertificateRequestPolicyAllowed is a set of attributes that are declared as
// permissible for a CertificateRequest to have those values present. It is
// permissible for a CertificateRequest to request _less_ than what is allowed,
// but _not more_, i.e. it is permissible for a CertificateRequest to request a
// subset of what is allowed.
// Empty fields or `nil` values declares that the equivalent CertificateRequest
// field _must_ be omitted or empty for the request to be permitted.
type CertificateRequestPolicyAllowed struct {
	// CommonName defines the X.509 Common Name that is permissible.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids a CommonName being requested.
	// An empty string is equivalent to `nil`.
	// +optional
	CommonName *string `json:"commonName,omitempty"`

	// DNSNames defines the X.509 DNS SANs that may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any DNSNames being requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	DNSNames *[]string `json:"dnsNames,omitempty"`

	// IPAddresses defines the X.509 IP SANs that may be requested
	// for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any IPAddresses from being
	// requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	IPAddresses *[]string `json:"ipAddresses,omitempty"`

	// URIs defines the X.509 URI SANs that may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any URIs from being requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	URIs *[]string `json:"uris,omitempty"`

	// EmailAddresses defines the X.509 Email SANs that may be
	// requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any EmailAddress from being
	// requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	EmailAddresses *[]string `json:"emailAddresses,omitempty"`

	// IsCA defines whether it is permissible for a CertificateRequest to have
	// the `spec.IsCA` field set to `true`.
	// An omitted field, value of `nil` or `false`, forbids the `spec.IsCA` field
	// from bring `true`.
	// A value of `true` permits CertificateRequests setting the `spec.IsCA` field
	// to `true`.
	// +optional
	IsCA *bool `json:"isCA,omitempty"`

	// Usages defines the list of permissible key usages that may appear
	// of the CertificateRequest `spec.keyUsages` field.
	// An omitted field or value of `nil` forbids any Usages being requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	Usages *[]cmapi.KeyUsage `json:"usages,omitempty"`

	// Subject defines the X.509 subject that is permissible. An omitted field or
	// value of `nil` forbids any Subject being requested.
	// +optional
	Subject *CertificateRequestPolicyAllowedX509Subject `json:"subject,omitempty"`
}

// CertificateRequestPolicyAllowedX509Subject declares the X.509 Subject
// attributes that are permissible for a CertificateRequest to request for this
// policy. It is permissible for CertificateRequests to request a subset of
// Allowed X.509 Subject attributes defined.
type CertificateRequestPolicyAllowedX509Subject struct {
	// Organizations define the X.509 Subject Organizations that may be requested
	// for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Organization from being
	// requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	Organizations *[]string `json:"organizations,omitempty"`

	// Countries define the X.509 Subject Countries that may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Countries from being
	// requested.
	// An empty slice `[]` is equivalent to `nil`.
	// +optional
	Countries *[]string `json:"countries,omitempty"`

	// OrganizationalUnits defines the X.509 Subject Organizational Units that
	// may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any OrganizationalUnits from
	// being requested.
	// An empty slice `[]` is equivalent to nil.
	// +optional
	OrganizationalUnits *[]string `json:"organizationalUnits,omitempty"`

	// Localities defines the X.509 Subject Localities that may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Localities from being
	// requested.
	// An empty slice `[]` is equivalent to nil.
	// +optional
	Localities *[]string `json:"localities,omitempty"`

	// Provinces defines the X.509 Subject Provinces that may be requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Provinces from being
	// requested.
	// An empty slice `[]` is equivalent to nil.
	// +optional
	Provinces *[]string `json:"provinces,omitempty"`

	// StreetAddresses defines the X.509 Subject Street Addresses that may be
	// requested for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Street Addresses from being
	// requested.
	// An empty slice `[]` is equivalent to nil.
	// +optional
	StreetAddresses *[]string `json:"streetAddresses,omitempty"`

	// PostalCodes defines the X.509 Subject Postal Codes that may be requested
	// for.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` forbids any Postal Codes from being
	// requested.
	// An empty slice `[]` is equivalent to nil.
	// +optional
	PostalCodes *[]string `json:"postalCodes,omitempty"`

	// SerialNumber defines the X.509 Subject Serial Number that may be requested
	// for.
	// An omitted field or value of `nil` forbids any Serial Number from being
	// requested.
	// +optional
	SerialNumber *string `json:"serialNumber,omitempty"`
}

// CertificateRequestPolicyConstraints define fields that, if defined, _must_
// be satisfied by the CertificateRequest for the request to be permissible by
// this policy. Fields that are omitted or have a value of `nil` will be
// satisfied by any value on the corresponding attribute on the request.
type CertificateRequestPolicyConstraints struct {
	// MinDuration defines the minimum duration a certificate may be requested
	// for.
	// Values are inclusive (i.e. a min value of `1h` will accept a duration of
	// `1h`). MinDuration and MaxDuration may be the same value.
	// An omitted field or value of `nil` permits any minimum duration.
	// If MinDuration is defined, a duration _must_ be requested on the
	// CertificateRequest.
	// +optional
	MinDuration *metav1.Duration `json:"minDuration,omitempty"`

	// MaxDuration defines the maximum duration a certificate may be requested
	// for.
	// Values are inclusive (i.e. a max value of `1h` will accept a duration of
	// `1h`). MaxDuration and MinDuration may be the same value.
	// An omitted field or value of `nil` permits any maximum duration.
	// If MaxDuration is defined, a duration _must_ be requested on the
	// CertificateRequest.
	// +optional
	MaxDuration *metav1.Duration `json:"maxDuration,omitempty"`

	// PrivateKey defines the shape of permissible private keys that may be used
	// for the request with this policy.
	// An omitted field or value of `nil` permits the use of any private key by
	// the requestor.
	// +optional
	PrivateKey *CertificateRequestPolicyConstraintsPrivateKey `json:"privateKey,omitempty"`
}

// CertificateRequestPolicyConstraintsPrivateKey defines constraints on what
// shape of private key is permissible for a CertificateRequest to have used
// for its request.
type CertificateRequestPolicyConstraintsPrivateKey struct {
	// Algorithm defines the allowed crypto algorithm that is used by the
	// requestor for their private key in their request.
	// An omitted field or value of `nil` permits any Algorithm.
	// +optional
	Algorithm *cmapi.PrivateKeyAlgorithm `json:"Algorithm,omitempty"`

	// MinSize defines the minimum key size a requestor may use for their private
	// key.
	// Values are inclusive (i.e. a min value of `2048` will accept a size
	// of `2048`). MinSize and MaxSize may be the same value.
	// An omitted field or value of `nil` permits any minimum size.
	// +optional
	MinSize *int `json:"minSize,omitempty"`

	// MaxSize defines the maximum key size a requestor may use for their private
	// key.
	// Values are inclusive (i.e. a min value of `2048` will accept a size
	// of `2048`). MaxSize and MinSize may be the same value.
	// An omitted field or value of `nil` permits any maximum size.
	// +optional
	MaxSize *int `json:"maxSize,omitempty"`
}

// CertificateRequestPolicyPluginData is configuration needed by the plugin
// approver to evaluate a CertificateRequest on this policy.
type CertificateRequestPolicyPluginData struct {
	// Values define a set of well-known, to the plugin, key value pairs that are
	// required for the plugin to successfully evaluate a request based on this
	// policy.
	// +optional
	Values map[string]string `json:"values,omitempty"`
}

// CertificateRequestPolicySelector is used for selecting over the
// CertificateRequests what this CertificateRequestPolicy is appropriate for,
// and if so, will be used to evaluate the request.
type CertificateRequestPolicySelector struct {
	// IssuerRef is used to match this CertificateRequestPolicy against processed
	// CertificateRequests. This policy will only be evaluated against a
	// CertificateRequest whose `spec.issuerRef` field matches
	// `spec.selector.issuerRef`. CertificateRequests will not be processed on
	// unmatched `issuerRef`, regardless of whether the requestor is bound by
	// RBAC.
	// Accepts wildcards "*".
	// Nil values are equivalent to "*",
	//
	// The following value will match _all_ `issuerRefs`:
	// ```
	// issuerRef: {}
	// ```
	//
	// Required field.
	IssuerRef *CertificateRequestPolicySelectorIssuerRef `json:"issuerRef"`
}

// CertificateRequestPolicySelectorIssuerRef defines the selector for matching
// on `issuerRef` of requests.
type CertificateRequestPolicySelectorIssuerRef struct {
	// Name is the wildcard selector to match the `spec.issuerRef.name` field on
	// requests.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` matches all.
	// +optional
	Name *string `json:"name,omitempty"`

	// Kind is the wildcard selector to match the `spec.issuerRef.kind` field on
	// requests.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` matches all.
	// +optional
	Kind *string `json:"kind,omitempty"`

	// Group is the wildcard selector to match the `spec.issuerRef.group` field
	// on requests.
	// Accepts wildcards "*".
	// An omitted field or value of `nil` matches all.
	// +optional
	Group *string `json:"group,omitempty"`
}

// CertificateRequestPolicyStatus defines the observed state of the
// CertificateRequestPolicy.
type CertificateRequestPolicyStatus struct {
	// List of status conditions to indicate the status of the
	// CertificateRequestPolicy.
	// Known condition types are `Ready`.
	// +optional
	Conditions []CertificateRequestPolicyCondition `json:"conditions,omitempty"`
}

// CertificateRequestPolicyCondition contains condition information for a
// CertificateRequestPolicyStatus.
type CertificateRequestPolicyCondition struct {
	// Type of the condition, known values are (`Ready`).
	Type CertificateRequestPolicyConditionType `json:"type"`

	// Status of the condition, one of ('True', 'False', 'Unknown').
	Status corev1.ConditionStatus `json:"status"`

	// LastTransitionTime is the timestamp corresponding to the last status
	// change of this condition.
	// +optional
	LastTransitionTime *metav1.Time `json:"lastTransitionTime,omitempty"`

	// Reason is a brief machine readable explanation for the condition's last
	// transition.
	// +optional
	Reason string `json:"reason,omitempty"`

	// Message is a human readable description of the details of the last
	// transition, complementing reason.
	// +optional
	Message string `json:"message,omitempty"`

	// If set, this represents the .metadata.generation that the condition was
	// set based upon.
	// For instance, if .metadata.generation is currently 12, but the
	// .status.condition[x].observedGeneration is 9, the condition is out of date
	// with respect to the current state of the CertificateRequestPolicy.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

// CertificateRequestPolicyConditionType represents a CertificateRequestPolicy
// condition value.
type CertificateRequestPolicyConditionType string

const (
	// CertificateRequestPolicyConditionReady indicates that the
	// CertificateRequestPolicy has successfully loaded the policy, and all
	// configuration including plugin options are accepted and ready for
	// evaluating CertificateRequests.
	CertificateRequestPolicyConditionReady CertificateRequestPolicyConditionType = "Ready"
)

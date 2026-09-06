package domain

type UsageRequest struct {
	ModelID      string
	IntendedUse  string
	Territory    string
	HasValidKYC  bool
	RequesterDID string
}

type LicensingDecision struct {
	Allowed bool
	Reason  string
}

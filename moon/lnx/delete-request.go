package lnx

type deleteRequest struct {
	Query query  `json:"query"`
	Limit uint64 `json:"limit"`
}

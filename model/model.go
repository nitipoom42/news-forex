package model

const (
	CurRencyUSD      = "USD"
	ImpactHigh       = "icon--ff-impact-red"
	ImpactMedium     = "icon--ff-impact-yel"
	ImpactLow        = "icon--ff-impact-ora"
	ImpactUnknown    = "icon--ff-impact-gra"
	ResImpactHigh    = "High"
	ResImpactMed     = "Medium"
	ResImpactLow     = "Low"
	ResImpactUnknown = "Unknown"
	AllDay           = "ALL DAY"
)

type NewsEvent struct {
	Date  string `json:"date"`
	Time  string `json:"time"`
	Title string `json:"title"`
}

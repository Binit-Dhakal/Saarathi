package domain

type CarPackage string

type Fare struct {
	Package    CarPackage
	TotalPrice int
}

// temporary fare estimate for given route for redis persistance
type FareQuote struct {
	Route Route  `json:"route"`
	Fares []Fare `json:"fares"`
}

type FareRecord struct {
	RouteID string
	Fare    Fare
}

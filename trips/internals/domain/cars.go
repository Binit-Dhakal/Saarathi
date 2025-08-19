package domain

const (
	SEDAN  CarPackage = "SEDAN"
	SUV    CarPackage = "SUV"
	VAN    CarPackage = "VAN"
	LUXURY CarPackage = "LUXURY"
)

type Car struct {
	Name          CarPackage
	BaseFare      int
	PerKmRate     int
	PerMinuteRate int
}

var CarRegistry = []Car{
	{Name: SEDAN, BaseFare: 50, PerKmRate: 40, PerMinuteRate: 1},
	{Name: SUV, BaseFare: 80, PerKmRate: 55, PerMinuteRate: 2},
	{Name: VAN, BaseFare: 100, PerKmRate: 65, PerMinuteRate: 2},
	{Name: LUXURY, BaseFare: 200, PerKmRate: 120, PerMinuteRate: 5},
}

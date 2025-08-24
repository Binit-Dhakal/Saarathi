package messagebus

func DriverRoutingKey(driverID string) string {
	return DriverInstancePrefix + driverID
}

func RideMatchingRoutingKey(instanceID string) string {
	return RideMatchingInstancePrefix + instanceID
}

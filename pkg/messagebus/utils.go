package messagebus

func DriverRoutingKey(instanceID string) string {
	return DriverInstancePrefix + instanceID
}

func RideMatchingRoutingKey(instanceID string) string {
	return RideMatchingInstancePrefix + instanceID
}

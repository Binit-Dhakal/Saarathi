package messagebus

func DriverRoutingKey(eventName string, instanceID string) string {
	return eventName + "." + instanceID
}

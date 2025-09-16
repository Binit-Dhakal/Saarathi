# right now temporary; will polish it later
config-create:
	kubectl create configmap db-migrations --from-file=migrations/ 

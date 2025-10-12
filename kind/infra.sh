#!/bin/bash

# redis url: redis://redis-master:6379
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install redis bitnami/redis --set architecture=standalone --set auth.enabled=false

# nats.default.svc.cluster.local:4222
helm repo add nats https://nats-io.github.io/k8s/helm/charts
helm install nats nats/nats

# postgres://devuser:devpass@userspq-postgresql.default.svc.cluster.local:5432/users_db
helm repo add bitnami https://charts.bitnami.com/bitnami
helm install userspq bitnami/postgresql --set primary.service.name=users_postgres --set auth.database=usersdb --set auth.username=devuser --set auth.password=devpass --set auth.postgresPassword=supersecret 

helm install tripspq bitnami/postgresql --set auth.password=devpass --set primary.service.name=trips_postgres --set auth.database=tripsdb --set auth.username=devuser --set auth.postgresPassword=supersecret 

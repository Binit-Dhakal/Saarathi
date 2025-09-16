default_registry("kind-registry:5001")

docker_build(
    'api-gateway', 
    context='.',
    dockerfile='./api-gateway/Dockerfile.dev',
)

docker_build(
    'trips-service',
    context='.',
    dockerfile='./trips/Dockerfile.dev'
)

docker_build(
    'users-service',
    context='.',
    dockerfile='./users/Dockerfile.dev'
)

k8s_yaml(kustomize('./k8s/overlays/dev'))

k8s_resource('api-gateway', port_forwards=8080)


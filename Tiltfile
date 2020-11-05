yaml = helm(
  'deployments/go-http-server',
  name='go-http-server-dev',
  namespace='go-http-server-dev',
  set=[
      "service.type=NodePort"
  ]
  )
k8s_yaml(yaml)
docker_build('registry.gitlab.com/safesurfer/go-http-server', '.', dockerfile="build/Dockerfile")
allow_k8s_contexts('in-cluster')

load('ext://ko', 'ko_build')
yaml = helm(
  'deployments/go-http-server',
  name='go-http-server-dev',
  namespace='go-http-server-dev',
  set=[
      "service.type=NodePort"
  ]
  )
k8s_yaml(yaml)
ko_build('registry.gitlab.com/BobyMCbobs/go-http-server',
         '.',
         deps=['./'])
allow_k8s_contexts('in-cluster')

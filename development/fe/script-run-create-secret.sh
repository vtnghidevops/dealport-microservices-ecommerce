run-this
kubectl create secret docker-registry auth-registry \ 
--docker-email=yourmail@gmail.com --docker-username=username-harbor \ 
--docker-password=password-harbor --docker-server=domain-harbor.com \ 
--namespace ecommerce
or
apiVersion: v1
kind: Secret
metadata:
  name: auth-registry
  namespace: ecommerce
type: kubernetes.io/dockerconfigjson
data:
  .dockerconfigjson: <base64-encoded-docker-config>

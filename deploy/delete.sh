set -a
. ./$1.env
set +a
cat k8s.yaml.tmpl | envsubst > k8s.yaml
kubectl delete -f k8s.yaml

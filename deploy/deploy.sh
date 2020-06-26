TAG=$(env LC_CTYPE=C tr -dc "a-zA-Z0-9" < /dev/urandom | head -c 10)
(cd .. && make docker-deploy TAG=$TAG)

set -a
DOCKER_IMAGE="zsiec/duck:${TAG}"
. ./$1.env
set +a

cat k8s.yaml.tmpl | envsubst > k8s.yaml
kubectl apply -f k8s.yaml

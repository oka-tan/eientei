#Namespace
kubectl create namespace eientei

#Ingress
#Consider replacing with a properly configured ingress
kubectl apply -f kubernetes/traefik.yaml

#s3 stateful single-node minio instance
#Consider using the minio operator for production deployments
#Otherwise remember to reconfigure the volume size and the passwords
kubectl apply -f kubernetes/s3.yaml

#Postgres stateful single-node instance
#Consider using a postgres operator for production deployments
#Otherwise remember to reconfigure the volume size and the passwords
kubectl apply -f kubernetes/postgres.yaml

#Load config files into configmaps
kubectl delete configmap -n eientei kaguya-config
kubectl create configmap -n eientei kaguya-config --from-file=kaguya.json
kubectl delete configmap -n eientei reisen-config
kubectl create configmap -n eientei reisen-config --from-file=reisen.json
kubectl delete configmap -n eientei moon-config
kubectl create configmap -n eientei moon-config --from-file=moon.json

#Stateless kaguya instance
#Booting more than one is not recommended in general
#I mean you can, but you can scrape every board with one instance no sweat
kubectl apply -f kubernetes/kaguya.yaml

#Stateless reisen instance
kubectl apply -f kubernetes/reisen.yaml

#Stateless moon instance
#Booting more than one is not recommended in general
#I mean you can, but you can index every board with one instance no sweat
kubectl apply -f kubernetes/moon.yaml

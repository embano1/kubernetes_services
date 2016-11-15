# kubernetes_services
Tiny program to ease the use of getting NodePorts for exposed Kubernetes services 
(e.g. not running an external Load Balancer, local testing, etc.). Queries the API server. You´d usually want to run "kubectl proxy" locally to proxy to the API server.

# Clone
git clone https://github.com/embano1/kubernetes_services.git  

# Build
cd kubernetes_services/  
go build -o kubernetes_services *.go  

# Run
./kubernetes_services  

# Help and options
./kubernetes_services -h

# Docker-Compose + Consul + Go 🎯

I haven't found too much information on how to setup a Go service and Consul running on top of Docker using `docker-compose`, so I decided to create this little example setup.

1. Clone it and run `docker build -t mgyongyosi/weather .` inside the repository's folder
2. Start with `docker-compose up`
3. Check the Consul Web-UI on `localhost:8500/ui`
4. Check the Go service on `localhost:8080/` and the 'heartbeat' on `localhost:8080/health` (it returns only `200 OK`)

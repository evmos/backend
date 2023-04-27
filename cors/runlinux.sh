#!/bin/bash
docker build -f Dockerfile.linux -t h/nginx .
docker run --add-host=host.docker.internal:host-gateway -p 80:80 h/nginx

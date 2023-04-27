#!/bin/bash
docker build -f Dockerfile.local -t h/nginx .
docker run -p 80:80 h/nginx

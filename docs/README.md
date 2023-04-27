# Docs

## Production configuration

### Screens running

- Redis:

```sh
redis-server
```

- Price:

```sh
cd dashboard-backend/cronjobs
python price.py
```

- Endpoints:

```sh
cd dashboard-backend/cronjobs
python cron.py
```

- API Backend:

```sh
cd dashboard-backend
go build
GOMAXPROCS=4 ./dashboard-backend
```

### Settings

- NginxEndpoint (`/etc/nginx/sites-available/newapi`)

```
limit_req_zone $binary_remote_addr zone=limitreqsbyaddr:20m rate=100r/s;

server {
    server_name goapi.hanchon.live;
    location / {
        limit_req zone=limitreqsbyaddr burst=100 nodelay;

        add_header Access-Control-Allow-Origin *;
        proxy_pass http://127.0.0.1:8081/;
        include proxy_params;


       if ($request_method = 'OPTIONS') {
        add_header 'Access-Control-Allow-Origin' '*';

        add_header 'Access-Control-Allow-Credentials' 'true';
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';

        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';

        add_header 'Access-Control-Max-Age' 86400;
        add_header 'Content-Type' 'text/plain charset=UTF-8';
        add_header 'Content-Length' 0;
        return 204; break;
     }

     if ($request_method = 'POST') {
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Credentials' 'true';
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';
     }
     if ($request_method = 'GET') {
        add_header 'Access-Control-Allow-Origin' '*';
        add_header 'Access-Control-Allow-Credentials' 'true';
        add_header 'Access-Control-Allow-Methods' 'GET, POST, OPTIONS';
        add_header 'Access-Control-Allow-Headers' 'DNT,X-CustomHeader,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type';
     }

    proxy_set_header Host      $host:$server_port;
    proxy_set_header X-Real-IP $remote_addr;

    }


    listen 443 ssl reuseport; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/goapi.hanchon.live/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/goapi.hanchon.live/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot

    ssl_session_cache shared:SSL:100m;
    ssl_session_cache shared:AdminSSL:100m;
}
server {
    if ($host = goapi.hanchon.live) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    server_name goapi.hanchon.live;
    listen 80 reuseport;
    return 404; # managed by Certbot


}
```

- Add to nginx config (`/etc/nginx/nginx.conf`)

```
worker_rlimit_nofile 200000;

events {
	worker_connections 200000;
}
```

## Add new networks

- Add the new endpoints to the `cronjobs/constants.py` file
- Add the new tokens to the `internal/constants/erc20.go` file
- Add the ibc channels and denom to the `constants/ibc.go` file
- Add the network to the `internal/constants/networks.go` file

## Architecture

![diagram](./architecture.png)

# Sample configs

- **Save configuration to:** /etc/nginx/sites-enabled/data-listener

- **Test commands for http to unix/http**

        curl -X POST -H "Content-Type: application/json" -d '{"key1": "value1", "key2": "value2"}' http://localhost:80
        go-wrk -c 450 -n 40000 -m POST -H "Content-Type: application/json" -p benchmarks/json/data1.json -i http://localhost:80/

## HTTP to HTTP

        server {
            listen 80;
            server_name localhost;

            location / {
                include proxy_params;
                proxy_pass http://127.0.0.1:8080;
            }
        }

## HTTP to Unix Socket

        server {
            listen 80;

            server_name localhost;

            location / {
                include proxy_params;
                proxy_pass <http://unix:/tmp/data-listener/fasthttp.sock>;
            }
        }

## HTTPS to HTTP

        server {
            listen 443 ssl;
            server_name localhost;
        
            ssl_certificate /path/to/your/certificate.crt;
            ssl_certificate_key /path/to/your/private.key;
        
            location / {
                include proxy_params;
                proxy_pass http://127.0.0.1:8080;;
            }
        }

## HTTPS to Unix Socket

        server {
            listen 443 ssl;
            server_name localhost;

            ssl_certificate /path/to/your/certificate.crt;
            ssl_certificate_key /path/to/your/private.key;

            location / {
                include proxy_params;
                proxy_pass http://unix:/tmp/data-listener/fasthttp.sock:;
            }
        }

## Additional Information

- include proxy_params;

        The proxy_params file is usually provided by default with Nginx installations,
        and it contains common proxy header configurations that are frequently used
        when setting up Nginx as a reverse proxy. These headers are important for
        forwarding information about the original client request to the proxied server.

        location / {
            proxy_pass http://127.0.0.1:8080;;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

# Create a systemd service unit and configure it

        1. sudo nano /etc/systemd/system/data-listener.service

## Create a systemd socket unit file and configure it

        2. sudo nano /etc/systemd/system/data-listener.socket

## Create executable

        3. go build -o data-listener-exec ./cmd/main.go
        3.1 cp https/ data-listener-exec log/ output/ config.toml <desired path> </tmp/data-listener/>

## create system user/group example - dl-development, and give it rights for the executable, root folder, socket folder

        4.1 sudo adduser dl-development
        4.2 sudo addgroup dl-development
        4.3 sudo chown dl-development:dl-development /opt/data-listener/data-listener-exec

## Reload systemd to pick up the c       hanges

        5. sudo systemctl daemon-reload  

## Enable and start the socket and service

        6.1 sudo systemctl daemon-reload
        6.2 sudo systemctl enable data-listener.socket
        6.3 sudo systemctl enable data-listener.service
        6.4 sudo systemctl start data-listener.socket
        6.5 sudo systemctl start data-listener.service

## Verify that the services are running

        7.1 sudo systemctl status data-listener.socket
        7.2 sudo systemctl status data-listener.service

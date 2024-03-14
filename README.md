# Data-Listener

*by Gudasoft*

## Overview

Data-Listener is a GoLang application designed to efficiently receive and resend packets of incoming requests. It offers flexibility through a highly configurable `config.toml` file and supports both direct and buffered packet handling.

## Functionality

- **Configurability**: Easily tweak the application behavior through the `config.toml` file.
- **Graceful Shutdowns**: Implements graceful shutdowns for a smoother termination process.
- **Configuration Reload**: Dynamically reload configurations while running live with the SIGHUP signal.

### Available Listeners

- **HTTP**: Handles incoming requests efficiently.

### Available Streamers

- **File**: Allows streaming of data to files.
- **AWS S3**: Enables streaming to Amazon S3.

## Running

1. **Prerequisites**: Ensure you have GoLang version 1.20.3 installed.

2. **Configuration**: Create a `config.toml` file in the main folder and configure your settings. You can refer to the example in `config-example.toml`.

3. **Application Start**: Run the following command inside the main directory to start the application:

   ```bash
   go run ./cmd/ .

4. **AWS Authentication**: If using AWS services, make sure you are authenticated in the AWS CLI.

### Feel free to explore and adapt Data-Listener to fit your specific use case. For any issues or suggestions, please reach out to us.

# Thinks to check and finish

- If you plan to store input data as files, create folders "output" & "output2" in the main directory

Create examples folder

add example for:

  - storing objects on s3 - store-on-s3.toml
  - store files locally - store-unbuffered-locally.toml
  - store buffered files locally - store-buffered-locally.toml
  - full config with working options - all-options.toml
  - roadmap config  - roadmap-config.toml
  - development config  - development-config.toml


# Releasing a new version

on sarah:

   sudo su -
   cd /etc/nginx
   export DATA_LISTENER_VERSION=1.0.0
   mkdir -p /etc/nginx/htpasswd-data-listener/
   touch /etc/nginx/htpasswd-data-listener/$DATA_LISTENER_VERSION
   ./get_password.sh -b /etc/nginx/htpasswd-data-listener/$DATA_LISTENER_VERSION $DATA_LISTENER_VERSION hellomarshmellow

Add this in sites-enabled/data-listener.conf

    location /releases/1.0.0 {
        autoindex on;
        autoindex_exact_size off;
        autoindex_localtime on;

        auth_basic "Restricted Access";
        auth_basic_user_file /etc/nginx/htpasswd-data-listener/1.0.0;
    }


   sudo systemctl reload nginx


Check it https://data-listener.gudasoft.com/releases/1.0.0



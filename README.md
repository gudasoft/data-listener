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

- **Prerequisites**: Ensure you have GoLang version 1.20.3 installed.

- **Configuration**: Create a `config.toml` file in the main folder and configure your settings. You can refer to the example in `config-example.toml`.

- **Application Start**: Run the following command inside the main directory to start the application:

```bash
   go run ./cmd/ .
```

- **AWS Authentication**: If using AWS services, make sure you are authenticated in the AWS CLI.

### Feel free to explore and adapt Data-Listener to fit your specific use case. For any issues or suggestions, please reach out to us.

# Things to check and finish

- If you plan to store input data as files, create folders "output" & "output2" in the main directory

Create examples folder add example for:

  - storing objects on s3 - store-on-s3.toml
  - store files locally - store-unbuffered-locally.toml
  - store buffered files locally - store-buffered-locally.toml
  - full config with working options - all-options.toml
  - roadmap config  - roadmap-config.toml
  - development config  - development-config.toml


# Downloading the latest release

Check it https://data-listener.gudasoft.com/releases/1.0.0




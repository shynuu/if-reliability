# Interface Reliability and Switch Script

## Overview

`if-reliability` is a command-line tool designed to monitor and manage the reliability of network interfaces. It provides functionality to check the connectivity of an interface and switch between LTE and WiFi based on connectivity status.

## Features

- Pings a specified endpoint to check connectivity.
- Automatically switches to WiFi if the LTE interface is unreliable.
- Configurable retry mechanism to determine connectivity loss.
- Command-line flags for easy configuration of interfaces and parameters.

## Usage

```bash
./if-reliability --lte-if <LTE_INTERFACE> --wifi-if <WIFI_INTERFACE> --wifi-ssid <WIFI_SSID> --wifi-password <WIFI_PASSWORD> --endpoint <ENDPOINT> --failure-delay <DELAY_IN_MS>
```

### Flags

- `--lte-if` (`-l`): The LTE interface to monitor (required).
- `--wifi-if` (`-w`): The WiFi interface to switch to (required).
- `--wifi-ssid` (`-s`): The SSID of the WiFi network (required).
- `--wifi-password` (`-p`): The password for the WiFi network (required).
- `--endpoint` (`-e`): The endpoint to ping for connectivity checks (required).
- `--retry` (`-r`): The number of failed attempts before switching to WiFi (default: 5).

## Installation

Ensure you have Go installed on your system. Clone the repository and build the tool using the provided `compile.sh` script:

```bash
git clone https://github.com/shynuu/if-reliability
cd if-reliability
./compile.sh main.go
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributions

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## Contact

For questions or suggestions, please contact Youssouf Drif <youssouf.drif@uni.lu>.

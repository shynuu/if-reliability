# Interface Reliability Tool

Interface Reliability Tool is a command-line utility designed to check the reliability of a network interface by pinging a specified endpoint and automatically switching to a WiFi connection if the primary connection fails.

## Features

- Ping a specified endpoint to check connectivity.
- Automatically switch to a specified WiFi network upon failure.
- Customize retry count for failure detection.
- Replace default route with the WiFi network's route upon successful connection.

## Installation

1. Clone the repository:

   ```
   git clone <repository-url>
   ```

2. Navigate to the project directory:

   ```
   cd <project-directory>
   ```

3. Build the project using the provided `compile.sh` script:

   ```
   ./compile.sh main.go
   ```

## Usage

Run the Interface Reliability Tool with the required flags:

```
./if-reliability --wifi-if <wifi-interface> --wifi-ssid <wifi-ssid> --wifi-password <wifi-password> --endpoint <endpoint> [--retry <retry-count>]
```

- `--wifi-if`: WiFi interface name (required)
- `--wifi-ssid`: WiFi SSID (required)
- `--wifi-password`: WiFi password (required)
- `--endpoint`: Endpoint to check connectivity (required)
- `--retry`: Number of retries before switching to WiFi (default: 5)

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## Contact

For questions or issues, please contact **Youssouf Drif** <youssouf.drif@uni.lu>.

## Project

The Interface Reliability Tool is developed as part of the SnT (Interdisciplinary Centre for Security, Reliability and Trust) initiatives at the University of Luxembourg. It aims to improve network interface management and reliability, aligning with the cutting-edge research and development objectives of SnT.

For more info please visit https://6gspacelab.uni.lu

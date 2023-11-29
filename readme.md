# UserDataProcessor

UserDataProcessor is a server application designed to handle and process user-generated data. It validates incoming JSON payloads, applies user-specific quotas, and stores or forwards data as needed.

## Features

- **Data Receipt**: Accept user data through HTTP POST requests.
- **Data Validation**: Ensure data integrity through JSON schema validation.
- **Quota Management**: Enforce data upload limits to prevent system overuse.
- **Request Throttling**: Limit user request frequency to maintain system performance.

## Getting Started

### Prerequisites

Ensure you have Go installed on your system. You can download it from [the Go website](https://golang.org/dl/).

### Installation

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/UserDataProcessor.git
    ```
2. Navigate to the project directory:
    ```bash
    cd UserDataProcessor
    ```
3. Compile the application:
    ```bash
    go build
    ```

### Usage
Run the compiled application:
```bash
./UserDataProcessor
```

The service will start and listen for requests at `localhost:8080`.

## API Endpoint

### POST /data
Receives JSON-formatted user data, validates it, and applies user quotas.

- Request Body
    ```JSON
    {
        "id": "<unique_data_id>",
        "userID": "<user_identifier>",
        "data": "<user_data_string>"
    }
    ```

- Response

    A JSON with the result of the operation and HTTP status code.


## Configuration
The application runs with default settings that limit each user to 100MB of data and 10 requests per minute. These settings can be customized in the application’s configuration file or environment variables.

## Testing
Run the following command to execute the unit tests:

```bash
go test ./...
```
## Contributing
Contributions are what make the open-source community an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

## License
Distributed under the MIT License. See LICENSE for more information.
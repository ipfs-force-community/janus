# Janus Backend

The backend service for the Janus project is responsible for synchronizing Filecoin chain data, processing messages, and serving APIs for visualization and analysis.

---

## Features

- **Chain Synchronization**: Syncs Filecoin chain data and processes messages.
- **Database Management**: Stores chain and miner data in a MySQL database.
- **API Services**: Provides RESTful APIs for accessing chain and miner statistics.
- **Indexer**: Periodically indexes chain data for analysis.


---

## Setup Instructions

### Prerequisites

- **Go**: Version 1.20 or higher
- **MySQL**: Ensure a MySQL database is running and accessible
- **Filecoin Node**: A running Filecoin node with API access

### Configuration

1. Copy the example configuration file:
   ```bash
   cp config/config_test.yaml config/config.yaml
   ```

2. Update the `config.yaml` file with your database.

### Installation

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Build the backend services:
   ```bash
   go build -o bin/api cmd/api/main.go
   go build -o bin/indexer cmd/indexer/main.go
   go build -o bin/janus cmd/janus/main.go
   ```

---

## Usage

### API Server

Start the API server:
```bash
./bin/api --config config/config.yaml --port 8080
```

### Indexer

Run the indexer to periodically sync chain data:
```bash
./bin/indexer --config config/config.yaml --interval 10  --node-endpoint 127.0.0.1:1234 --node-token xxxxx
```

### Janus Backend

Run the main backend service:
```bash
./bin/janus --config config/config.yaml --start-epoch 5260000 --end-epoch 5261000
```


---

## API Endpoints

### `/miners`

- **Method**: `GET`
- **Description**: Retrieves daily statistics of new miners.
- **Query Parameters**:
  - `interval`: Number of days to retrieve data for (e.g., `7d`).

---

## Contributing

Contributions are welcome! Please follow the project's coding guidelines and submit pull requests for review.

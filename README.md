# Trading Ace

## Overview

Integrate Uniswap V2 Swap event and give reward to users depending on swap amount in one week

**Trading Ace** is a *simple app demonstrating how to process event from web3 product*. This project aims to integrating
uniswapV2 events and do some customization reward logics. It is built using *golang*.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Testing](#testing)
- [License](#license)
- [Contact](#contact)

## Features

- **Integrate UniSwapV2 Contract**: Integrate UniswapV2 Swap event from pool USDC-WETH by websocket
- **Onboarding/Share Pool Task Support**
    - Onboarding task
        - User will get 100 points when they swap at least 1000 USDC
        - Only once per user
    - Share pool task
        - For user who have completed onboarding task
        - User will get reward points based on the swap amount proportion to the total swap amount in the pool
        - Calculated on a weekly basis
- **Support Realtime Event Processing**
    - Listen to the Swap event from UniswapV2 pool USDC-WETH
    - Use `asynq` to enqueue the event to redis and process it asynchronously
- **Calculate Shared Pool Tasks by Scheduler**
    - Use `go-cron` to schedule the task to calculate the shared pool tasks weekly
- **Query API Support**
    - Get user reward points history
        - path: `GET /api/rewards?user_address=&start_time=&end_time=`
        - query params:
            - user_address: user address `string`
            - start_time: start time of the query period `string` `RFC3339`
            - end_time: end time of the query period `string` `RFC3339`
        -
        example: `GET /api/rewards?user_address=0x1234567890&start_time=2021-09-01T00:00:00Z&end_time=2021-09-30T23:59:59Z`
    - Get tasks of user
        - path: `GET /api/tasks`
        - query params:
            - user_address: user address `string`
            - start_time: start time of the query period `string` `RFC3339`
            - end_time: end time of the query period `string` `RFC3339`

## Installation

### Prerequisites

- Go SDK 1.22
- Docker

### Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/project-name.git
    ```
2. Install dependency
    ```bash
    go mod tidy
    ```
3. follow [Configuration](#configuration), set up environment variables

4. follow [Configuration](#configuration) and `config.example.json` to set up configuration
   file `configuration.{APP_ENV}.json`


1. Run the project development services
    ```bash
    docker-compose up -d
    ```

2. Build and Run the project
    ```bash
    go run src/main.go
    ```

## Configuration

- Supported environment variables
    - **APP_ENV**
        - Environment of the application (development, production, staging), default is `development`
    - **CONFIG_FOLDER**
        - Configuration file name, default is `./config`

- Description of configuration keys in `config.example.json`
- Depend on your **APP_ENV**, the app will load settings from `configuration.{APP_ENV}.json`

```json
{
  "database": {
    // database configuration
    "driver": "postgres",
    // database driver
    "host": "localhost",
    // database host
    "port": 5435,
    // database port
    "username": "postgres",
    // database username
    "password": "postgres",
    // database password
    "dbname": "trading_ace"
    // database name
  },
  "redis": {
    // redis configuration
    "job": {
      // redis job configuration
      "host": "localhost",
      // redis host
      "port": 6377,
      // redis port
      "db": 0
      // redis db
    }
  },
  "ethereum_node": {
    // ethereum node configuration
    "socket": "wss://mainnet.infura.io/ws/v3/socket"
    // ethereum node websocket
  },
  "campaign": {
    // campaign configuration
    "start_time": "2024-09-01",
    // campaign start time
    "weeks": 4
    // campaign weeks
  }
}
```

## Testing

1. create a new test database
     ```bash
     docker run --name test-db -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=trading_ace_test -p 5432:5432 -d postgres
     ```

2. follow [Configuration](#configuration), set up test configuration file `configuration.test.json`

3. Run the test
    ```bash
    sh ./scripts/run_test_coverage.sh
    ```

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details

## Contact

- **Author**: Hui Chih Wang
- **Email**: taya87136@gmail.com


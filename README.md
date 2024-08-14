# Trading Ace

## Overview

Integrate Uniswap V2 Swap event and give reward to users depending on swap amount in one week

**Trading Ace** is a *simple app demonstrating how to process event from web3 product*. This project aims to integrating
uniswapV2 events and do some customization reward logics. It is built using *golang*.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Contributing](#contributing)
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
- **Query API Support**
  - Get user reward points history
    - path: `GET /api/rewards/:userId`
  - Get tasks of user
    - path: `GET /api/tasks/:userId`

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

| Key               | Sub-Key    | Description                          |
|-------------------|------------|--------------------------------------|
| **database**      | driver     | Database driver being used           |
|                   | host       | Database host address                |
|                   | port       | Database port number                 |
|                   | username   | Database username                    |
|                   | password   | Database password                    |
|                   | dbname     | Database name                        |
| **ethereum_node** | socket     | WebSocket endpoint for Ethereum node |
| **campaign**      | start_time | Start date of the campaign           |
|                   | weeks      | Duration of the campaign in weeks    |

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


Welcome

## Prerequisites

Before getting started, ensure you have the following prerequisites installed on your machine:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Setup
0. Config
    copy the .env.example to .env and change the values

1. Clone the repository:
   ```bash
   git clone https://github.com/shandler/go-expert-temperature.git
   ```
   
2. Navigate to the project directory:
    ```bash
    cd go-expert-temperature
    ```
    
3. Test:
   - Test
    ```bash
    make test 
    ```
    -  Coverage
    ```bash
    make cover 
    ```
    
4. upload application with dev image
    ```
    make dev-run
    ```
    
5. upload application with prod image
    ```
    make prod-run
    ```
    
6. deploy the application to GCP, remember to configure the GCP CLI and configure the .env environment variable, copy the .env.example to .env and change the values
    ```
    make deploy
    ```
    
7. Application URL in production
    ```
    https://temperature-prod-4z6n3mffoq-rj.a.run.app/zip-code?zipCode=07987110
    

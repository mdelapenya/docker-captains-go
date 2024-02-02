# Sample golang application with Testcontainers setup

This is a sample application implementing the [TodoBackend](https://todobackend.com/) functionality. 
It uses Testcontainers for setting up local environment for running the application. 

## Running
1. Clone this repository or download and unzip it: 
2. `git clone https://github.com/mdelapenya/docker-captains-go && cd docker-captains-go`
3. Run the application: 
4. `go run -tags dev -v ./...`
5. Open the application in the browser: [link](http://localhost:8080/?http://localhost:8080/todos)
6. Run `docker ps` to see the database container running.
7. Check out [dev_mode.go](https://github.com/mdelapenya/docker-captains-go/blob/main/dev_mode.go#L27) to see how Testcontainers is used to provide a development environment. 

**Optionally**:
1. Install [Testcontainers Desktop](https://testcontainers.com/desktop/) (free)
2. Configure Postgres service in the [Testcontainers Desktop app](https://testcontainers.com/desktop/docs/#debug-testcontainers-based-services)
3. Connect to the database from your IDE / database viewer using. 
4. Open terminal to the service (from the taskbar app menu):
5. Run: `psql -U postgres todos` to connect to the database
6. Execute: `select * from todos;` to see your data.   

## Requirements
* Go 1.21
* Docker environment

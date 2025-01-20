# Meepow

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/AraanBranco/meepow)
![GitHub License](https://img.shields.io/github/license/AraanBranco/meepow)

## Description
Meepow is a powerful tool designed to manage multiple lobbies for Dota 2. It provides an easy-to-use interface for creating, managing, and monitoring lobbies, making it ideal for tournament organizers, community managers, and casual players who want to host their own games. With Meepow, you can automate lobby creation, manage player invites, and ensure smooth operation of your Dota 2 events. The tool is built with Go and leverages Docker for easy deployment and scalability.

## Deploy a Meepow
Follow these steps to deploy the Meepow application:

### 1. Clone the Repository
First, clone the repository to your local machine:

```sh
git clone https://github.com/AraanBranco/meepow.git
cd meepow
```

### 2. Install Dependencies
Install the necessary dependencies:

```sh
go mod download
```

### 3. Configure Environment Variables
Change `config/config.yaml` with the envs to start the API and bots (via docker)

To start the BOT, you need to send the Lobby ID so that when it starts it searches for data in `redis`, just start docker with the following env:
```sh
MEEPOW_LOBBY_ID=12345
```
If you want to pass more variables to docker, send them with the `MEEPOW_` prefix and the application will recognize the environment variable.

### 4. Run the API
Once the application is running, you can access the API at http://localhost:3000 (or the port you specified in the environment variables).

### 5. Run the BOT
When sending a request to `/new-lobby`, a docker will be started along with the Loddy ID and login/password to start the session and lobby.

## Contributing
We welcome contributions! Please refer to our [Contribution Guide]() for more details.

## License
This project is licensed under the MIT license. See the [LICENSE](https://github.com/AraanBranco/meepow/blob/main/LICENSE) file for more info.
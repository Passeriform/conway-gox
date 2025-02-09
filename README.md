# Conway's Game of Life - Multiplayer HTMX Edition

This project is a room-based multiplayer implementation of *Conway's Game of Life* that runs in the browser. It utilizes [HTMX](https://htmx.org/) to handle state updates efficiently and a server to synchronize client states. Users can join rooms to share functionalities such as pausing the simulation and saving/loading states.

Find a running copy here: [https://conway-gox.passeriform.com/](https://conway-gox.passeriform.com/)

## Features

- **HTMX-powered state updates**: Seamless interactions with selective page reloads.
- **Real-time synchronization**: A server ensures all clients in a room stay in sync.
- **Room system**: Clients within the same room share simulation controls.
- **Save & Load states**: Persist simulations for later use.
- **Pause/Resume functionality**: Room-wide simulation control.

## Installation

```sh
# Clone the repository
git clone https://github.com/Passeriform/conway-gox.git

cd conway-gox

# Run the server
go run web/server.go

# Access your game in the browser at http://localhost:8080
```

## Usage

- **Creating a room**: A user can start a new simulation room.
- **Joining a room**: Other users can enter an existing room by simply using the URL
- **Interacting with the simulation**: Users can control the simulation, and save/load states together.

## Technologies Used

- **HTMX** for frontend interactions
- **WebSockets** for real-time updates

## Contributing

Feel free to fork this repository, open issues, or submit pull requests.

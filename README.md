# README

This document outlines the system and library requirements, build instructions, and core functionality for both the CLI and server components of the Cribbage application.

## CLI

### Description

The CLI provides a terminal-based user interface for interacting with the Cribbage game server. It allows users to log in, join or create games, view their hand, play cards, and track the game state.

### System and Library Requirements

-   Go (version specified in `cli/go.mod`)
-   [github.com/charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea): A Go framework for building terminal apps.
-   [github.com/charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss): Style definitions for the TUI.
-   Other dependencies are listed in [cli/go.mod](cli/go.mod) and [cli/go.sum](cli/go.sum).

### Build Instructions

1.  Ensure you have Go installed and configured correctly.
2.  Navigate to the `cli` directory:

    ```sh
    cd cli
    ```
3.  Build the CLI application:

    ```sh
    go build -o crib_cli main.go
    ```

    Alternatively, use the [BuildGame](http://_vscodecontentref_/0) mage target:

    ```sh
    mage BuildGame
    ```

    This command builds the CLI for multiple platforms and places the binaries in the [builds](http://_vscodecontentref_/1) directory.  This requires [magefile.go](http://_vscodecontentref_/2) and [build.go](http://_vscodecontentref_/3) to be configured correctly.
4.  The executable `crib_cli` will be created in the [cli](http://_vscodecontentref_/4) directory, or in the [builds](http://_vscodecontentref_/5) directory if using the mage target.

### Core Functionality

-   **Login:** Authenticates a user with the server.
-   **Lobby:** Displays a list of open matches and allows users to join existing matches or create new ones.
-   **Game View:** Presents the game board, player hands, and other relevant game information.
-   **Input Handling:** Parses user input to perform actions such as playing cards or cutting the deck.
-   **State Management:** Manages the CLI's state and updates the view based on server responses.

## Server

### Description

The server component provides the backend logic and API endpoints for the Cribbage game. It handles user authentication, match management, game state, and communication between clients.

### System and Library Requirements

-   Go (version specified in [go.mod](http://_vscodecontentref_/6))
-   PostgreSQL: Database for storing game data.
-   [github.com/labstack/echo/v4](https://github.com/labstack/echo/v4): High performance, extensible, minimalist Go web framework.
-   [github.com/swaggo/echo-swagger](https://github.com/swaggo/echo-swagger):  Swagger middleware for Echo framework.
-   Other dependencies are listed in [go.mod](http://_vscodecontentref_/7) and [go.sum](http://_vscodecontentref_/8).

### Build Instructions

1.  Ensure you have Go installed and configured correctly.
2.  Navigate to the [server](http://_vscodecontentref_/9) directory:

    ```sh
    cd server
    ```
3.  Build the server application:

    ```sh
    go build -o crib_server main.go
    ```

    Alternatively, use the [BuildServer](http://_vscodecontentref_/10) mage target:

    ```sh
    mage BuildServer
    ```

    This command builds the server and places the binary in the [builds](http://_vscodecontentref_/11) directory. This requires [magefile.go](http://_vscodecontentref_/12) and [build.go](http://_vscodecontentref_/13) to be configured correctly.
4.  The executable `crib_server` will be created in the [server](http://_vscodecontentref_/14) directory, or in the [builds](http://_vscodecontentref_/15) directory if using the mage target.

### Running the Server

1.  Ensure PostgreSQL is installed and running.  Configure the database connection string as needed.
2.  Run the server executable:

    ```sh
    ./crib_server
    ```

### Core Functionality

-   **API Endpoints:** Provides RESTful API endpoints for user authentication, match management, and game actions.  See [router.go](http://_vscodecontentref_/16) and [main.go](http://_vscodecontentref_/17) for endpoint definitions.
-   **Match Management:** Creates, joins, and manages Cribbage matches.
-   **Game Logic:** Implements the rules of Cribbage, including card dealing, playing, and scoring.
-   **Database Interaction:** Stores and retrieves game data from the PostgreSQL database.
-   **Real-time Updates:** (Potentially) Implements real-time communication with clients using WebSockets or similar technologies.
-   **Swagger Documentation:** Generates Swagger documentation for the API endpoints.
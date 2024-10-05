# Minesweeper Game

## Features

- **Classic Gameplay**: Experience the original Minesweeper game with familiar mechanics.
- **Cell Revealing and Flagging**: Floating bubble action chooser for revealing and flagging cells.
- **Charts and Statistics**: View game statistics with charts representing wins, losses, and incomplete games.
- **Game State Management**: Load and save game sessions using UUIDs.
- **Full Server-Side Rendering**: Enjoy SSR and HTMX.

## Technologies Used

- **Go (Golang)**: Server-side logic.
- **SQLC**: SQL compiler for type-safe database interactions.
- **Goose**: Database migration tool.
- **HTMX**: For AJAX requests and HTML swapping.
- **Go-Templates**: For server-side rendered HTML.
- **Go-Echarts**: SSR charts.
- **SQLite**: Database for storing game data.
- **Tailwind CSS**: Styling framework.

## Installation

### Prerequisites

- [**Go**](https://golang.org/dl/): Ensure you have Go installed on your machine.
- [**Goose**](https://github.com/pressly/goose): DB database migration tool.
- [**SQLC**](https://sqlc.dev/): Generating type-safe database code from SQL.

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/Oskarowski/minesweeper.git
   cd minesweeper
   ```
2. **Create the `.env` File**:
   Copy the `.ENV_example` file to `.env` and fill in the required values:

   ```bash
   cp .env.example .env
   ```

3. **Build CSS**:
   Compile the Tailwind CSS files by running:

   ```bash
   npm run build:css
   ```

4. **Apply Database Migrations**:
   Use Goose to apply the database migrations:

   ```bash
   goose -dir db/migrations sqlite3 ./db/minesweeper.db up
   ```

5. **Generate SQL Code**:
   Run SQLC to generate type-safe database code:

   ```bash
   npm run build:sql
   ```

6. **Start the Server**:
   Run the Go server:

   ```bash
   go run main.go
   ```

7. **Access the Game**:
   Open your browser and go to \`http://localhost:8080\` to play the game.

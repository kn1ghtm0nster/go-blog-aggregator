# Blog Aggregator (Gator)

A command-line RSS feed aggregator built in Go that allows users to follow RSS feeds, collect posts, and browse them in the terminal.

## Prerequisites

Before running Gator, you'll need to have the following installed:

### 1. PostgreSQL

- **macOS**: `brew install postgresql`
- **Ubuntu/Debian**: `sudo apt-get install postgresql postgresql-contrib`
- **Windows**: Download from [postgresql.org](https://www.postgresql.org/download/)

### 2. Go (version 1.21 or higher)

- Download and install from [golang.org](https://golang.org/download/)
- Verify installation: `go version`

## Installation

Install the Gator CLI using Go:

```bash
go install github.com/yourusername/blog-aggregator@latest
```

_Note: Replace `yourusername` with your actual GitHub username when you publish this._

## Setup

### 1. Database Setup

Create a PostgreSQL database for Gator:

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE gator;

# Create user (optional)
CREATE USER gator_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE gator TO gator_user;
```

### 2. Configuration File

Create a configuration file at `~/.gatorconfig.json`:

```json
{
  "db_url": "postgres://username:password@localhost/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username` and `password` with your PostgreSQL credentials.

### 3. Run Database Migrations

If building from source, run the database migrations:

```bash
# Install goose for migrations
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/schema postgres "your-connection-string" up
```

## Usage

### User Management

```bash
# Register a new user
gator register <username>

# Login as a user
gator login <username>

# List all users
gator users

# Reset database (deletes all users and data)
gator reset
```

### Feed Management

```bash
# Add a new RSS feed (automatically follows it)
gator addfeed <name> <url>

# List all feeds in the database
gator feeds

# Follow an existing feed
gator follow <url>

# Unfollow a feed
gator unfollow <url>

# List feeds you're following
gator following
```

### Content Aggregation

```bash
# Start the aggregator (runs continuously)
gator agg <duration>

# Example: fetch feeds every 1 minute
gator agg 1m

# Example: fetch feeds every 30 seconds
gator agg 30s
```

### Browse Posts

```bash
# View latest 2 posts (default)
gator browse

# View latest 10 posts
gator browse 10
```

## Example Workflow

1. **Setup and login:**

   ```bash
   gator register alice
   gator login alice
   ```

2. **Add some feeds:**

   ```bash
   gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
   gator addfeed "Hacker News" "https://news.ycombinator.com/rss"
   gator addfeed "Boot.dev Blog" "https://blog.boot.dev/index.xml"
   ```

3. **Start aggregating:**

   ```bash
   # Run in background to collect posts
   gator agg 1m
   ```

4. **Browse collected posts:**
   ```bash
   # In another terminal
   gator browse 5
   ```

## Features

- **Multi-user support**: Multiple users can use the same database
- **RSS feed parsing**: Supports standard RSS 2.0 feeds
- **Continuous aggregation**: Automatically fetches new posts at specified intervals
- **Feed following**: Users can follow/unfollow feeds independently
- **Post browsing**: View collected posts in a clean terminal format
- **Duplicate handling**: Automatically ignores duplicate posts
- **Privacy-focused**: Cascading deletes ensure user data is completely removed

## Safety Notes

- Use reasonable intervals (30s-1m minimum) when running the aggregator to avoid overwhelming RSS servers
- The aggregator runs continuously until stopped with `Ctrl+C`
- Always test with a small number of feeds first

## Development

If you want to build from source:

```bash
git clone https://github.com/yourusername/blog-aggregator.git
cd blog-aggregator
go build -o gator .
```

## License

This project is licensed under the MIT License.

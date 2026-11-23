# Gator
Blog Aggregator to scrape posts off the internet and view them via the CLI.

## Usage

```bash
gator <command>
```

Commands Available:

* `login <name>`: Login to the CLI with an existing user profile.
* `register <name>`: Register a new user profile.
* `reset`: Reset the database.
* `users` : List all registered users.
* `agg <timeBetweenReqs>`: Start the Aggregation loop that fetches feeds.
* `addfeed <name> <URL>`: Add a feed. (URL needs to be a RSS-Feed-URL)
* `feeds`: Show available feeds.
* `follow <URL>`: Follow feeds created by other users.
* `following`: Show registered feeds in the database.
* `unfollow <URL>`: Unfollow a feed.
* `browse <limit>`: Browse aggregated feeds.
## Installation

1. Install using go install:
    ```bash
    go install https://github.com/MMeier05/habit_tracker.git
    ```
2. Install Postgres v15 or later:
    ###Linux/WSL (Debian)
    ```bash
    sudo apt update
    sudo apt install postgresql postgresql-contrib
    ```
    Set a password (Only on Linux/WSL):
    ```bash
    sudo passwd postgres
    ```

    ###macOS with brew:
    ```bash
    brew install postgresql@15
    ```
3. Add a .gatorconfig file in your home directory.
    ```bash
    cd ~ 
    touch .gatorconfig.json
    echo '''
    { 
        "db_url": "postgres://postgres:postgres@localhost:5432/gator?sslmode=disable",
        "current_user_name": "username_set_automatically"
    }
    ''' > .gatorconfig.json
    ```


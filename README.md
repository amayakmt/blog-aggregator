# Gator

A terminal-based RSS feed aggregator. Add feeds, follow them, and browse posts — all from the command line.

## Prerequisites

- [Go](https://go.dev/dl/) 1.21+
- [PostgreSQL](https://www.postgresql.org/download/)

## Installation

```bash
go install github.com/amayakmt/blog-aggregator@latest
```

This compiles and installs the binary to your `$GOPATH/bin`. Make sure that directory is in your `$PATH` — if `gator` isn't found after install, add `export PATH=$PATH:$(go env GOPATH)/bin` to your shell profile.

## Configuration

Gator reads its config from `~/.gatorconfig.json`. Create it manually or let the program create it on first run. It expects two fields:

```json
{
  "db_url": "postgres://username:@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username` with your Postgres username and make sure the `gator` database exists:

```sql
CREATE DATABASE gator;
```

Then run the migrations using [Goose](https://github.com/pressly/goose):

```bash
cd sql/schema
goose postgres "postgres://username:@localhost:5432/gator?sslmode=disable" up
```

## Usage

### User management

```bash
gator register <username>   # create an account and log in
gator login <username>      # switch to an existing user
gator users                 # list all users
```

### Feeds

```bash
gator addfeed "Hacker News" "https://news.ycombinator.com/rss"   # add a feed and follow it
gator feeds                                                        # list all feeds
gator follow <url>                                                 # follow an existing feed
gator unfollow <url>                                               # unfollow a feed
gator following                                                    # list feeds you follow
```

### Aggregation

```bash
gator agg 30s   # start the aggregator, fetching feeds every 30 seconds
gator agg 5m    # every 5 minutes
```

Leave `agg` running in one terminal while using other commands in another. Press `Ctrl+C` to stop.

### Browsing posts

```bash
gator browse        # show 2 most recent posts from feeds you follow
gator browse 20     # show 20 posts
```
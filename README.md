# enable_pg_logs

Tired of manually enabling Postgres query logging? I wrote a binary to turn
it on for you.

### What it does

1) Create the `pg_log` directory

2) Create the `conf.d` directory

3) Write `conf.d/logging.conf` with the following settings:

    logging_collector = on
    log_rotation_size = 200MB
    log_duration = on
    log_lock_waits = on
    log_statement = 'all'

4) Add `include = 'conf.d/logging.conf` to postgresql.conf.

And you're done! (well you might need to restart Postgres, but mostly done)

### Install

```bash
go get -u github.com/kevinburke/enable_pg_logs
```

### Usage

Just run the binary:

```bash
enable_pg_logs
```

### How it works

We find the data directory by running `psql -c 'SHOW data_directory' postgres`.
This may not work if you don't have a database named `postgres`, or if `psql`
is not on the PATH for the running user.

### Errata

- conf.d and pg_log are created using the current user's group and gid. These
may not be the ones that you want to create it with.

- Assumes postgresql.conf is located in the data directory.

- If lines are present after the `include = 'conf.d/logging.conf'` line, they
  may override the settings above.

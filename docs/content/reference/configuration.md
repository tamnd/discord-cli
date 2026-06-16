---
title: "Configuration"
description: "Environment variables, defaults, and the data directory."
weight: 20
---

discord needs almost no configuration: it runs anonymously against public
data out of the box. The settings below let you tune politeness and storage.

## Defaults

| Setting | Default | Flag |
|---|---|---|
| Requests | paced and retried on 429/5xx | `--rate`, `--retries` |
| Per-request timeout | 30s | `--timeout` |
| On-disk cache | under the data directory | `--no-cache` to bypass |

## The data directory

Caches and any record store live under one data directory, chosen in this order:

1. `--data-dir`
2. `DISCORD_DATA_DIR`
3. `$XDG_DATA_HOME/discord`
4. `~/.local/share/discord`

## Environment variables

Every flag has an environment fallback, prefixed `DISCORD_` in
upper case with dashes as underscores. For example:

```bash
export DISCORD_RATE=1s        # same as --rate 1s
export DISCORD_DATA_DIR=~/data/discord
```

Flags win over environment variables, which win over the built-in defaults.

## Sending records to a store

`--db` tees every emitted record into a store as a side effect of reading, so a
session fills a local database without a separate import step:

```bash
discord page <path> --db out.db        # SQLite file
discord page <path> --db 'postgres://...'
```

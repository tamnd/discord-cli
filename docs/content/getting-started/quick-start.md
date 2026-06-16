---
title: "Quick start"
description: "Fetch your first record with discord."
weight: 30
---

Once `discord` is on your `PATH`, fetch a page. The argument is the path
of the page on discord.com (everything after the host), or a full URL:

```bash
discord page <path>
```

By default you get an aligned table. Ask for JSON when you want to pipe it:

```bash
$ discord page <path> -o json
[
  {
    "id": "<path>",
    "url": "https://discord.com/<path>",
    "title": "<path>",
    "body": "..."
  }
]
```

## Shape the output

The same flags work on every command:

```bash
discord page <path> --fields id,url        # keep only these columns
discord page <path> --template '{{.Body}}' # just the body text
discord page <path> -o jsonl | jq .url     # one object per line, into jq
```

`-o` takes `table`, `json`, `jsonl`, `csv`, `tsv`, `url`, or `raw`. Left to
`auto`, it prints a table to a terminal and JSONL into a pipe, so the same
command reads well by hand and parses cleanly downstream. See
[output formats](/reference/output/) for the full contract.

## Follow the links

`links` lists the pages a page links to, and each one is a path you can fetch in
turn:

```bash
discord links <path> -n 10                 # the first ten links
discord links <path> -o url                # just the URLs
discord links <path> -o url | head -3 | xargs -n1 discord page
```

## Serve it instead

The same operations are available over HTTP and to agents over MCP:

```bash
discord serve --addr :7777 &
curl -s 'localhost:7777/v1/page/<path>'          # NDJSON, one record per line
discord mcp                                # MCP over stdio: page, links
```

## What to build next

This scaffold ships one example type, `page`, wired end to end so the whole
chain works today. To make it really about discord, model the records you
care about in `discord/` and declare their operations in
`discord/domain.go`. Each one you add shows up as a command here, a route
under `serve`, and a tool under `mcp`, with no extra wiring. The
[guides](/guides/) cover the common jobs.

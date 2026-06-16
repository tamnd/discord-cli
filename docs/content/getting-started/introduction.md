---
title: "Introduction"
description: "What discord is and how it is put together."
weight: 10
---

Look up Discord servers and read messages from the terminal

discord is a single binary. It speaks to discord over plain HTTPS,
shapes the responses into clean records, and gets out of your way. There is
nothing to sign up for and nothing to run alongside it.

## How it is built

- A **library package** (`discord`) holds the HTTP client and the typed
  data models. It paces requests, sets an honest User-Agent, and retries the
  transient failures any public site throws under load.
- A **domain** (`discord/domain.go`) declares each operation once on the
  [any-cli/kit](https://github.com/tamnd/any-cli) framework. That single
  declaration becomes a CLI command, an HTTP route, an MCP tool, and a
  resource-URI dereference. It is the one place you add to the tool.
- A thin **`cmd/discord`** hands the assembled app to `kit.Run`, which
  builds the command tree and the serve and mcp surfaces.

## One operation, four surfaces

Because an operation is surface-neutral, the same `page` you run on the command
line is also a route and a tool:

```bash
discord page <path>                  # the command
discord serve --addr :7777           # GET /v1/page/<path>
discord mcp                          # the page tool, over stdio
ant get discord://page/<path>        # the URI dereference (via a host)
```

You write the fetch and the record shape; the surfaces come for free.

## Scope

discord is a read-only client over data discord already serves
publicly. It reads that data and shapes it for you. That narrow scope keeps it a
single small binary with no database, no daemon, and no setup.

Next: [install it](/getting-started/installation/), then take the
[quick start](/getting-started/quick-start/).

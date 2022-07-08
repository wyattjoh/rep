# rep

![Tests](https://github.com/wyattjoh/rep/workflows/Test/badge.svg)
[![Go Report](https://goreportcard.com/badge/github.com/wyattjoh/rep)](https://goreportcard.com/report/github.com/wyattjoh/rep)
[![GitHub release](https://img.shields.io/github/release/wyattjoh/rep.svg)](https://github.com/wyattjoh/rep/releases/latest)

The [rep](https://github.com/wyattjoh/rep) tool is designed to assist with
working with reproductions of issues for open source maintainers. Currently this
will:

1. Prompt for the issue details including GitHub reproduction repository
2. Clones to a specified folder using a known format `${GitHub Issue}_${Camel Case Issue Description}`
3. Installs dependancies using `pnpm` (future to use configurable tool)

On first run it prompts for some configuration details, saving them to `~/reprc.json`.

## Installation

You can use the standard Go utility to get the binary and compile it yourself:

```bash
go get -u github.com/wyattjoh/rep/...
```

From Homebrew:

```bash
brew install wyattjoh/stable/rep
```

Or by visiting the [Releases](https://github.com/wyattjoh/rep/releases/latest)
page to download a pre-compiled binary for your arch/os.

## License

MIT

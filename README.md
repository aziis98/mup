# μploader - Go Micro Uploader

A simple file uploader written in Golang. It only uses [Go Chi](https://github.com/go-chi/chi) and [pflag](https://github.com/spf13/pflag) as dependencies.

## Installation

### Git

```bash
$ git clone https://github.com/aziis98/mup
$ cd mup
$ go build -v -o bin/mup .
```

### Go Install

```bash
$ go install github.com/aziis98/mup
$ mup
```

### [TODO] Static Binary from Release

I have set up an action that continuously builds the binary, here is an install script for Linux that installs the latest version to `~/.local/bin/mup`

```bash
$ wget -qO- https://raw.githubusercontent.com/aziis98/mup/main/install.sh | sh
$ mup
```

## Usage

```bash
# Start the server on port 5000 in the "Uploads" directory
$ mup
```


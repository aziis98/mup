# μploader - Go Micro Uploader

<p align="center">
<img src="https://github.com/user-attachments/assets/268d853f-9b69-4fa1-853e-e645818c3f6d" alt="screenshot" />
</p>


A simple file uploader that can be used to easily move and share files across the local network between devices with a web browser. 

It only uses [Go Chi](https://github.com/go-chi/chi) and [pflag](https://github.com/spf13/pflag) as dependencies and the releases provide a statically linked binary for Linux.

## Installation

### Static Binary from Release

Run the following command to install the latest version to `~/.local/bin/mup`

```bash
wget -qO- https://raw.githubusercontent.com/aziis98/mup/main/install | sh
```

Then you can run `mup` from anywhere in your terminal, the default upload directory is `Uploads` so this can even be run directly from the home folder (only the files inside `Uploads` are served to the client).

### Git

```bash
$ git clone https://github.com/aziis98/mup
$ cd mup

# Run the server
$ go run -v .

# Build the binary
$ go build -v -o bin/mup .
```

## Usage

```bash
# Start the server on port 5000 in the "Uploads" directory
$ mup
```


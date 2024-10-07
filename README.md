# μploader - Micro Uploader

[![Build and Release](https://github.com/aziis98/mup/actions/workflows/release.yml/badge.svg)](https://github.com/aziis98/mup/actions/workflows/release.yml)

<p align="center">
<img width="800px" src="https://github.com/user-attachments/assets/4501778a-109e-4478-95a6-e617b7fc3160" alt="screenshot" />
</p>

A simple file uploader that can be used to easily move and share files across the local network between devices with a web browser. 

It only uses [Go Chi](https://github.com/go-chi/chi) and [pflag](https://github.com/spf13/pflag) as dependencies and the releases provide a statically linked binary for Linux.

**Motivation.** Sometimes I want to move files between my pc and a device I do not own that has an old browser version (that generally means expired https certificates, oh and without any cables). When I try to search for a tool like this I always find random outdated projects that aren't easy to setup. So I made this tool that can be easily installed on all linux systems.

## Demo

I have setup a demo hosted on glitch.me that clears the uploads every couple of minutes. This is just to preview the UI, use at your own risk:
 
<p align="center">
<a href="https://checkered-pyrite-garage.glitch.me">https://checkered-pyrite-garage.glitch.me</a>
</p>

## Installation

### Static Binary from Release

Run the following command to install the latest version to `~/.local/bin/mup`

```bash
curl -sSL https://raw.githubusercontent.com/aziis98/mup/main/install | sh
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

### Dockerfile

I provide this just to easily deploy on a local server. I **do not recomend to expose this publicly** on the web as there is no auth or password and there is no upload limit to the number of files and all files in the `Uploads/` folder are public by default for now.

```bash shell
$ docker build -t mup .
$ docker run -p 5000:5000 -v $PWD/Uploads:/Uploads mup

# On a LAN go to "<ip from below>:5000"
$ ip addr | grep 'inet 192'
```

## Usage

```bash
$ mup --help
Usage:
  mup [OPTIONS] [UPLOAD_FOLDER]

A micro file uploader, the default upload folder is 'Uploads'

Options:
  -h, --host string           Host to run the server on (default "0.0.0.0")
  -s, --max-upload-size int   Maximum upload size in MB (default 100)
  -p, --port int              Port to run the server on (default 5000)
```

## To Do

- [ ] Decrease minimum requirements (make this work without js, now needed for the upload with live progress bar)

- [ ] Update Github Action to make a release for the Raspberry Pi (`aarch64` and `armv7l`)

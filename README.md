# My personal page

## Introduction

Welcome to my personal page!

This project showcases a comprehensive web application built using the latest
tools and technologies in the industry: GoFiber, HTMX, Tailwind, DaisyUI, and
more.

## How to run this project

To run this project, you need to have Go installed on your machine. You can find the installation instructions [here](https://golang.org/doc/install).

Once you have Go installed, you can run the following commands to start the project:

```bash
## Usage

make            # install dependencies and run
make install    # install all dependencies
make run        # run the server
make deploy     # deploy the server
make release    # build the server version
make fmt        # format the code
make t          # run the tests
make encrypt    # encrypt the .env file
make decrypt    # decrypt the .env.gpg file
```

## Tech Stack

-   🚀 HTMX: Powering our modern frontend

-   📊 Golang with Fiber: Powering our backend services

-   🐳 Docker for distributing

-   ☁️ Hosted by [fly.io](https://fly.io)

### HTMX

I've utilized HTMX to build a reusable and modular frontend that makes it easy
to create and maintain complex and highly-dynamic user interfaces. With HTMX, I
was able to streamline my development process and deliver a fast, responsive,
and highly interactive user experience.

### Go Templates

The project also leverages Go Templates to render dynamic HTML content. By using
templates, it can easily generate HTML content based on data from our backend
services, ensuring consistency across our application's UI.

## Technical Details

Here are some of the technical details behind our FullStack Go project:

### Project Strcuture

```sh
cmd             # entry point of the project (server, cli, etc)
  server      # server entry point

internal
  config      # app configs
  constants   # app constants
  cron        # cron job module
  github      # github service
  lua         # lua runtime
  models      # data models (gorm)
  server      # create our backend service

```

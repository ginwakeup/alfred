# Alfred
Alfred is a development environment manager for local applications, based on Docker-Compose and written in Go.

The idea behind Alfred it to have:
1. A location (repository) for docker-compose files
2. A CLI (Alfred) that given a config file: `alfred.yaml`, can understand your project dependencies, and spin up the right services using the repository docker-compose files, all in the same network.

For example, let's say you are starting a FastAPI Project in Python, and you need a MongoDB and Redis, you run: `alfred init`, provide the right information and get the following in your project root:

> alfred.yaml
```
project:
  name: hello-world
  compose: ./docker-compose.yaml

dependencies_root: /Users/Iacopo/GolandProjects/Alfred/examples/simple/systems
dependencies:
  - mongo

network:
  name: alfred-dev
```

Then you can run the development environment by executing: `alfred dev`

## Get Started
1. Place `alfred` in your `PATH`

To create a new project, run:
> alfred init --project-name <name> --project-root <project-root> --dependencies <comma,sep,deps>

e.g.
> alfred init --project-name my-project --project-root /my-project/root --dependencies mongo,redis

## Features
1. Define development environments with an `alfred.yaml` and define a full development environment for your application.

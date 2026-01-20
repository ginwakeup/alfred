# Alfred
Alfred is a development environment manager for local applications, based on docker-compose.

The idea behind Alfred is to have:
1. A location (repository) for docker-compose files. This can be either on git, or a local file-system.
2. A CLI (Alfred) that given a config file: `alfred.yaml`, can understand your project dependencies, and spin up the right services using the repository docker-compose files, all in the same network.

For example, let's say you are starting a FastAPI Project in Python, and you need a MongoDB and Redis, you run: `alfred init`
and provide the necessary info for your project, e.g. that you need `mongo`.

The following gets generated:

> /my-project-dir/alfred.yaml
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

Then you can run it by executing: `alfred run /my-project-dir`

Alfred will:
1. Collect and cache all your dependencies docker-compose files from either git or the file-system.
2. Override the compose files to ensure services are in the same network and can communicate.
3. Execute your dependencies and project compose.

## Why?
The reason behind Alfred is to provide some sort of layer to retrieve docker-compose files for backend systems from
repositories/registries, to ensure the developer can get up to speed with a project without having to create a new compose for each project.

Alfred also makes it possible to have a centralized repository where all docker-compose files can be stored and shared
with a team of developers, so versions can be kept consistent for simple development environments.

## Get Started
1. Place `alfred` executable in your `PATH`

To create a new project using a git repo, run:
> alfred init --project-root <project-root> --project-name <project-name> --dependencies-location https://github.com/repository.git --dependencies-repo-type git --dependencies mongo,kafka

Or, you can also point to a file system location where your docker-compose files live:

> alfred init --project-root <project-root> --project-name <project-name> --dependencies-location /my/compose/files/root --dependencies-repo-type local --dependencies mongo,kafka

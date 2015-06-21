# VoicePipe [![Build Status](https://travis-ci.org/y-taka-23/voicepipe.svg)](https://travis-ci.org/y-taka-23/voicepipe)

Build parameterized Docker images from a single Dockerfile.

## Description

As you know, you often need to manage multiple Docker images,
depending on versions of applications or configurations for environments
(e.g. development, staging, production etc.)

The `docker build` command, however, has no options to inject such parameters.
Thus you must manage multiple Dockerfiles separately, and of course they bother you.

VoicePipe generates and builds multiple Dockerfiles from a single Dockerfile
by overriding `ENV` instructions.
So you can define parameters as environment variables in your Dockerfile.

## Requirement

* [Docker](https://www.docker.com/)

## Installation

```
$ go get https://github.com/y-taka-23/voicepipe
```

## Example Usage

Let us build multiple [Nginx](https://registry.hub.docker.com/_/nginx/) images,
each of them has different `index.html` files, from a single Dockerfile by VoicePipe.

Assume the following directory structure.
You will find it under the `example` in this repository.

```
+--- Dockerfile
+--- index
|    +--- index_develop.html
|    +--- index_latest.html
|    `--- index_production.html
`--- voicepipe.yml
```

First, define your application with a `Dockerfile`:

```
$ vi Dockerfile
FROM nginx
ENV INDEX_HTML index_latest.html
COPY index/${INDEX_HTML} /usr/share/nginx/html/index.html
```

Next, define the parameters for each images in `voicepipe.yml`:

```
$ vi voicepipe.yml
repository: user/nginx
images:
  - tag: develop
    description: "for the development environment"
    parameters:
      - name: INDEX_HTML
        value: index_develop.html
  - tag: production
    description: "for the production environment"
    parameters:
      - name: INDEX_HTML
        value: index_production.html
```

Finally, run VoicePipe in the `example` directory.

```
$ voicepipe build
```
VoicePipe will build the Docker images
`user/nginx:develop` and `user/nginx:production`
with the corresponding `index_*.html` files.

For more detail, see `voicepipe help`.

## License

[MIT](https://github.com/y-taka-23/voicepipe/blob/master/LICENSE)

## Author

[y-taka-23](https://github.com/y-taka-23)


VoicePipe
=========

VoicePipe builds parameterized Docker images from a single Dockerfile.

[![Build Status](https://travis-ci.org/y-taka-23/voicepipe.svg)](https://travis-ci.org/y-taka-23/voicepipe)

## Example Usage

Assume the following directory structure:

```
+--- Dockerfile
+--- conf
|    +--- develop.conf
|    `--- production.conf
`--- voicepipe.yml
```

First, you define your application with a `Dockerfile`:

```
FROM nginx
ENV NGINX_CONF dummyvalue
COPY conf/${NGINX_CONF} /etc/nginx/nginx.conf
```

Next, you define the parameters for each images in `voicepipe.yml`:

```
repository: user/nginx
images:
  - tag: develop
    parameters:
      - name: NGINX_CONF
        value: develop.conf
  - tag: production
    parameters:
      - name: NGINX_CONF
        value: procudtion.conf
```

Finally, run `voicepipe` in the directory.
VoicePipe will build the Docker images
`user/nginx:develop` and `user/nginx:production`
with the corresponding `*.conf` files.

## Requirement

* [Docker](https://www.docker.com/)

## Installation

```
go get https://github.com/y-taka-23/voicepipe
```

## License

[MIT](https://github.com/y-taka-23/voicepipe/blob/master/LICENSE)

## Author

[y-taka-23](https://github.com/y-taka-23)


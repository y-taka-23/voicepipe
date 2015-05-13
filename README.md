VoicePipe
=========

VoicePipe builds parameterized Docker images from a single Dockerfile.

[![Build Status](https://travis-ci.org/y-taka-23/voicepipe.svg)](https://travis-ci.org/y-taka-23/voicepipe)

## Example Usage

Assume the following directory structure:

```
+--- Dockerfile
+--- index
|    +--- index_develop.html
|    `--- index_production.html
`--- voicepipe.yml
```

First, you define your application with a `Dockerfile`:

```
FROM nginx
ENV INDEX_HTML dummyvalue
COPY index/${INDEX_HTML} /usr/share/nginx/html/index.html
```

Next, you define the parameters for each images in `voicepipe.yml`:

```
repository: user/nginx
images:
  - tag: develop
    parameters:
      - name: INDEX_HTML
        value: index_develop.html
  - tag: production
    parameters:
      - name: INDEX_HTML
        value: index_production.html
```

Finally, run `voicepipe` in the directory.
VoicePipe will build the Docker images
`user/nginx:develop` and `user/nginx:production`
with the corresponding `index_*.html` files.

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


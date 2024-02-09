# Go Products API

## Where is the documentation?

The documentation can be found on:

``http://localhost:8000/api/v1/docs/index.html``

## How to generate documentation?

It is necessary to install swag by using the command:

``go install github.com/swaggo/swag/cmd/swag@latest``

Then:

``swag init -g cmd/server/main.go``
# Overview

This repo is used as an example of how Tilt can be used to manage tasks within a microservice monorepo. This is roughly styled to be similar to how we at SOON_ use Tilt.

This is the companion repo to our [Tilt At SOON_](TODO) blog post.

## Structure

There are two microservices (if you can even call them that):

- `image-service/` - Provides `/totallyRandomImage` which gives you a "totally random" image.
- `name-service/` - Provides `/name` which gives you a random name from a list of names, as well as allows you to add names to the list.

Beyond that, there's a `docker-compose.yaml` file to provide a datastore-emulator container, as well as a UI for it.

Finally there is of course the `Tiltfile` which sets up the ever so wonderful Tilt for us.

## Using this repo

Currently this example repo hasn't been made compatible with Windows (just use WSL), but should work on Linux and Mac.

Please install the following:

* [Tilt](tilt.dev)
* [Docker](https://docs.docker.com/get-docker/)
* [Go](https://go.dev/dl/)

Then simply run `tilt up && tilt down` within your terminal. You can press Space within your terminal to open the Tilt UI in your browser.

## Tilt

Tilt is configured to do the following things:

* Execute the `docker-compose.yaml` file. You can see each container individually in the "Unlabeled" section of the Tilt UI.
  * If you examine the `datastore-ui` resource within Tilt, you'll see Tilt automatically displays a link above the log window so you don't need to keep the README up to date with the link - devs can discover it naturally in Tilt!

* Provides a generic way to define Go microservices. The Tiltfile has more comments on how this works.

* For each service:

  * A `:test` resource is created, which will run `go test ./...` anytime the source code or go mod files change.
    * This resource has a `Show Cover Profile` button to easily access the go coverage visualiser.
  * A `:lint` resource, which will run [golangci-lint](https://github.com/golangci/golangci-lint) against the service when the source code changes.
  * A `:serve` resource, which will run your service locally, and restart iw when the source code changes.
    * If configured with test requests, buttons will appear above the log window to easily perform them. Have a look at the name-service for example!
    * Additionally a link to the main API endpoint will be added above the log window for easy discovery.

## What should I do with this repo?

Play around with it! Launch Tilt; mess around with the code and see how Tilt automatically performs tasks when you do so; see if this pre-configured Tilt workflow feels good to use; and maybe even try adding a new microservice or two into the setup.

e.g.

1. `tilt up && tilt down`.

2. Press Space to open the UI.

3. Open the `name-service:serve` resource.

4. Use the `Add Name` button to add a few names.

5. Press the link above the log window, and refresh a few times to see your new names show up.

6. Open the `datastore-ui` resource in Tilt.

7. Press the link above the log window, and see that its storing the names for `name-service`.

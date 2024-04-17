#### Extensions and misc configuration ####

load("ext://uibutton", "cmd_button", "text_input")

#### Constants ####

golangci_lint_version = "v1.57-alpine"
projects = { # Keys are used for naming resources
    "image-service": {
        "directory": "image-service",   # Which directory to serve from
        "cmd": "go run ./cmd/main.go",  # What command to run to serve the application
        "api": "/totallyRandomImage",   # What the main API endpoint is
        "env": {                        # Additional env vars to pass
            "HOST": "localhost:10000"
        },
        "test_requests": {}
    },

    "name-service": {
        "directory": "name-service",
        "cmd": "go run ./cmd/main.go",
        "api": "/name",
        "env": {
            "HOST": "localhost:10001",

            "DATASTORE_PROJECT_ID": "tilt-demo",
            "DATASTORE_EMULATOR_HOST": "localhost:10080",
            "DATASTORE_HOST": "http://localhost:10080",
            "DATASTORE_EMULATOR_HOST_PATH": "localhost:10080/datastore",
        },

        "test_requests": {          # Adds a button for each request
            "Add Name": {           # Keys are used as the button name
                "method": "POST",
                "fields": ["name"]  # Each field is given a text_input, and stuffed into JSON
            }
        }
    },
}

#### Functions ####

# Describes a generic Go service, with auto reload; testing, and linting.
def use_go_service(name, dir, api, cmd, env, requests):
    # Calculate some useful constants
    watch = [
        "{}/go.mod".format(dir),
        "{}/go.sum".format(dir),
        "{}/cmd/".format(dir),
        "{}/internal/".format(dir),
    ]

    serve_resource = "{}:serve".format(name)
    test_resource = "{}:test".format(name)
    lint_resource = "{}:lint".format(name)

    # Register main resources
    
    local_resource(
        serve_resource,
        serve_cmd=cmd,
        serve_env=env,
        serve_dir=dir,
        deps=watch,
        labels=[name],
        allow_parallel=True,
        links=[
            link("http://{}{}".format(env["HOST"], api))
        ]
    )

    local_resource(
        test_resource,
        cmd="go test ./... -cover -coverprofile=cover.out",
        dir=dir,
        deps=watch,
        labels=[name],
        allow_parallel=True
    )
    cmd_button(
        "{}:cover".format(name),
        dir=dir,
        resource=test_resource,
        icon_name="open_in_new", # https://fonts.google.com/icons?selected=Material+Symbols+Outlined:open_in_new:FILL@0;wght@400;GRAD@0;opsz@48
        text="Show Cover Profile",
        argv=["go", "tool", "cover", "-html", "cover.out"]
    )

    local_resource(
        lint_resource,
        cmd="docker run -v .:/app -w /app --rm golangci/golangci-lint golangci-lint run --color always",
        dir=dir,
        deps=watch,
        labels=[name],
        allow_parallel=True
    )

    # Create test requests buttons
    for requestName in requests:
        request = requests[requestName]
        inputs = []
        payload = "{" # JSON payload in the form: { "input1: "$input1", "input2": "$input2", etc. }
        for fieldName in request["fields"]:
            inputs.append(text_input(fieldName))
            payload += '"{0}": "${0}",'.format(fieldName)
        payload = payload[:-1] + "}"

        cmd_button(
            "{}:request:{}".format(name, requestName),
            resource=serve_resource,
            text=requestName,
            argv=[
                "bash",
                "-c",
                """
                curl \\
                    -s \\
                    -X {method} \\
                    -H "Content-Type=application/json" \\
                    -d "{payload}" \\
                    http://{host}{api}
                """.format(
                    method=request["method"], 
                    payload=payload.replace('"', '\\"'), 
                    host=env["HOST"], 
                    api=api
                )
            ],
            inputs=inputs
        )

#### Main ####

# Register each project.
for name in projects:
    config = projects[name]
    use_go_service(name, config["directory"], config["api"], config["cmd"], config["env"], config["test_requests"])

docker_compose("docker-compose.yaml")

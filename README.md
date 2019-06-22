# Infuse - the global template parser

Infuse is a CLI tool and a go library that can be imported in other Go projects. Its main purpose is to handle template parsing for different template formats in a generic way.

## Getting started

If interested in downloading the code and using it as a dependency in another Go project, simply do

```bash
go get github.com/jucardi/infuse
```

## The command line tool

The command line tool is a convenient terminal command that allows to parse a template, passing the a JSON or YAML representation to be used as values within the template.


### How to install it

There are a few options to install the CLI tool.

#### Download the release file

##### For UNIX (Mac and Linux)

1) Download the binary

    ```bash
    sudo curl -L https://github.com/jucardi/infuse/releases/download/v1.0.0.0/infuse-$(uname -s)-$(uname -m) -o /usr/local/bin/infuse
    ```

2) Apply executable permissions

    ```bash
    sudo chmod +x /usr/local/bin/infuse
    ```

3) Test the installation

    ```bash
    infuse --help
    ```

#### Install using Go

```bash
go get -u github.com/jucardi/infuse/cmd/infuse
```

#### Installing the CLI from the source code directory

Using `make` run

```bash
make install
```

This `make` recipe will compile the binary for your local architecture and place it under `$GOPATH/bin` via `go install`

#### Building the CLI binaries from code

Using `make` run

```bash
make compile-all
```

This `make` recipe will build the binaries for Linux, Mac and Windows under the `build` directory.

### Uninstallation

#### If installed using the first method (downloading the release)

```bash
sudo rm /usr/local/bin/infuse
```

#### If installed using Go

```bash
rm $GOPATH/bin/infuse
```

### Usage

```bash
infuse [flags] [template file]
```

#### Flags

All the flags are optional, including the data input. Infuse supports using a template to parse declared environment variables, so they can also be used as data to be parsed in.

##### Data Input flags

The input flags indicate how the data to be parsed into the templates is read. Note that only one input type allowed (-f, -s, -u). May change in future releases to allow merging multiple data sources.

- **`-f` or `--file`:** *A JSON or YAML file to use as an input for the data to be parsed*
- **`-u` or `--url`:** *A URL for HTTP GET to a JSON or a YAML file. Useful to parse data from config servers*
- **`-s` or `--string`:** *A JSON or YAML string representation*

##### Target flags

The target flags indicate where the parsed template will be output

- **`-o` or `--output`:** *Indicate an output file. If not specified, the resulting template will be printed to StdOut*

##### Template definitions flags

The template definition flags allow auxiliary template files to be loaded so they can be used in the primary template.

- **`-d` or `--definition`:** *File path of another template to be imported and used by the primary template to be parsed. This flag can be used multiple times to load multiple template definitions*
- **`-p` or `--pattern`:** *Search pattern to load multiple template definitions, for example `-p ./templates/*`*

### Examples

```bash
infuse -f service-config.yml -o docker-compose.yml -d global/mongo.tmpl -d global/redis.tmpl docker-compose.tmpl
```

```bash
infuse -u http://config-server.local/something/service-config.json -o docker-compose.yml -p global/* docker-compose.tmpl
```

```bash
infuse -o something.config something.tmpl
```

#### When are template definitions useful?

When defining certain global templates that can be reused in other primary templates.

For example, a `docker-compose.yml` that defines a service stack where one of the containers is a database definition. The database definition can be a separate template and be imported in the `docker-compose.yml` file when generating the service stack definition.

Given the following files:

**service-config.yml** - *The configuration for a specific service*

```yaml
name: some-service
version: RC-1.0.0
port: 1234
develop_mode: true
resources:
  replicas: 3
  cpus: '0.5'
  memory: 512M
db:
  port: 4321
```

**service.tmpl** - *Primary template that defines a service stack. Generic and can be used for multiple microservices, provided the right configuration*

```yaml
version: "3.3"
services:
  api:
    image: registry.local/services/{{ .name }}:{{ .version }}
    ports:
      - '{{ .port }}:{{ .port }}'
    healthcheck:
      test: curl -f http://localhost:{{ .port }}{{ .health_uri | default "/health" }} || echo 1
      interval: 5s
      timeout: 20s
      retries: 3
    networks:
      - {{ .name }}-net
    deploy:
      replicas: {{ .resources.replicas | default 2 }}
      placement:
        constraints:
          - node.role == worker
          - node.labels.workertype == services
      resources:
        limits:
          cpus: {{ .resources.cpus }}
          memory: {{ .resources.memory }}
      restart_policy:
        condition: any
{{ if .db }}
{{ template "mongo.tmpl" . }}
{{ end }}

networks:
  {{ .name }}-net:
```

**mongo.tmpl** - *Template used as definition to be imported. It defines a generic mongo database config which can be reused in other templates*

```yaml
  mongo:
    image: mongo:latest
{{ if .develop_mode }}
    ports:
      - {{ .db.port }}:27017
{{ end }}
    volumes:
      - {{ .name }}-mongodb:/data/db
    networks:
      - {{ .name }}-net
    deploy:
      replicas: 1
      placement:
        constraints:
          - node.role == worker
          - node.labels.workertype == dbe
      restart_policy:
        condition: any

volumes:
  {{ .name }}-mongodb:
    driver: local-persist
    driver_opts:
      mountpoint: /mnt/db-data/{{ .name }}-mongodb

```

## The template library

### Custom helpers

- `default`: Indicates a default value if the input data does not contain a specific value.

    Usage:
    ```html
    {{ .some_value | default "something" }}
    ```
    If `some_value` is not defined, the value `something` will be used.

- `map`: Helps define a `map` object which can be passed as the data object for another template definition

    Usage:
    ```html
    {{- $someObj := map "X" "chuck" "Y" "norris" "Z" 1234 "something" .something -}}
    ```
    In this case, a new dictionary is being created and set to an object called `$someObj`. The number of members in the map definition must be an even number, since the definition will be done as `key`, `value`, `key`, `value`...

    For this example, the equivalent JSON object will be
    ```json
    {
        "X": "chuck",
        "Y": "norris",
        "Z": 1234,
        "something": "[the value for the 'something' key in the original data input]"
    }
    ```

- `dict`: Is an alias of `map`, stands for "dictionary"

- `template`: Indicates a loaded template definition.

    Usage:
    ```html
    {{ template "template_name" $dataObj }}
    ```
    If using the CLI, the template name will always match the file name of the template definition. `$dataObj` is the object that contains the data that will be used in the sub-template. If using the original data object, simply pass `.` instead of `$dataObj`, or if the data object is a sub-key of the original input, pass `.sub_key` instead

- `env`: Reads an environment value and parses it into the template

    Usage:
    ```
    {{ env "SOME_ENVIRONMENT_VARIABLE" }}
    ```
    If `SOME_ENVITONMENT_VARIABLE=something` the result of the example above will be `something`

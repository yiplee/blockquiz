## BlockQuiz

[Usage](./docs/usage.md)

[API Doc](./docs/api.md)

### Build

```shell script
make build
```

### Config

请参考 config.template.yaml

### Create Database Tables

```shell script
./blockquiz migrate
```

### Run API Server

```shell script
./blockquiz --debug --config your.config.file.path http --port 8080
````

### Run Engine

```shell script
./blockquiz --debug --config your.config.file.path run
````

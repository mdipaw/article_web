## Get started

1. Download golang
2. Create $GOPATH
3. Download Docker dekstop

```shell
mkdir -p $HOME/go export GOPATH="$HOME/go"

# Folder contains your golang source codes
mkdir -p $GOPATH/src

# Folder contains the binaries when you install an go based executable
mkdir -p $GOPATH/bin

# Folder contains the Go packages you install
mkdir -p $GOPATH/pkg
```

## Cloning project

```shell
cd $GOPATH/src
git clone https://github.com/difaagh/article_web.git
```

## Download Dependencies

```shell
cd $GOPATH/src/article_web

go mod download
```

## Setup a shell script to execute inline env

Most of the configs are stored in the standard OS environment variables. You can set the values of the environment
variables with whichever method you prefer.

For convenience of switching between different config values, you can add the following bash function to your `.bashrc`
or `.zshrc`, then you can prefix the commands you're running with `with_env file.env <your command here>` to run that
command with the environment variables set in the `file.env`.

```shell
with_env () {
    eval "$(grep -vE '^#' "$1" | xargs)" "${@:2}"
}
```

For example, if you have a `local.env` file containing the following:

```
DATABASE_HOST=localhost
DATABASE_USER=pg
```

, and your run the command `with_env local.env go run ./bin/web/main.go`, it would be equivalent
to `DATABASE_HOST=localhost DATABASE_USER=pg go run ./bin/web/main.go`, which will start the executable with those following
environment variable values.

You can maintain different config files as needed. e.g. you can keep a `local.env` file to run the local server
with `with_env local.env go run main.go`, together with a `test.env` that you can use with `with_env test.env go test`
to run the tests with a different database.


## Test
```shell
cd $GOPATH/src/article_web
docker-compose -f docker-compose-test.yml up
with_env env.example go test ./... > debug.test
docker-compose -f docker-compose-test.yml down
```

## Start project
```shell
cd $GOPATH/src/article_web
docker-compose -f docker-compose.yml up
```


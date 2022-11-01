# SimplePasswords - Vaults

Password vault microservice.

## System Requirements
- Latest version of Go

## Download

Clone this repository: `git clone https://github.com/liobrdev/simplepasswords_vaults.git`

## Required Environment Variables

Each of the environment variables in the following table **must** be an **absolute path** to an existent `UTF-8`-encoded text file. The contents of the first line of this text file will be parsed to a Go data type (empty file contents will be parsed with a default value). This data will then be saved to the AppConfig field whose name matches that of the corresponding environment variable.

| **Required Environment Variable** | **Description of File Contents** | **Data Type** | **Default Value** |
|-----------------------------------|----------------------------------|---------------|-------------------|
| GO_FIBER_ENVIRONMENT | Should be either `development`, `production`, or `testing`. | `string` | `"development"` |
| GO_FIBER_BEHIND_PROXY | Configures Fiber server setting `EnableTrustedProxyCheck bool`. Should be either `true` or `false`. | `bool` | `false` |
| GO_FIBER_PROXY_IP_ADDRESSES | Configures Fiber server setting `TrustedProxies []string`. Should be comma-separated IP address string(s). | `[]string` | `[]string{""}` |
| GO_FIBER_VAULTS_DB_USER | PostgreSQL database user. | `string` | `""` |
| GO_FIBER_VAULTS_DB_PASSWORD | PostgreSQL database password. | `string` | `""` |
| GO_FIBER_VAULTS_DB_HOST | PostgreSQL database host. | `string` | `""` |
| GO_FIBER_VAULTS_DB_PORT | PostgreSQL database port. | `string` | `""` |
| GO_FIBER_VAULTS_DB_NAME | PostgreSQL database name. | `string` | `""` |
| GO_FIBER_REDIS_PASSWORD | Redis cache password. | `string` | `""` |
| GO_FIBER_SECRET_KEY | Secret key for various app-level encryption methods. | `string` | `""` |
| GO_FIBER_SERVER_HOST | Fiber app will be run from this host. | `string` | `"localhost"` |
| GO_FIBER_SERVER_PORT | Fiber app will be run from host using this port. | `string` | `"8080"` |

### Methods For Setting Environment Variables

Required environment variables may be sourced by:

1. setting variables via the command line at compile time,

and/or,

2. setting variables in the shell environment before compile time,

and/or,

3. including a `.env` file in the root application folder before compile time.

Each environment variable **must** be set using at least one of these three methods. Variables set in the shell environment will override duplicate variables included in a `.env` file, and variables set via the command line will override duplicate variables set in the shell environment and/or included in a `.env` file. That is, environment variable sources have the following precedence: command line, *then* shell, *then* `.env` file.

#### An example using all three methods:

With the following `.env` file present in the root application folder:

```bash
# .env

GO_FIBER_ENVIRONMENT=/path/to/secret_files/environment_3
GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_3
GO_FIBER_PROXY_IP_ADDRESSES=/path/to/secret_files/proxy_ip_addresses
GO_FIBER_VAULTS_DB_USER=/path/to/secret_files/db_user
GO_FIBER_VAULTS_DB_PASSWORD=/path/to/secret_files/db_password
GO_FIBER_VAULTS_DB_HOST=/path/to/secret_files/db_host
GO_FIBER_VAULTS_DB_PORT=/path/to/secret_files/db_port
GO_FIBER_VAULTS_DB_NAME=/path/to/secret_files/db_name
GO_FIBER_REDIS_PASSWORD=/path/to/secret_files/redis_password
GO_FIBER_SECRET_KEY=/path/to/secret_files/secret_key
GO_FIBER_SERVER_HOST=/path/to/secret_files/server_host
GO_FIBER_SERVER_PORT=/path/to/secret_files/server_port
```

Then running the following commands:

```bash
export GO_FIBER_ENVIRONMENT=/path/to/secret_files/environment_2
export GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_2
GO_FIBER_BEHIND_PROXY=/path/to/secret_files/behind_proxy_1 go build
```

The resulting executable will be built with `GO_FIBER_ENVIRONMENT` set to `/path/to/secret_files/environment_2`, `GO_FIBER_BEHIND_PROXY` set to `/path/to/secret_files/behind_proxy_1`, and all other required variables set to their corresponding values in the `.env` file.

### Rationale

The scheme described above is particularly convenient for use with Docker Compose `secrets` configuration:

```yaml
# docker-compose.yml

version: '3.9'

services:
    vaults:
        build:
            context: ./vaults
        command: ./simplepasswords_vaults
        secrets:
            - redis_password
            - vaults_behind_proxy
            - vaults_db_host
            - vaults_db_name
            - vaults_db_password
            - vaults_db_port
            - vaults_db_user
            - vaults_environment
            - vaults_proxy_ip_addresses
            - vaults_secret_key
            - vaults_server_host
            - vaults_server_port
        environment:
            GO_FIBER_ENVIRONMENT: /run/secrets/vaults_environment
            GO_FIBER_BEHIND_PROXY: /run/secrets/vaults_behind_proxy
            GO_FIBER_PROXY_IP_ADDRESSES: /run/secrets/vaults_proxy_ip_addresses
            GO_FIBER_VAULTS_DB_HOST: /run/secrets/vaults_db_host
            GO_FIBER_VAULTS_DB_NAME: /run/secrets/vaults_db_name
            GO_FIBER_VAULTS_DB_PASSWORD: /run/secrets/vaults_db_password
            GO_FIBER_VAULTS_DB_PORT: /run/secrets/vaults_db_port
            GO_FIBER_VAULTS_DB_USER: /run/secrets/vaults_db_user
            GO_FIBER_REDIS_PASSWORD: /run/secrets/redis_password
            GO_FIBER_SECRET_KEY: /run/secrets/vaults_secret_key
            GO_FIBER_SERVER_HOST: /run/secrets/vaults_server_host
            GO_FIBER_SERVER_PORT: /run/secrets/vaults_server_port
        ports:
            - 8080:8080
        depends_on:
            - db_vaults
            - redis
    db_vaults:
        # ...
    redis:
        # ...
secrets:
    redis_password:
        file: ./secret_files/redis_password.txt
    vaults_behind_proxy:
        file: ./secret_files/vaults_behind_proxy.txt
    vaults_db_host:
        file: ./secret_files/vaults_db_host.txt
    vaults_db_name:
        file: ./secret_files/vaults_db_name.txt
    vaults_db_password:
        file: ./secret_files/vaults_db_password.txt
    vaults_db_port:
        file: ./secret_files/vaults_db_port.txt
    vaults_db_user:
        file: ./secret_files/vaults_db_user.txt
    vaults_environment:
        file: ./secret_files/vaults_environment.txt
    vaults_secret_key:
        file: ./secret_files/vaults_secret_key.txt
    vaults_server_host:
        file: ./secret_files/vaults_server_host.txt
    vaults_server_port:
        file: ./secret_files/vaults_server_port.txt
    vaults_proxy_ip_addresses:
        file: ./secret_files/vaults_proxy_ip_addresses.txt
```

## Install Dependencies, Build, & Run

After configuring required environment variables as explained above, run the following commands from the root application folder:

```bash
go mod download
go build
./simplepasswords_vaults
```

The server should now be running at whichever host and port were loaded from `GO_FIBER_SERVER_HOST` and `GO_FIBER_SERVER_PORT` environment variables respectively.

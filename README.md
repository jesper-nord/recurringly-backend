# recurringly-backend

The backend service supporting [Recurringly](https://recurringly.xyz).

## Development

Requires a Postgres instance running locally.

Create an `.env` file in project root with the following contents:

```
PORT=<application port, defaults to 8090>
DB_HOST=<db host, defaults to localhost>
DB_PORT=<db port, defaults to 5432>
DB_NAME=<db name, defaults to postgres>
DB_USERNAME=<username>
DB_PASSWORD=<password>
JWT_SIGNING_SECRET=<any string>
```

## See also
[recurringly-web](https://github.com/jesper-nord/recurringly-web), the web frontend for Recurringly.

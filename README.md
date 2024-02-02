# recurringly-backend

The backend service supporting [Recurringly](https://recurringly.xyz).

## Development

Requires a Postgres instance running locally.

Create an `.env` file in project root:

```
APP_ENV=local
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME=postgres
DB_USERNAME=<username>
DB_PASSWORD=<password>
JWT_SIGNING_SECRET=<any string>
```

## See also
[recurringly-web](https://github.com/jesper-nord/recurringly-web), the web frontend for Recurringly.

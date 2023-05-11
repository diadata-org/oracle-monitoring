Copy the file `env.example` to `.env`, in the root directory, and then fill it.

## Deployment

```shell
DB_CONNECTION_STRING='postgres://user:password@localhost:5432/dbname?sslmode=require'
psql -h <database_host> -U <database_user> -f /path/to/create_database.sql
```

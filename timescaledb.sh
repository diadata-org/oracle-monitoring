docker run --name oracle-monitoring-db -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=root -e POSTGRES_DB=oracle_monitoring_db -p 5432:5432 -d timescale/timescaledb:latest-pg13

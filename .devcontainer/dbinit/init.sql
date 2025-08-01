CREATE DATABASE devdb;

\connect devdb;
\i /docker-entrypoint-initdb.d/schema/init.sql

CREATE DATABASE testdb;

\connect testdb;
\i /docker-entrypoint-initdb.d/schema/init.sql

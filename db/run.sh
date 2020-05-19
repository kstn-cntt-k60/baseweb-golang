PGPASSWORD=1 psql -U postgres -d baseweb < schema.sql
PGPASSWORD=1 psql -U postgres -d baseweb < seed.sql
PGPASSWORD=1 psql -U postgres -d baseweb < demo.sql

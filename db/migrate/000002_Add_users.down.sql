alter table if EXISTS "accounts" drop CONSTRAINT if EXISTS "owner_currency_key";

alter table if EXISTS "accounts" drop CONSTRAINT if EXISTS "accounts_owner_fkey";

DROP TABLE IF EXISTS users;

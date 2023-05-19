CREATE TABLE IF NOT EXISTS "accounts" (
    "id" varchar(28) PRIMARY KEY,
    "first_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "phone" varchar UNIQUE NOT NULL,
    "birth_date" bigint NOT NULL,
    "gender" varchar NOT NULL,
    "show_me" varchar NOT NULL,
    "university" varchar,
    "nsfw" bool NOT NULL,
    "ethnicity" varchar NOT NULL,
    "interests" varchar[],
    "picture" varchar[],
    "verify_yourself" bool default false,
    "about_me" varchar,
    "time_zone" varchar,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
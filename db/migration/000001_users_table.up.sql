CREATE TABLE "users" (
  "id" varchar(28) PRIMARY KEY,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar UNIQUE NOT NULL,
  "age" bigint NOT NULL,
  "gender" varchar NOT NULL,
  "ethnicity" varchar[] NOT NULL,
  "nsfw" bool NOT NULL,
  "metadata" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("age");

CREATE INDEX ON "users" ("gender");

CREATE INDEX ON "users" ("ethnicity");

COMMENT ON COLUMN "users"."metadata" IS 'includes settings data like preference, audio call price, video call price';

COMMENT ON COLUMN "users"."ethnicity" IS 'must be one of American Indian, Black/African Descent, East Asian, Hispanic/Latino, Middle Easter, Pacific Islander, South Asian, Southeast Asian, White/Caucasian, Other';

CREATE INDEX ON "users" ("nsfw");
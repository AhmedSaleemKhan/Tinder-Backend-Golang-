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

CREATE TABLE "accounts" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar(28) NOT NULL,
  "balance" bigint NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "accounts" ("user_id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");


CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "entries" ("account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "entries"."type" IS 'must be one of chat,audio,video,payin,payout';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");


CREATE TABLE "transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" bigint NOT NULL,
  "to_account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

COMMENT ON COLUMN "transfers"."type" IS 'must be one of chat,audio,video';

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "user_id" varchar(28) NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
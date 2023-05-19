CREATE TABLE IF NOT EXISTS "favourites" (
    "fav_id" uuid PRIMARY KEY,
    "user_id" varchar(28) NOT NULL,
    "target_id" varchar(28) NOT NULL,
    "fav_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE ONLY "favourites" 
    ADD CONSTRAINT "user_id_fkey" FOREIGN KEY (user_id) REFERENCES accounts(id) ON DELETE CASCADE;
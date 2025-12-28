CREATE TABLE "user" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "email" varchar NOT NULL
);

ALTER TABLE "notes" ADD FOREIGN KEY ("owner") REFERENCES "user" ("username");
ALTER TABLE "tags" ADD FOREIGN KEY ("owner") REFERENCES "user" ("username");
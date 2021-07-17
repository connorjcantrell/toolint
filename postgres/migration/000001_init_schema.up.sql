CREATE TABLE "users" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL
);

CREATE TABLE "tools" (
  "id" uuid PRIMARY KEY,
  "name" varchar NOT NULL,
  "model" varchar,
  "make" varchar NOT NULL,
  "category" varchar NOT NULL
);

CREATE TABLE "tool_entries" (
  "id" uuid PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "tool_id" uuid NOT NULL,
  "condition" int NOT NULL DEFAULT 0
);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);

ALTER TABLE "tool_entries" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "tool_entries" ADD FOREIGN KEY ("tool_id") REFERENCES "tools" ("id");

CREATE INDEX ON "tools" ("make");

CREATE INDEX ON "tools" ("model");

CREATE INDEX ON "tools" ("category");

CREATE INDEX ON "tool_entries" ("user_id");

CREATE INDEX ON "tool_entries" ("condition");

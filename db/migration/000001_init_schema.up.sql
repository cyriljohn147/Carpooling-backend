CREATE TABLE "users" (
  "id" SERIAL PRIMARY KEY,
  "name" varchar,
  "email" varchar,
  "phone" varchar,
  "balance" decimal(12,2) DEFAULT 0.00,
  "created_at" timestamp DEFAULT NOW()
);

CREATE TABLE "trips" (
  "id" SERIAL PRIMARY KEY,
  "driver_id" integer NOT NULL,
  "date" date,
  "time" time,
  "total_cost" decimal,
  "created_at" timestamp DEFAULT NOW()
);

CREATE TABLE "participants" (
  "id" SERIAL PRIMARY KEY,
  "trip_id" integer NOT NULL,
  "user_id" integer NOT NULL,
  "share_amount" decimal
);

ALTER TABLE "trips"
  ADD CONSTRAINT "trip_driver"
  FOREIGN KEY ("driver_id") REFERENCES "users" ("id");

ALTER TABLE "participants"
  ADD CONSTRAINT "trip_participants"
  FOREIGN KEY ("trip_id") REFERENCES "trips" ("id");

ALTER TABLE "participants"
  ADD CONSTRAINT "participant_users"
  FOREIGN KEY ("user_id") REFERENCES "users" ("id");

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "email" varchar(50),
  "password_hash" varchar(200),
  "role" varchar(10)
);

CREATE TABLE "customer_profiles" (
  "user_id" uuid PRIMARY KEY,
  "first_name" varchar(50),
  "last_name" varchar(50),
  "address" varchar(100),
  "phone_number" varchar(50),
  "created_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "driver_profiles" (
  "user_id" uuid PRIMARY KEY,
  "first_name" varchar(50),
  "last_name" varchar(50),
  "bike" varchar(20),
  "license_plate" varchar(20),
  "phone_number" varchar(50),
  "created_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "order_applicants" (
  "order_id" uuid,
  "driver_id" uuid,
  "created_at" timestamp,
  PRIMARY KEY ("order_id", "driver_id")
);

CREATE TABLE "restaurants" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "name" varchar(50),
  "address" varchar(100),
  "created_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "items" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "restaurant_id" uuid,
  "name" varchar(50),
  "stock" int,
  "price" int,
  "created_at" timestamp DEFAULT (NOW()),
  "updated_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "ratings" (
  "order_id" uuid PRIMARY KEY,
  "rating" int,
  "created_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "orders" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "customer_id" uuid,
  "driver_id" uuid,
  "order_status" varchar(10),
  "delivery_fee" int,
  "total_fee" int,
  "created_at" timestamp DEFAULT (NOW()),
  "updated_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "order_items" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "order_id" uuid,
  "item_id" uuid,
  "quantity" int,
  "created_at" timestamp DEFAULT (NOW())
);

CREATE TABLE "ledgers" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid,
  "amount" int,
  "reason" varchar(50),
  "created_at" timestamp DEFAULT (NOW())
);

CREATE UNIQUE INDEX "idx_users_email_unique" ON "users" ("email");

CREATE UNIQUE INDEX "idx_customer_profiles_user_id_unique" ON "customer_profiles" ("user_id");

CREATE UNIQUE INDEX "idx_driver_profiles_user_id_unique" ON "driver_profiles" ("user_id");

CREATE INDEX "idx_items_restaurant_id" ON "items" ("restaurant_id");

CREATE UNIQUE INDEX "idx_ratings_order_id_unique" ON "ratings" ("order_id");

CREATE INDEX "idx_orders_customer_id" ON "orders" ("customer_id");

CREATE INDEX "idx_orders_driver_id" ON "orders" ("driver_id");

CREATE UNIQUE INDEX "idx_order_items_order_id_item_id_unique" ON "order_items" ("order_id", "item_id");

CREATE INDEX "idx_ledgers_user_id" ON "ledgers" ("user_id");

ALTER TABLE "customer_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "driver_profiles" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_applicants" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_applicants" ADD FOREIGN KEY ("driver_id") REFERENCES "driver_profiles" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "items" ADD FOREIGN KEY ("restaurant_id") REFERENCES "restaurants" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ratings" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD FOREIGN KEY ("customer_id") REFERENCES "customer_profiles" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "orders" ADD FOREIGN KEY ("driver_id") REFERENCES "driver_profiles" ("user_id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_items" ADD FOREIGN KEY ("order_id") REFERENCES "orders" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "order_items" ADD FOREIGN KEY ("item_id") REFERENCES "items" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "ledgers" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

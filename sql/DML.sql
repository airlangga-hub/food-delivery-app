-- 1. Create Users (Password: 'password123')
INSERT INTO "users" ("id", "email", "password_hash", "role") VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'customer@example.com', '$2a$10$20pZXNEqygov7jTWT53Tse0BVmqnbTgr90yoTwRBYcr.HylzjQjEy', 'customer'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'driver@example.com', '$2a$10$20pZXNEqygov7jTWT53Tse0BVmqnbTgr90yoTwRBYcr.HylzjQjEy', 'driver');

-- 2. Create Profiles
INSERT INTO "customer_profiles" ("user_id", "first_name", "last_name", "address", "phone_number") VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'John', 'Doe', 'Jl. Jalan No. 123, Jakarta', '08123456789');

INSERT INTO "driver_profiles" ("user_id", "first_name", "last_name", "bike", "license_plate", "phone_number") VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'Jane', 'Speedy', 'Honda Vario', 'B 1234 ABC', '08987654321');

-- 3. Create Restaurant & Items
INSERT INTO "restaurants" ("id", "name", "address") VALUES
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Gourmet Burger Central', 'Jl. Restoran No. 456, Jakarta');

INSERT INTO "items" ("restaurant_id", "name", "stock", "price") VALUES
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'Classic Cheeseburger', 100, 45000),
('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'French Fries', 200, 15000);

-- 4. Create an Order
INSERT INTO "orders" ("id", "customer_id", "driver_id", "order_status", "delivery_fee", "total_fee") VALUES
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'completed', 10000, 60000);

-- 5. Order Details
INSERT INTO "order_items" ("order_id", "item_id", "quantity") VALUES
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', (SELECT id FROM items WHERE name = 'Classic Cheeseburger'), 1),
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', (SELECT id FROM items WHERE name = 'French Fries'), 1);

-- 6. Ratings & Finance
INSERT INTO "ratings" ("order_id", "rating") VALUES
('d3eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 5);

INSERT INTO "ledgers" ("user_id", "amount", "reason") VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', -60000, 'customer_order'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 60000, 'driver_complete_order');
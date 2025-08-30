# Database Schema (Summary)

This summarizes the Postgres schema implemented in `database/schema/tables.sql`. See the ER diagram for relationships.

## Tables

- Clients (`clients`)
  - id UUID PK (auth user id), name, email UNIQUE, phone, address, role, created_at, updated_at
- Veterinarians (`veterinarians`)
  - id UUID PK (auth user id), name, email UNIQUE, phone, clinic_address, available_hours JSONB, role, timestamps
- Pets (`pets`)
  - id UUID PK, owner_id → clients.id, name, type, breed, date_of_birth, weight, timestamps
- Medical Records (`medical_records`)
  - id UUID PK, pet_id → pets.id, veterinarian_id → veterinarians.id, date_of_visit, reason_for_visit, diagnosis, medication_prescribed TEXT[], notes, timestamps
- QR Codes (`qr_codes`)
  - id UUID PK, pet_id → pets.id, qr_code_data, public_url UNIQUE, encoded_content JSONB, is_active, timestamps
- Appointments (`appointments`)
  - id UUID PK, client_id → clients.id, veterinarian_id → veterinarians.id, pet_id → pets.id, appointment_date, duration_minutes, reason, status, notes, timestamps
- Products (`products`)
  - id UUID PK, veterinarian_id → veterinarians.id, name, description, category, price, stock_quantity, sku UNIQUE, brand, weight, dimensions JSONB, is_prescription_required, is_active, images TEXT[], timestamps
- Orders (`orders`)
  - id UUID PK, client_id → clients.id, veterinarian_id → veterinarians.id, total_amount, status, payment_status, payment_method, shipping_address, delivery_method, notes, timestamps
- Order Items (`order_items`)
  - id UUID PK, order_id → orders.id, product_id → products.id, quantity, unit_price, total_price, created_at

## Indexes (selected)

- pets(owner_id)
- medical_records(pet_id), medical_records(veterinarian_id)
- clients(email), veterinarians(email)
- qr_codes(pet_id), qr_codes(public_url)
- appointments(client_id, veterinarian_id, pet_id, appointment_date)
- products(veterinarian_id, category, sku)
- orders(client_id, veterinarian_id)
- order_items(order_id, product_id)

## Notes

- RLS disabled; auth handled in Go backend via JWT and role checks.
- JSONB fields capture flexible structures (vet available_hours, QR encoded_content, product dimensions).
- Timestamps default to `now()` and most IDs default to `gen_random_uuid()`.

-- Pet Management System Database Schema
-- This file documents the required tables for the Supabase database
-- Clients table (pet owners)
CREATE TABLE IF NOT EXISTS clients (
    id UUID PRIMARY KEY,
    -- This will be the Supabase auth.users ID
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    address TEXT,
    role VARCHAR(50) DEFAULT 'client',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Veterinarians table
CREATE TABLE IF NOT EXISTS veterinarians (
    id UUID PRIMARY KEY,
    -- This will be the Supabase auth.users ID
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(50),
    clinic_address TEXT,
    available_hours JSONB,
    -- Array of working hours objects
    role VARCHAR(50) DEFAULT 'veterinarian',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Pets table
CREATE TABLE IF NOT EXISTS pets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    -- e.g., "Dog", "Cat", "Bird"
    breed VARCHAR(255),
    date_of_birth DATE NOT NULL,
    weight DECIMAL(5, 2),
    -- in kg
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Medical records table
CREATE TABLE IF NOT EXISTS medical_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    veterinarian_id UUID NOT NULL REFERENCES veterinarians(id) ON DELETE CASCADE,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    date_of_visit TIMESTAMP WITH TIME ZONE NOT NULL,
    reason_for_visit TEXT NOT NULL,
    diagnosis TEXT,
    medication_prescribed TEXT [],
    -- Array of medication strings
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- QR Codes table for pet identification system
CREATE TABLE IF NOT EXISTS qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    qr_code_data TEXT NOT NULL,
    -- Base64 encoded QR code image or unique identifier
    public_url TEXT NOT NULL UNIQUE,
    -- Public URL for pet profile access
    encoded_content JSONB NOT NULL,
    -- JSON containing pet info, owner contact details
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Appointments table for scheduling system
CREATE TABLE IF NOT EXISTS appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    veterinarian_id UUID NOT NULL REFERENCES veterinarians(id) ON DELETE CASCADE,
    pet_id UUID NOT NULL REFERENCES pets(id) ON DELETE CASCADE,
    appointment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    duration_minutes INTEGER DEFAULT 30,
    reason TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'scheduled',
    -- scheduled, completed, cancelled, rescheduled
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Products table for e-commerce platform
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    veterinarian_id UUID NOT NULL REFERENCES veterinarians(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    -- food, medicine, accessories, toys, etc.
    price DECIMAL(10, 2) NOT NULL,
    stock_quantity INTEGER DEFAULT 0,
    sku VARCHAR(100) UNIQUE,
    brand VARCHAR(255),
    weight DECIMAL(5, 2),
    -- product weight in kg
    dimensions JSONB,
    -- {length, width, height}
    is_prescription_required BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    images TEXT [],
    -- Array of image URLs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Orders table for purchase management
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    veterinarian_id UUID NOT NULL REFERENCES veterinarians(id) ON DELETE CASCADE,
    -- seller
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending',
    -- pending, confirmed, processing, shipped, delivered, cancelled
    payment_status VARCHAR(50) DEFAULT 'pending',
    -- pending, paid, failed, refunded
    payment_method VARCHAR(50),
    -- card, cash, bank_transfer
    shipping_address TEXT,
    delivery_method VARCHAR(50) DEFAULT 'pickup',
    -- pickup, shipping
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Order items table for individual products in orders
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_pets_owner_id ON pets(owner_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id ON medical_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_veterinarian_id ON medical_records(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_appointment_id ON medical_records(appointment_id);
CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
CREATE INDEX IF NOT EXISTS idx_veterinarians_email ON veterinarians(email);
-- New indexes for PRD features
CREATE INDEX IF NOT EXISTS idx_qr_codes_pet_id ON qr_codes(pet_id);
CREATE INDEX IF NOT EXISTS idx_qr_codes_public_url ON qr_codes(public_url);
CREATE INDEX IF NOT EXISTS idx_appointments_client_id ON appointments(client_id);
CREATE INDEX IF NOT EXISTS idx_appointments_veterinarian_id ON appointments(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_appointments_pet_id ON appointments(pet_id);
CREATE INDEX IF NOT EXISTS idx_appointments_date ON appointments(appointment_date);
CREATE INDEX IF NOT EXISTS idx_products_veterinarian_id ON products(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_products_category ON products(category);
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_orders_client_id ON orders(client_id);
CREATE INDEX IF NOT EXISTS idx_orders_veterinarian_id ON orders(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_product_id ON order_items(product_id);
-- Row Level Security (RLS) is disabled as mentioned in the requirements
-- The Go backend will handle all authorization logic
-- Example available_hours JSON structure for veterinarians:
-- [
--   {
--     "day_of_week": "Monday",
--     "start": "09:00",
--     "end": "17:00"
--   },
--   {
--     "day_of_week": "Tuesday", 
--     "start": "09:00",
--     "end": "17:00"
--   }
-- ]
-- Example qr_codes encoded_content JSON structure:
-- {
--   "pet_name": "Buddy",
--   "pet_type": "Dog",
--   "owner_name": "John Doe", 
--   "owner_phone": "+1234567890",
--   "owner_email": "john@example.com",
--   "owner_address": "123 Main St",
--   "emergency_contact": "+1987654321",
--   "medical_alerts": ["Diabetes", "Allergic to peanuts"],
--   "public_profile_url": "https://petmgt.com/pets/public/abc123"
-- }
-- Example products dimensions JSON structure:
-- {
--   "length": 25.5,
--   "width": 15.0,
--   "height": 10.0,
--   "unit": "cm"
-- }

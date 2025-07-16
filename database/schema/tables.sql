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
    date_of_visit TIMESTAMP WITH TIME ZONE NOT NULL,
    reason_for_visit TEXT NOT NULL,
    diagnosis TEXT,
    medication_prescribed TEXT [],
    -- Array of medication strings
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_pets_owner_id ON pets(owner_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_pet_id ON medical_records(pet_id);
CREATE INDEX IF NOT EXISTS idx_medical_records_veterinarian_id ON medical_records(veterinarian_id);
CREATE INDEX IF NOT EXISTS idx_clients_email ON clients(email);
CREATE INDEX IF NOT EXISTS idx_veterinarians_email ON veterinarians(email);
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

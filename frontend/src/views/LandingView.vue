<template>
  <div class="min-h-screen bg-seasalt">
    <!-- Hero Section -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center">
        <h1 class="text-6xl md:text-8xl font-bold text-rich-black mb-8">
          PetCare
        </h1>
        
        <p class="text-xl md:text-2xl text-gray-600 mb-8 max-w-3xl mx-auto">
          The complete pet management system for pet owners, veterinarians, and pet care professionals. 
          Manage medical records, appointments, and more in one place.
        </p>
        
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <RouterLink to="/signup">
            <Button size="lg" class="w-full sm:w-auto">
              Get Started Free
            </Button>
          </RouterLink>
          <RouterLink to="/login">
            <Button variant="outline" size="lg" class="w-full sm:w-auto">
              Sign In
            </Button>
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Features Section -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center mb-16">
        <h2 class="text-3xl md:text-4xl font-bold text-rich-black mb-4">
          Everything you need for pet care
        </h2>
        <p class="text-lg text-gray-600 max-w-2xl mx-auto">
          From medical records to appointment scheduling, PetCare provides all the tools 
          you need to keep your pets healthy and happy.
        </p>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <!-- Pet Records -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">Pet Records</h3>
          <p class="text-gray-600">
            Keep detailed records of your pets including medical history, vaccinations, 
            and important health information all in one place.
          </p>
        </Card>

        <!-- QR Codes -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">QR Codes</h3>
          <p class="text-gray-600">
            Generate QR codes for your pets that provide instant access to emergency 
            contact information and medical alerts.
          </p>
        </Card>

        <!-- Appointments -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">Appointments</h3>
          <p class="text-gray-600">
            Schedule and manage veterinary appointments with integrated calendar 
            and reminder notifications.
          </p>
        </Card>

        <!-- Medical Records -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">Medical Records</h3>
          <p class="text-gray-600">
            Veterinarians can create and update medical records with detailed 
            visit notes and treatment plans.
          </p>
        </Card>

        <!-- Product Store -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">Pet Products</h3>
          <p class="text-gray-600">
            Browse and purchase pet food, medicine, and accessories directly 
            from veterinarians and pet stores.
          </p>
        </Card>

        <!-- For Professionals -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <div class="w-8 h-8 bg-aquamarine rounded-full"></div>
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">For Professionals</h3>
          <p class="text-gray-600">
            Veterinarians and pet care professionals get powerful tools for 
            managing clients, schedules, and business operations.
          </p>
        </Card>
      </div>
    </div>

    <!-- CTA Section -->
    <div class="bg-white py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        <h2 class="text-3xl md:text-4xl font-bold text-rich-black mb-4">
          Ready to get started?
        </h2>
        <p class="text-lg text-gray-600 mb-8 max-w-2xl mx-auto">
          Join thousands of pet owners and veterinarians who trust PetCare 
          for their pet management needs.
        </p>
        
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <RouterLink to="/signup">
            <Button size="lg" class="w-full sm:w-auto">
              Create Free Account
            </Button>
          </RouterLink>
          <RouterLink to="/login">
            <Button variant="outline" size="lg" class="w-full sm:w-auto">
              Sign In to Your Account
            </Button>
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Footer -->
    <Footer />
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Footer from '@/components/Footer.vue'
</script> 

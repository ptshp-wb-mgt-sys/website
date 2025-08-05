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

      <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
        <!-- Pet Records -->
        <Card class="p-8 text-center">
          <div class="w-16 h-16 bg-aquamarine-100 rounded-full flex items-center justify-center mx-auto mb-6">
            <Heart class="w-8 h-8 text-aquamarine" />
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
            <QrCode class="w-8 h-8 text-aquamarine" />
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
            <Calendar class="w-8 h-8 text-aquamarine" />
          </div>
          <h3 class="text-xl font-semibold text-rich-black mb-4">Appointments</h3>
          <p class="text-gray-600">
            Schedule and manage veterinary appointments with integrated calendar 
            and reminder notifications.
          </p>
        </Card>
      </div>
    </div>

    <!-- Additional Features Section -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-20">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
        <!-- For Pet Owners -->
        <Card class="p-8">
          <div class="flex items-start space-x-4">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center flex-shrink-0">
              <Users class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <h3 class="text-xl font-semibold text-rich-black mb-3">For Pet Owners</h3>
              <p class="text-gray-600 mb-4">
                Access your pet's complete medical history, schedule appointments, 
                and purchase recommended products all in one place.
              </p>
              <ul class="space-y-2 text-sm text-gray-600">
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Medical history tracking
                </li>
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Appointment scheduling
                </li>
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Emergency QR codes
                </li>
              </ul>
            </div>
          </div>
        </Card>

        <!-- For Professionals -->
        <Card class="p-8">
          <div class="flex items-start space-x-4">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center flex-shrink-0">
              <Stethoscope class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <h3 class="text-xl font-semibold text-rich-black mb-3">For Professionals</h3>
              <p class="text-gray-600 mb-4">
                Veterinarians and pet care professionals get powerful tools for 
                managing patient records and growing their practice.
              </p>
              <ul class="space-y-2 text-sm text-gray-600">
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Digital medical records
                </li>
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Product marketplace
                </li>
                <li class="flex items-center">
                  <Check class="w-4 h-4 text-aquamarine mr-2" />
                  Practice management tools
                </li>
              </ul>
            </div>
          </div>
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
import { Heart, QrCode, Calendar, Users, Stethoscope, Check } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Footer from '@/components/Footer.vue'
</script> 

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">Appointments</h1>
      <Button>Book Appointment</Button>
    </div>

    <!-- Appointment Sections -->
    <div class="space-y-8">
      <!-- Upcoming Appointments -->
      <section>
        <h2 class="text-xl font-semibold text-rich-black mb-4">Upcoming</h2>
        <div class="space-y-4">
          <Card class="p-6">
            <div class="flex items-center justify-between">
              <div class="space-y-1">
                <h3 class="font-semibold text-rich-black">Buddy - Annual Checkup</h3>
                <p class="text-sm text-gray-600">Dr. Sarah Johnson</p>
                <p class="text-sm text-gray-600">Today, 2:00 PM - 2:30 PM</p>
              </div>
              <div class="flex space-x-2">
                <Button variant="outline" size="sm">Reschedule</Button>
                <Button variant="ghost" size="sm">Cancel</Button>
              </div>
            </div>
          </Card>

          <Card class="p-6">
            <div class="flex items-center justify-between">
              <div class="space-y-1">
                <h3 class="font-semibold text-rich-black">Max - Vaccination</h3>
                <p class="text-sm text-gray-600">Dr. Mike Wilson</p>
                <p class="text-sm text-gray-600">Tomorrow, 10:00 AM - 10:30 AM</p>
              </div>
              <div class="flex space-x-2">
                <Button variant="outline" size="sm">Reschedule</Button>
                <Button variant="ghost" size="sm">Cancel</Button>
              </div>
            </div>
          </Card>
        </div>
      </section>

      <!-- Past Appointments -->
      <section>
        <h2 class="text-xl font-semibold text-rich-black mb-4">Past Appointments</h2>
        <div class="space-y-4">
          <Card class="p-6">
            <div class="flex items-center justify-between">
              <div class="space-y-1">
                <h3 class="font-semibold text-rich-black">Luna - Dental Cleaning</h3>
                <p class="text-sm text-gray-600">Dr. Sarah Johnson</p>
                <p class="text-sm text-gray-600">Last week, Jan 15, 2024</p>
              </div>
              <div class="flex space-x-2">
                <Button variant="outline" size="sm">View Records</Button>
                <Button variant="ghost" size="sm">Book Again</Button>
              </div>
            </div>
          </Card>

          <Card class="p-6">
            <div class="flex items-center justify-between">
              <div class="space-y-1">
                <h3 class="font-semibold text-rich-black">Buddy - Emergency Visit</h3>
                <p class="text-sm text-gray-600">Dr. Mike Wilson</p>
                <p class="text-sm text-gray-600">Last month, Dec 28, 2023</p>
              </div>
              <div class="flex space-x-2">
                <Button variant="outline" size="sm">View Records</Button>
                <Button variant="ghost" size="sm">Book Again</Button>
              </div>
            </div>
          </Card>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
</script> 

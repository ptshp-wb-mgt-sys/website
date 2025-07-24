<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
</script>

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">Dashboard</h1>
      <Button>Add New Pet</Button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card class="p-6">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Total Pets</p>
            <p class="text-2xl font-bold text-rich-black">12</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Appointments Today</p>
            <p class="text-2xl font-bold text-rich-black">3</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Medical Records</p>
            <p class="text-2xl font-bold text-rich-black">45</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Orders</p>
            <p class="text-2xl font-bold text-rich-black">8</p>
          </div>
        </div>
      </Card>
    </div>

    <!-- Recent Activity -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Recent Pets</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium text-rich-black">Buddy</p>
              <p class="text-sm text-gray-600">Golden Retriever</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium text-rich-black">Luna</p>
              <p class="text-sm text-gray-600">Persian Cat</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Upcoming Appointments</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium text-rich-black">Buddy - Checkup</p>
              <p class="text-sm text-gray-600">Today, 2:00 PM</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
          <div class="flex items-center justify-between">
            <div>
              <p class="font-medium text-rich-black">Max - Vaccination</p>
              <p class="text-sm text-gray-600">Tomorrow, 10:00 AM</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>

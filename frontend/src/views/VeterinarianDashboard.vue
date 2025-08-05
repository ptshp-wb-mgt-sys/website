<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">Veterinarian Dashboard</h1>
        <p class="text-gray-600 mt-1">Manage your practice and patients</p>
      </div>
      <Button>
        <Plus class="w-4 h-4 mr-2" />
        New Medical Record
      </Button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <Calendar class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Today's Appointments</p>
            <p class="text-2xl font-bold text-rich-black">8</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <Users class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Active Patients</p>
            <p class="text-2xl font-bold text-rich-black">127</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <FileText class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Records This Week</p>
            <p class="text-2xl font-bold text-rich-black">23</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <DollarSign class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Revenue This Month</p>
            <p class="text-2xl font-bold text-rich-black">$4,250</p>
          </div>
        </div>
      </Card>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Today's Schedule -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Today's Schedule</h2>
          <Button variant="ghost" size="sm">Manage Availability</Button>
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">10:00 AM - Buddy (Golden Retriever)</p>
              <p class="text-sm text-gray-600">Annual checkup • John Smith</p>
              <p class="text-xs text-green-600">Confirmed</p>
            </div>
            <Button variant="ghost" size="sm">Start Visit</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">11:30 AM - Max (German Shepherd)</p>
              <p class="text-sm text-gray-600">Follow-up • Sarah Johnson</p>
              <p class="text-xs text-blue-600">Waiting</p>
            </div>
            <Button variant="ghost" size="sm">Start Visit</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">2:00 PM - Luna (Persian Cat)</p>
              <p class="text-sm text-gray-600">Vaccination • Mike Davis</p>
            </div>
            <Button variant="ghost" size="sm">View Details</Button>
          </div>
        </div>
      </Card>

      <!-- Recent Patients -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Recent Patients</h2>
          <Button variant="ghost" size="sm">View All Patients</Button>
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-aquamarine-100 rounded-full flex items-center justify-center">
                <Heart class="w-5 h-5 text-aquamarine" />
              </div>
              <div>
                <p class="font-medium text-rich-black">Charlie</p>
                <p class="text-sm text-gray-600">Labrador • Last visit: 2 days ago</p>
              </div>
            </div>
            <Button variant="ghost" size="sm">View Records</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-aquamarine-100 rounded-full flex items-center justify-center">
                <Heart class="w-5 h-5 text-aquamarine" />
              </div>
              <div>
                <p class="font-medium text-rich-black">Bella</p>
                <p class="text-sm text-gray-600">Siamese Cat • Last visit: 1 week ago</p>
              </div>
            </div>
            <Button variant="ghost" size="sm">View Records</Button>
          </div>
        </div>
      </Card>

      <!-- Product Sales -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Product Sales</h2>
          <Button variant="ghost" size="sm">Manage Products</Button>
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Premium Dog Food</p>
              <p class="text-sm text-gray-600">5 units sold this week</p>
              <p class="text-xs text-green-600">$275 revenue</p>
            </div>
            <Button variant="ghost" size="sm">Restock</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Flea Treatment</p>
              <p class="text-sm text-gray-600">3 units sold this week</p>
              <p class="text-xs text-green-600">$180 revenue</p>
            </div>
            <Button variant="ghost" size="sm">View Details</Button>
          </div>
        </div>
      </Card>

      <!-- Quick Actions -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Quick Actions</h2>
        <div class="grid grid-cols-2 gap-3">
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto">
            <FileText class="w-6 h-6 mb-2" />
            <span class="text-sm">New Record</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto">
            <Calendar class="w-6 h-6 mb-2" />
            <span class="text-sm">Set Availability</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto">
            <Package class="w-6 h-6 mb-2" />
            <span class="text-sm">Add Product</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto">
            <QrCode class="w-6 h-6 mb-2" />
            <span class="text-sm">Scan QR Code</span>
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Plus, Heart, Calendar, FileText, Users, DollarSign, Package, QrCode } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
</script>

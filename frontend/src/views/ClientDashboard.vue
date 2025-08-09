<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">My Pet Dashboard</h1>
        <p class="text-gray-600 mt-1">Manage your pets and appointments</p>
      </div>
      <Button @click="openAddPetModal">
        <Plus class="w-4 h-4 mr-2" />
        Add New Pet
      </Button>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <Heart class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">My Pets</p>
            <p class="text-2xl font-bold text-rich-black">{{ petsCount }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <Calendar class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Upcoming Appointments</p>
            <p class="text-2xl font-bold text-rich-black">{{ upcomingCount }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <FileText class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Medical Records</p>
            <p class="text-2xl font-bold text-rich-black">12</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <ShoppingBag class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Recent Orders</p>
            <p class="text-2xl font-bold text-rich-black">1</p>
          </div>
        </div>
      </Card>
    </div>

    <!-- Main Content Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- My Pets -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">My Pets</h2>
          <Button variant="ghost" size="sm" @click="goToMyPets">View All</Button>
        </div>
        <div class="space-y-3">
          <div v-for="pet in myPetsPreview" :key="pet.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-aquamarine-100 rounded-full flex items-center justify-center">
                <Heart class="w-5 h-5 text-aquamarine" />
              </div>
              <div>
                <p class="font-medium text-rich-black">{{ pet.name }}</p>
                <p class="text-sm text-gray-600">{{ pet.breed }} • {{ formatAge(pet.date_of_birth) }}</p>
              </div>
            </div>
            <div class="flex space-x-2">
              <Button variant="ghost" size="sm" @click="goToQr(pet.id)">
                <QrCode class="w-4 h-4" />
              </Button>
              <Button variant="ghost" size="sm" @click="goToMyPets">View</Button>
            </div>
          </div>
          <div v-if="myPetsPreview.length === 0" class="text-sm text-gray-600">No pets yet.</div>
        </div>
      </Card>

      <!-- Upcoming Appointments -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Upcoming Appointments</h2>
          <Button variant="ghost" size="sm">Book New</Button>
        </div>
        <div class="space-y-3">
          <div v-for="appt in upcomingPreview" :key="appt.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">{{ formatApptTitle(appt) }}</p>
              <p class="text-sm text-gray-600">{{ formatApptWhen(appt) }}</p>
              <p v-if="appt.status === 'confirmed'" class="text-xs text-aquamarine">Confirmed</p>
            </div>
            <Button variant="ghost" size="sm" @click="goToAppointments">Details</Button>
          </div>
          <div v-if="upcomingPreview.length === 0" class="text-sm text-gray-600">No upcoming appointments.</div>
        </div>
      </Card>

      <!-- Recent Medical Records -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Recent Medical Records</h2>
          <Button variant="ghost" size="sm">View All</Button>
        </div>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Buddy - Routine Checkup</p>
              <p class="text-sm text-gray-600">Dr. Smith • 3 days ago</p>
              <p class="text-xs text-green-600">All healthy ✓</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Luna - Vaccination</p>
              <p class="text-sm text-gray-600">Dr. Johnson • 1 week ago</p>
              <p class="text-xs text-blue-600">Vaccinated</p>
            </div>
            <Button variant="ghost" size="sm">View</Button>
          </div>
        </div>
      </Card>

      <!-- Quick Actions -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Quick Actions</h2>
        <div class="grid grid-cols-2 gap-3">
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToAppointments">
            <Calendar class="w-6 h-6 mb-2" />
            <span class="text-sm">Book Appointment</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="openAddPetModal">
            <Plus class="w-6 h-6 mb-2" />
            <span class="text-sm">Add Pet</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToProducts">
            <ShoppingBag class="w-6 h-6 mb-2" />
            <span class="text-sm">Shop Products</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToMyPets">
            <QrCode class="w-6 h-6 mb-2" />
            <span class="text-sm">QR Codes</span>
          </Button>
        </div>
      </Card>
      <!-- Add Pet Modal -->
      <AddPetModal
        :is-open="showAddPetModal"
        @close="closeAddPetModal"
        @pet-added="handlePetAdded"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus, Heart, Calendar, FileText, ShoppingBag, QrCode } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import AddPetModal from '@/components/AddPetModal.vue'
import { useRouter } from 'vue-router'
import { usePetsStore } from '@/stores/pets'
import { useAppointmentsStore, type Appointment } from '@/stores/appointments'

// Modal state
const showAddPetModal = ref(false)

/**
 * Open the add-pet modal; simple toggle.
 */
const openAddPetModal = () => {
  showAddPetModal.value = true
}

/**
 * Close the add-pet modal.
 */
const closeAddPetModal = () => {
  showAddPetModal.value = false
}

/**
 * When a pet is added successfully, close the modal.
 */
const handlePetAdded = () => {
  closeAddPetModal()
}

// Stores
const petsStore = usePetsStore()
const appointmentsStore = useAppointmentsStore()
const router = useRouter()

/**
 * Number of pets for stats card.
 */
const petsCount = computed(() => petsStore.petsCount)

/**
 * Preview of pets (first 2) for dashboard card.
 */
const myPetsPreview = computed(() => petsStore.pets.slice(0, 2))

/**
 * Count of upcoming appointments for stats card.
 */
const upcomingCount = computed(() => appointmentsStore.upcomingAppointments.length)

/**
 * Preview of upcoming appointments (first 2).
 */
const upcomingPreview = computed(() => appointmentsStore.upcomingAppointments.slice(0, 2))

/**
 * On mount, ensure data is loaded.
 */
onMounted(async () => {
  if (petsStore.pets.length === 0) {
    await petsStore.fetchPets()
  }
  if (appointmentsStore.appointments.length === 0) {
    await appointmentsStore.fetchAppointments()
  }
})

/**
 * Format age from date string.
 */
const formatAge = (dateOfBirth: string) => {
  const birthDate = new Date(dateOfBirth + 'T00:00:00')
  const today = new Date()
  const years = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    return `${years - 1} years`
  }
  if (years === 0) {
    const months = monthDiff + (today.getDate() >= birthDate.getDate() ? 0 : -1)
    return months <= 0 ? 'Less than 1 month' : `${months} months`
  }
  return `${years} years`
}

/**
 * Helpers to format appointments.
 */
const formatApptTitle = (appt: Appointment) => `${appt.reason}`
const formatApptWhen = (appt: Appointment) => new Date(appt.appointment_date).toLocaleString()

/**
 * Navigation helpers.
 */
const goToMyPets = () => router.push({ name: 'my-pets' })
const goToAppointments = () => router.push({ name: 'book-appointment' })
const goToQr = (petId: string) => router.push({ name: 'my-pets' })
const goToProducts = () => router.push({ name: 'browse-products' })
</script>


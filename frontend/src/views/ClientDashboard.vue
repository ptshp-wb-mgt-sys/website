<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">My Pet Dashboard</h1>
        <p class="text-gray-600 mt-1">Manage your pets and appointments</p>
      </div>
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
            <p class="text-2xl font-bold text-rich-black">{{ recordsCount }}</p>
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
              <Button variant="ghost" size="sm" @click="goToPetProfile(pet.id)">View</Button>
            </div>
          </div>
          <div v-if="myPetsPreview.length === 0" class="text-sm text-gray-600">No pets yet.</div>
        </div>
      </Card>

      <!-- Upcoming Appointments -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Upcoming Appointments</h2>
          <Button variant="ghost" size="sm" @click="goToAppointments">Book New</Button>
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
          <Button variant="ghost" size="sm" @click="goToMyPets">View All</Button>
        </div>
        <div class="space-y-3">
          <div
            v-for="rec in recentRecords.slice(0, 2)"
            :key="rec.id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
          >
            <div>
              <p class="font-medium text-rich-black">{{ petNameById[rec.pet_id] || 'Pet' }} — {{ rec.reason_for_visit }}</p>
              <p class="text-sm text-gray-600">{{ formatDate(rec.date_of_visit) }}</p>
              <p v-if="rec.diagnosis" class="text-xs text-gray-600">{{ rec.diagnosis }}</p>
            </div>
            <Button variant="ghost" size="sm" @click="goToPetProfile(rec.pet_id)">View</Button>
          </div>
          <div v-if="recentRecords.length === 0" class="text-sm text-gray-600">No records yet.</div>
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
            <Heart class="w-6 h-6 mb-2" />
            <span class="text-sm">View Pets</span>
          </Button>
        </div>
      </Card>
      <!-- Add Pet Modal -->
      <AddPetModal
        :is-open="showAddPetModal"
        @close="closeAddPetModal"
        @pet-added="handlePetAdded"
      />

      <!-- QR Preview Modal -->
      <QRPreviewModal :is-open="showQRModal" :pet-id="qrPetId" @close="() => { showQRModal = false; qrPetId = '' }" />
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
import { useMedicalRecordsStore, type MedicalRecord } from '@/stores/medicalRecords'
import QRPreviewModal from '@/components/QRPreviewModal.vue'
import { useQRCodesStore } from '@/stores/qrcodes'

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
const recordsStore = useMedicalRecordsStore()
const qrStore = useQRCodesStore()

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

// Recent medical records across user's pets
const recentRecords = computed<MedicalRecord[]>(() => {
  const loaded: MedicalRecord[] = []
  for (const pet of petsStore.pets) {
    const bucket = recordsStore.getCachedForPet(pet.id)
    if (bucket && bucket.length) loaded.push(...bucket)
  }
  return loaded
    .sort((a, b) => new Date(b.date_of_visit).getTime() - new Date(a.date_of_visit).getTime())
})

// Count for stats card
const recordsCount = computed(() => recentRecords.value.length)

// Quick lookup of pet names by id
const petNameById = computed<Record<string, string>>(() => {
  const map: Record<string, string> = {}
  for (const p of petsStore.pets) map[p.id] = p.name
  return map
})

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
  // Preload medical records for each pet (best-effort)
  await Promise.all(petsStore.pets.map(p => recordsStore.fetchByPetId(p.id)))
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
const goToProducts = () => router.push({ name: 'browse-products' })
const goToPetProfile = (petId: string) => router.push({ name: 'pet-profile', params: { id: petId } })

// Small utils
const formatDate = (iso: string) => new Date(iso).toLocaleDateString()

// QR Preview modal state
const showQRModal = ref(false)
const qrPetId = ref<string>('')

/**
 * Open QR preview modal for a pet; generates QR if missing.
 */
const openQRModal = async (petId: string) => {
  try {
    await qrStore.getOrCreateForPet(petId)
    qrPetId.value = petId
    showQRModal.value = true
  } catch (e) {
    console.error('Failed to open QR preview', e)
  }
}

/**
 * Dashboard QR button handler; just delegates to modal logic.
 */
const goToQr = async (petId: string) => {
  await openQRModal(petId)
}
</script>

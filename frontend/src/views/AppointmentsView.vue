<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">{{ pageTitle }}</h1>
        <p class="text-gray-600 mt-1">{{ pageSubtitle }}</p>
      </div>
      <Button v-if="userStore.isVeterinarian" @click="showAvailability = true">
        <Settings class="w-4 h-4 mr-2" />
        Manage Availability
      </Button>
    </div>

    <!-- CLIENT VIEW: Appointment Booking -->
    <template v-if="userStore.isClient">
      <!-- Quick Book Section -->
      <Card v-if="!showReschedule" class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Book New Appointment</h2>
          <Button :disabled="!canFindSlots" @click="loadSlots">Find Available Times</Button>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4 items-start">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Select Pet</label>
            <select v-model="form.pet_id" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
              <option value="" disabled>Select pet</option>
              <option v-for="pet in pets" :key="pet.id" :value="pet.id">{{ pet.name }} ({{ pet.breed }})</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Select Veterinarian</label>
            <div ref="vetDropdownRef" class="relative">
              <button type="button" @click="toggleVetOpen" class="w-full flex items-center justify-between px-3 py-2 border border-gray-300 rounded-lg bg-white">
                <span class="truncate text-left">
                  {{ selectedVetLabel || 'Select veterinarian' }}
                </span>
                <svg class="w-4 h-4 text-gray-500 ml-2" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                  <path fill-rule="evenodd" d="M5.23 7.21a.75.75 0 011.06.02L10 10.94l3.71-3.71a.75.75 0 111.06 1.06l-4.24 4.24a.75.75 0 01-1.06 0L5.21 8.29a.75.75 0 01.02-1.08z" clip-rule="evenodd" />
                </svg>
              </button>

              <div v-if="vetOpen" class="absolute z-10 mt-1 w-full bg-white border border-gray-200 rounded-lg shadow-lg">
                <div class="p-2 border-b border-gray-200 space-y-2">
                  <input
                    v-model="vetSearch"
                    type="text"
                    placeholder="Search by name or clinic"
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg"
                  />
                  <select v-model="vetLocation" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
                    <option value="">All cities</option>
                    <option v-for="loc in vetLocations" :key="loc" :value="loc">{{ loc }}</option>
                  </select>
                </div>
                <ul class="max-h-64 overflow-auto py-1">
                  <li
                    v-for="vet in filteredVeterinarians"
                    :key="vet.id"
                    @click="selectVet(vet.id)"
                    class="px-3 py-2 hover:bg-gray-50 cursor-pointer flex justify-between items-center"
                  >
                    <span class="truncate">{{ vet.name }}</span>
                    <span v-if="vet.clinic_address" class="ml-2 text-xs text-gray-500 truncate">{{ vet.clinic_address }}</span>
                  </li>
                  <li v-if="filteredVeterinarians.length === 0" class="px-3 py-2 text-sm text-gray-600">No veterinarians found</li>
                </ul>
              </div>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Reason for Visit</label>
            <div class="space-y-2">
              <select v-model="form.reason" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
                <option value="" disabled>Select reason</option>
                <option>Routine Checkup</option>
                <option>Vaccination</option>
                <option>Emergency</option>
                <option>Follow-up</option>
                <option>Illness</option>
                <option>Injury</option>
                <option>Dental Cleaning</option>
                <option>Spay/Neuter</option>
                <option>Lab Work</option>
                <option>Behavior Consultation</option>
                <option>Grooming</option>
                <option>Medication Refill</option>
                <option>Surgery Consultation</option>
                <option value="other">Other</option>
              </select>
              <input
                v-if="form.reason === 'other'"
                v-model="customReason"
                type="text"
                placeholder="Type your reason"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg"
              />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date</label>
            <input v-model="date" type="date" :min="todayISODate" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
          </div>
        </div>
        <p v-if="uiError" class="text-sm text-red-600 mt-3">{{ uiError }}</p>
        <div v-if="slots.length" class="grid grid-cols-2 md:grid-cols-4 gap-2 mt-4">
          <Button v-for="s in slots" :key="s.start_time" variant="outline" :disabled="isSlotInPast(s)" @click="selectSlot(s)">
            {{ formatTimeHM(s.start_time, false) }}
          </Button>
        </div>
        <p v-else-if="triedLoadSlots" class="text-sm text-gray-600 mt-2">No available time slots for the selected day.</p>
        <div class="mt-4" v-if="selectedSlot">
          <Button @click="book">Confirm Booking</Button>
        </div>
      </Card>

      <!-- Reschedule Appointment -->
      <Card v-if="showReschedule" class="p-6 mt-4" data-reschedule>
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Reschedule Appointment</h2>
          <div class="space-x-2">
            <Button variant="outline" @click="showReschedule = false">Back to Booking</Button>
            <Button variant="ghost" @click="showReschedule = false">Close</Button>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <input v-model="date" type="date" :min="todayISODate" class="px-3 py-2 border border-gray-300 rounded-lg" />
          <Button @click="loadSlots">Find Available Times</Button>
        </div>
        <div v-if="slots.length" class="grid grid-cols-2 md:grid-cols-4 gap-2 mt-4">
          <Button v-for="s in slots" :key="s.start_time" variant="outline" @click="selectSlot(s)">
            {{ formatTimeHM(s.start_time, false) }}
          </Button>
        </div>
        <div class="mt-4 flex space-x-2" v-if="selectedSlot">
          <Button @click="confirmReschedule">Confirm</Button>
          <Button variant="ghost" @click="showReschedule = false">Cancel</Button>
        </div>
      </Card>

      <!-- My Appointments -->
      <div class="space-y-8">
        <section>
          <h2 class="text-xl font-semibold text-rich-black mb-4">My Upcoming Appointments</h2>
           <div class="space-y-4">
            <Card v-for="appt in upcomingAppointments" :key="appt.id" class="p-6">
              <div class="flex items-center justify-between">
                <div class="space-y-1">
                  <h3 class="font-semibold text-rich-black">{{ appt.reason }} • <span class="text-gray-600">{{ petLabel(appt.pet_id) }}</span></h3>
                  <p class="text-sm text-gray-600">{{ formatDateTimeMDYHM(appt.appointment_date) }}</p>
                  <p class="text-xs" :class="appt.status === 'confirmed' ? 'text-green-600' : 'text-gray-600'">{{ appt.status }}</p>
                </div>
                <div class="flex space-x-2">
                  <Button variant="outline" size="sm" @click="openReschedule(appt)">Reschedule</Button>
                  <Button variant="outline" size="sm" @click="$router.push({ name: 'pet-profile', params: { id: appt.pet_id } })">View Records</Button>
                  <Button variant="ghost" size="sm" :disabled="appt.status === 'cancelled'" class="text-red-600 disabled:opacity-50" @click="cancel(appt.id)">Cancel</Button>
                </div>
              </div>
            </Card>
          </div>
        </section>

        <section>
          <h2 class="text-xl font-semibold text-rich-black mb-4">Past Appointments</h2>
          <div class="space-y-4">
            <Card v-for="appt in pastAppointments" :key="appt.id" class="p-6">
              <div class="flex items-center justify-between">
                <div class="space-y-1">
                  <h3 class="font-semibold text-rich-black">{{ appt.reason }}</h3>
                  <p class="text-sm text-gray-600">{{ formatDateTimeMDYHM(appt.appointment_date) }}</p>
                  <p class="text-xs" :class="appt.status === 'completed' ? 'text-green-600' : 'text-gray-600'">{{ appt.status }}</p>
                </div>
                <div class="flex space-x-2">
                  <Button variant="outline" size="sm" @click="$router.push({ name: 'pet-profile', params: { id: appt.pet_id } })">View Records</Button>
                  <Button variant="ghost" size="sm" @click="openReschedule(appt)">Book Again</Button>
                </div>
              </div>
            </Card>
          </div>
        </section>
      </div>
    </template>

    <!-- VETERINARIAN VIEW: Schedule Management -->
    <template v-if="userStore.isVeterinarian">
      <!-- Availability Editor -->
      <Card v-if="showAvailability" class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Edit Availability</h2>
          <div class="space-x-2">
            <Button variant="outline" @click="addAvailabilityRow">Add Schedule</Button>
            <Button @click="saveAvailability">Save</Button>
            <Button variant="ghost" @click="showAvailability = false">Close</Button>
          </div>
        </div>
        <div class="space-y-2">
          <div class="hidden md:grid grid-cols-12 gap-2 text-xs text-gray-500 px-1">
            <div class="col-span-3">Day</div>
            <div class="col-span-4">Start</div>
            <div class="col-span-4">End</div>
            <div class="col-span-1 text-right">Actions</div>
          </div>

          <div v-for="(row, idx) in availability" :key="idx" :class="['grid grid-cols-1 md:grid-cols-12 gap-2 items-center p-2 rounded-lg border', isInvalid(row) ? 'border-red-300 bg-red-50' : 'border-gray-200 bg-white']">
            <select v-model="row.day_of_week" class="col-span-12 md:col-span-3 px-3 py-2 border border-gray-300 rounded-lg">
              <option>Sun</option>
              <option>Mon</option>
              <option>Tue</option>
              <option>Wed</option>
              <option>Thu</option>
              <option>Fri</option>
              <option>Sat</option>
            </select>
            <input v-model="row.start" type="time" step="900" class="col-span-12 md:col-span-4 px-3 py-2 border border-gray-300 rounded-lg" />
            <input v-model="row.end" type="time" step="900" class="col-span-12 md:col-span-4 px-3 py-2 border border-gray-300 rounded-lg" />
            <div class="col-span-12 md:col-span-1 flex md:justify-end">
              <Button variant="ghost" size="sm" class="text-red-600" @click="removeAvailabilityRow(idx)">Remove</Button>
            </div>
            <div v-if="isInvalid(row)" class="md:col-span-12 col-span-12 text-xs text-red-600 px-1">End time must be after start time.</div>
          </div>
        </div>
      </Card>

      <!-- Today's Schedule (only show if there are appointments today) -->
      <Card v-if="todaysAppointments.length > 0" class="p-6">
        <div class="mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Today's Schedule - {{ currentDate }}</h2>
        </div>
        
        <div class="space-y-3">
          <div v-for="appt in todaysAppointments" :key="appt.id" class="flex items-center justify-between p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <div class="space-y-1">
              <p class="font-medium text-rich-black">{{ formatTimeHM(appt.appointment_date, false) }} - {{ appt.reason }} • <span class="text-gray-600 text-sm">{{ petLabel(appt.pet_id) }}</span></p>
              <span class="inline-flex items-center text-xs px-2 py-0.5 rounded-full"
                :class="{
                  'bg-green-100 text-green-700': appt.status === 'completed',
                  'bg-yellow-100 text-yellow-700': appt.status === 'scheduled' || appt.status === 'rescheduled',
                  'bg-gray-100 text-gray-700': appt.status !== 'completed' && appt.status !== 'scheduled' && appt.status !== 'rescheduled',
                }">
                {{ appt.status }}
              </span>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" @click="$router.push({ name: 'pet-profile', params: { id: appt.pet_id } })">View</Button>
              <Button variant="outline" size="sm" :disabled="appt.status === 'completed'" @click="markCompleted(appt.id)">Mark Completed</Button>
              <Button size="sm" @click="openQR(appt.pet_id)" title="QR Code"><QrCode class="w-4 h-4" /></Button>
            </div>
          </div>
        </div>
      </Card>

      <!-- Upcoming Appointments (small) -->
      <Card v-if="vetUpcoming.length > 0" class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Upcoming Appointments</h2>
        <div class="space-y-3">
          <div v-for="appt in vetUpcoming" :key="appt.id" class="flex items-center justify-between p-3 bg-gray-50 border border-gray-200 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">{{ formatDateTimeMDYHM(appt.appointment_date) }} - {{ appt.reason }} • <span class="text-gray-600">{{ petLabel(appt.pet_id) }}</span></p>
              <p class="text-xs text-gray-600 capitalize">Status: {{ appt.status }}</p>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" @click="$router.push({ name: 'pet-profile', params: { id: appt.pet_id } })">View</Button>
              <Button variant="ghost" size="sm" class="text-red-600" :disabled="appt.status === 'cancelled'" @click="cancel(appt.id)">Cancel</Button>
            </div>
          </div>
        </div>
      </Card>

      <!-- Availability Management -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">This Week's Availability</h2>
        <div class="grid grid-cols-7 gap-2">
          <div v-for="day in availabilityDays" :key="day.name" class="text-center">
            <p class="text-sm font-medium text-gray-700 mb-2">{{ day.name }}</p>
            <div class="space-y-1">
              <div v-if="day.slots.length === 0" class="text-xs p-1 rounded bg-gray-100 text-gray-600">Closed</div>
              <div v-else v-for="slot in day.slots" :key="slot" class="text-xs p-1 rounded bg-green-100 text-green-700">
                {{ slot }}
              </div>
            </div>
          </div>
        </div>
      </Card>

      <!-- Recent Patients -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Recent Patients</h2>
        <div class="space-y-3">
          <div v-if="recentCompleted.length === 0" class="text-sm text-gray-600">No recent completed appointments.</div>
          <div v-else v-for="appt in recentCompleted" :key="appt.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">{{ appt.reason }}</p>
              <p class="text-sm text-gray-600">{{ formatDateTimeMDYHM(appt.appointment_date) }}</p>
            </div>
            <Button variant="ghost" size="sm" @click="$router.push({ name: 'pet-profile', params: { id: appt.pet_id } })">View Records</Button>
          </div>
        </div>
      </Card>
    </template>

    <!-- ADMIN VIEW: System Overview -->
    <template v-if="userStore.isAdmin">
      <!-- System Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
        <Card class="p-6">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
              <Calendar class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-600">Today's Appointments</p>
              <p class="text-2xl font-bold text-rich-black">23</p>
            </div>
          </div>
        </Card>
        <Card class="p-6">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
              <Users class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-600">Active Veterinarians</p>
              <p class="text-2xl font-bold text-rich-black">8</p>
            </div>
          </div>
        </Card>
        <Card class="p-6">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
              <TrendingUp class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-600">This Week</p>
              <p class="text-2xl font-bold text-rich-black">156</p>
            </div>
          </div>
        </Card>
        <Card class="p-6">
          <div class="flex items-center">
            <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
              <AlertCircle class="w-6 h-6 text-aquamarine" />
            </div>
            <div>
              <p class="text-sm font-medium text-gray-600">Cancellations</p>
              <p class="text-2xl font-bold text-rich-black">3</p>
            </div>
          </div>
        </Card>
      </div>

      <!-- All Appointments -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">All Appointments</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Buddy - Dr. Johnson - 10:00 AM</p>
              <p class="text-sm text-gray-600">John Smith • Checkup • Confirmed</p>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm">View</Button>
              <Button variant="ghost" size="sm">Manage</Button>
            </div>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Luna - Dr. Wilson - 2:00 PM</p>
              <p class="text-sm text-gray-600">Mike Davis • Vaccination • Scheduled</p>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm">View</Button>
              <Button variant="ghost" size="sm">Manage</Button>
            </div>
          </div>
        </div>
      </Card>
    </template>
    <QRPreviewModal :isOpen="qrOpen" :petId="qrPetId" @close="qrOpen = false" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import { useAppointmentsStore, type CreateAppointmentRequest, type TimeSlot } from '@/stores/appointments'
import { 
  Plus, Settings, Calendar, Clock, Users, TrendingUp, AlertCircle, QrCode 
} from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import QRPreviewModal from '@/components/QRPreviewModal.vue'
import { useQRCodesStore } from '@/stores/qrcodes'
import { usePetsStore } from '@/stores/pets'
import { formatTimeHM, formatDateTimeMDYHM } from '@/lib/utils'

const userStore = useUserStore()
const petsStore = usePetsStore()
const apptStore = useAppointmentsStore()
const qrStore = useQRCodesStore()

// Booking UI state
const showAvailability = ref(false)
const showReschedule = ref(false)
const rescheduleTarget = ref<string | null>(null)
const veterinarians = ref<{ id: string; name: string; clinic_address?: string }[]>([])
const vetSearch = ref('')
const vetLocation = ref('')
const vetOpen = ref(false)
const vetDropdownRef = ref<HTMLElement | null>(null)
const slots = ref<TimeSlot[]>([])
const selectedSlot = ref<TimeSlot | null>(null)
const date = ref<string>('')
const form = ref<CreateAppointmentRequest>({ veterinarian_id: '', pet_id: '', appointment_date: '', reason: '' })
const customReason = ref('')
const uiError = ref<string | null>(null)
const triedLoadSlots = ref(false)
const canFindSlots = computed(() => !!form.value.veterinarian_id && !!form.value.pet_id && !!form.value.reason && !!date.value)

// Availability editor model seeded from profile to avoid initial flicker
const initialHours = (userStore.isVeterinarian && (userStore.profile as any)?.available_hours)
  ? ((userStore.profile as any).available_hours as Array<{ day_of_week: string; start: string; end: string }>)
  : []
const availability = ref<{ day_of_week: string; start: string; end: string }[]>(
  Array.isArray(initialHours)
    ? initialHours.map(h => ({ day_of_week: h.day_of_week, start: h.start, end: h.end }))
    : [],
)

onMounted(async () => {
  if (userStore.isClient) {
    if (petsStore.pets.length === 0) await petsStore.fetchPets()
    await apptStore.fetchAppointments()
    const vets = await apptStore.listVeterinarians()
    veterinarians.value = vets.map(v => ({ id: v.id, name: v.name, clinic_address: v.clinic_address }))
  }
  if (userStore.isVeterinarian || userStore.isAdmin) {
    // Load once per session unless nothing is cached yet
    if (apptStore.appointments.length === 0) {
      await apptStore.fetchAppointments()
    }
    // Load pets on demand for vet views when showing labels (best effort)
    try { await petsStore.fetchPets() } catch (_) {}
    // Also fetch labels for pets not owned by the vet
    try {
      const ids = apptStore.appointments.map(a => a.pet_id)
      await ensurePetLabelsFor(ids)
    } catch (_) {}
  }
  // Sync with profile after mount as well
  if (userStore.isVeterinarian) {
    const hours = (userStore.profile as any)?.available_hours as Array<{ day_of_week: string; start: string; end: string }> | undefined
    if (Array.isArray(hours)) {
      availability.value = hours.map(h => ({ day_of_week: h.day_of_week, start: h.start, end: h.end }))
    }
  }
})

// Close custom dropdown on outside click
function onClickOutside(event: MouseEvent) {
  const target = event.target as Node
  if (!vetDropdownRef.value) return
  if (!vetDropdownRef.value.contains(target)) {
    vetOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

// Cleanup listener
watch(vetOpen, (open) => {
  if (!open) return
})

onMounted(() => {})


// Reflect profile availability updates to the local availability model
watch(
  () => (userStore.profile as any)?.available_hours,
  (hours) => {
    if (!userStore.isVeterinarian) return
    if (Array.isArray(hours)) {
      availability.value = hours.map(h => ({ day_of_week: h.day_of_week, start: h.start, end: h.end }))
    }
  },
)

/**
 * Computed properties for dynamic content
 */
const pageTitle = computed(() => {
  if (userStore.isClient) return 'Appointments'
  if (userStore.isVeterinarian) return 'My Schedule'
  if (userStore.isAdmin) return 'All Appointments'
  return 'Appointments'
})

const pageSubtitle = computed(() => {
  if (userStore.isClient) return 'Book and manage your pet appointments'
  if (userStore.isVeterinarian) return 'Manage your schedule and patient appointments'
  if (userStore.isAdmin) return 'System-wide appointment management and analytics'
  return ''
})

const currentDate = computed(() => {
  return new Date().toLocaleDateString('en-US', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
})

// Compute availability display from saved availability rows
const availabilityDays = computed(() => {
  // Start week on Sunday
  const days = ['Sun','Mon','Tue','Wed','Thu','Fri','Sat']
  const map = days.map(name => ({ name, slots: [] as string[] }))
  for (const row of availability.value) {
    const idx = days.indexOf(row.day_of_week)
    if (idx >= 0) {
      map[idx].slots.push(`${row.start} - ${row.end}`)
    }
  }
  return map
})

// Derived lists
const upcomingAppointments = computed(() => apptStore.upcomingAppointments)
const pastAppointments = computed(() => apptStore.pastAppointments)
const todaysAppointments = computed(() => apptStore.todaysAppointments)
// Vet filters
const filteredVeterinarians = computed(() => {
  const q = vetSearch.value.trim().toLowerCase()
  const loc = vetLocation.value
  return veterinarians.value.filter(v => {
    const matchesLoc = !loc || extractCity(v.clinic_address).toLowerCase() === loc.toLowerCase()
    if (!q) return matchesLoc
    const hay = `${v.name} ${(v.clinic_address || '')}`.toLowerCase()
    return matchesLoc && hay.includes(q)
  })
})

const vetLocations = computed(() => {
  const set = new Set<string>()
  for (const v of veterinarians.value) {
    const city = extractCity(v.clinic_address)
    if (city) set.add(city)
  }
  return Array.from(set).sort((a, b) => a.localeCompare(b))
})

/**
 * Pull out the city from a clinic address.
 * We expect the city to be the SECOND comma-separated segment.
 * Falls back to the first segment if there's only one.
 */
const extractCity = (address?: string): string => {
  if (!address) return ''
  const parts = address.split(',').map(p => p.trim()).filter(Boolean)
  if (parts.length >= 2) return parts[1]
  if (parts.length === 1) return parts[0]
  return ''
}

const selectedVetLabel = computed(() => {
  const id = form.value.veterinarian_id
  if (!id) return ''
  const vet = veterinarians.value.find(v => v.id === id)
  if (!vet) return ''
  return vet.clinic_address ? `${vet.name} — ${vet.clinic_address}` : vet.name
})

const toggleVetOpen = () => {
  vetOpen.value = !vetOpen.value
}

const selectVet = (id: string) => {
  form.value.veterinarian_id = id
  vetOpen.value = false
}
// For veterinarians: small upcoming widget (next 5), exclude today's appointments
const vetUpcoming = computed(() => {
  const t = new Date()
  t.setHours(0, 0, 0, 0)
  t.setDate(t.getDate() + 1) // start of tomorrow
  const startOfTomorrow = t.getTime()
  return apptStore.upcomingAppointments
    .filter(a => new Date(a.appointment_date).getTime() >= startOfTomorrow)
    .slice(0, 5)
})

// Recently completed appointments for vet
const recentCompleted = computed(() => {
  return apptStore.pastAppointments.filter(a => a.status === 'completed').slice(0, 5)
})

// Pets for select
const pets = computed(() => petsStore.pets)

// Cache of pet labels for vet views (since vets don't own the pets list)
const petLabels = ref<Record<string, string>>({})

async function ensurePetLabelsFor(ids: string[]): Promise<void> {
  const unique = Array.from(new Set(ids)).filter(id => !petLabels.value[id])
  await petsStore.warmPetLabels(unique)
  unique.forEach(id => {
    const label = petsStore.getPetLabelSync(id)
    if (label) petLabels.value[id] = label
  })
}

function petLabel(petId: string): string {
  const sync = petsStore.getPetLabelSync(petId)
  return sync || petLabels.value[petId] || ''
}

/**
 * Load available time slots for the selected vet and date.
 */
const loadSlots = async () => {
  uiError.value = null
  triedLoadSlots.value = false
  if (!form.value.veterinarian_id || !date.value) {
    uiError.value = 'Please select a veterinarian and a date.'
    return
  }
  const iso = new Date(date.value + 'T00:00:00').toISOString()
  const fetched = await apptStore.getAvailableSlots(form.value.veterinarian_id, iso)
  const now = Date.now()
  // Filter out past times (including earlier today)
  slots.value = fetched.filter(s => new Date(s.start_time).getTime() >= now)
  triedLoadSlots.value = true
}

/**
 * Choose a specific slot; sets the appointment start time.
 */
const selectSlot = (slot: TimeSlot) => {
  selectedSlot.value = slot
  form.value.appointment_date = slot.start_time
}

/**
 * Book the appointment with the selected options.
 */
const book = async () => {
  await apptStore.createAppointment({
    veterinarian_id: form.value.veterinarian_id,
    pet_id: form.value.pet_id,
    appointment_date: form.value.appointment_date,
    reason: form.value.reason === 'other' ? customReason.value : form.value.reason,
  })
  // Reset UI
  form.value = { veterinarian_id: '', pet_id: '', appointment_date: '', reason: '' }
  date.value = ''
  slots.value = []
  selectedSlot.value = null
  triedLoadSlots.value = false
  customReason.value = ''
}

// Cancellation
const cancel = async (id: string) => {
  await apptStore.cancelAppointment(id)
}

// Mark appointment completed
const markCompleted = async (id: string) => {
  await apptStore.updateAppointment(id, { status: 'completed' })
}

// Reschedule flow
const openReschedule = async (appt: { id: string; veterinarian_id: string; appointment_date?: string }) => {
  rescheduleTarget.value = appt.id
  form.value.veterinarian_id = appt.veterinarian_id
  if (appt.appointment_date) {
    const d = new Date(appt.appointment_date)
    const yyyy = d.getFullYear()
    const mm = String(d.getMonth() + 1).padStart(2, '0')
    const dd = String(d.getDate()).padStart(2, '0')
    date.value = `${yyyy}-${mm}-${dd}`
  }
  showReschedule.value = true
  await nextTick()
  // Auto-load slots for the selected day and vet
  try {
    await loadSlots()
  } catch (_) {}
  nextTick(() => {
    // Optionally scroll to reschedule panel
    const el = document.querySelector('[data-reschedule]') as HTMLElement | null
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  })
}

const confirmReschedule = async () => {
  if (!rescheduleTarget.value || !form.value.appointment_date) return
  await apptStore.updateAppointment(rescheduleTarget.value, {
    appointment_date: form.value.appointment_date,
    status: 'rescheduled',
  })
  showReschedule.value = false
  rescheduleTarget.value = null
  form.value.appointment_date = ''
  slots.value = []
  selectedSlot.value = null
}

// Availability editor actions
const addAvailabilityRow = () => {
  availability.value.push({ day_of_week: 'Mon', start: '09:00', end: '17:00' })
}
const removeAvailabilityRow = (idx: number) => {
  availability.value.splice(idx, 1)
}
const saveAvailability = async () => {
  if (!userStore.profile?.id) return
  // Prevent saving invalid rows
  const invalid = availability.value.some(a => isInvalid(a))
  if (invalid) {
    return
  }
  await apptStore.setAvailability(userStore.profile.id, {
    available_hours: availability.value.map(a => ({ day_of_week: a.day_of_week, start: a.start, end: a.end })),
  })
  showAvailability.value = false
  // Update local profile cache so UI persists after reloads
  ;(userStore.profile as any).available_hours = availability.value.map(a => ({ day_of_week: a.day_of_week, start: a.start, end: a.end }))
}

// Simple validation helper
const isInvalid = (row: { start: string; end: string }) => {
  if (!row.start || !row.end) return true
  return row.end <= row.start
}

// Helpers for UI constraints
const todayISODate = computed(() => {
  const t = new Date()
  const yyyy = t.getFullYear()
  const mm = String(t.getMonth() + 1).padStart(2, '0')
  const dd = String(t.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
})

const isSlotInPast = (slot: TimeSlot) => {
  return new Date(slot.start_time).getTime() < Date.now()
}

// QR modal state & opener
const qrOpen = ref(false)
const qrPetId = ref<string>('')
const openQR = async (petId: string) => {
  qrPetId.value = petId
  try {
    await qrStore.getOrCreateForPet(petId)
  } catch (_) {}
  qrOpen.value = true
}
</script> 

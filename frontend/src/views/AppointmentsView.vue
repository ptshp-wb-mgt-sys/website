<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">{{ pageTitle }}</h1>
        <p class="text-gray-600 mt-1">{{ pageSubtitle }}</p>
      </div>
      <Button v-if="userStore.isClient" @click="showBook = true">
        <Plus class="w-4 h-4 mr-2" />
        Book Appointment
      </Button>
      <Button v-else-if="userStore.isVeterinarian">
        <Settings class="w-4 h-4 mr-2" />
        Manage Availability
      </Button>
    </div>

    <!-- CLIENT VIEW: Appointment Booking -->
    <template v-if="userStore.isClient">
      <!-- Quick Book Section -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Book New Appointment</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Select Pet</label>
            <select v-model="form.pet_id" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
              <option value="" disabled>Select pet</option>
              <option v-for="pet in pets" :key="pet.id" :value="pet.id">{{ pet.name }} ({{ pet.breed }})</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Select Veterinarian</label>
            <select v-model="form.veterinarian_id" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
              <option value="" disabled>Select veterinarian</option>
              <option v-for="vet in veterinarians" :key="vet.id" :value="vet.id">{{ vet.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Reason for Visit</label>
            <select v-model="form.reason" class="w-full px-3 py-2 border border-gray-300 rounded-lg">
              <option value="" disabled>Select reason</option>
              <option>Routine Checkup</option>
              <option>Vaccination</option>
              <option>Emergency</option>
              <option>Follow-up</option>
            </select>
          </div>
        </div>
        <div class="mt-4 flex space-x-2">
          <input v-model="date" type="date" class="px-3 py-2 border border-gray-300 rounded-lg" />
          <Button @click="loadSlots">Find Available Times</Button>
        </div>
        <div v-if="slots.length" class="grid grid-cols-2 md:grid-cols-4 gap-2 mt-4">
          <Button v-for="s in slots" :key="s.start_time" variant="outline" @click="selectSlot(s)">
            {{ new Date(s.start_time).toLocaleTimeString() }}
          </Button>
        </div>
        <div class="mt-4" v-if="selectedSlot">
          <Button @click="book">Confirm Booking</Button>
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
                  <h3 class="font-semibold text-rich-black">{{ appt.reason }}</h3>
                  <p class="text-sm text-gray-600">{{ new Date(appt.appointment_date).toLocaleString() }}</p>
                  <p class="text-xs" :class="appt.status === 'confirmed' ? 'text-green-600' : 'text-gray-600'">{{ appt.status }}</p>
                </div>
                <div class="flex space-x-2">
                  <Button variant="outline" size="sm">Reschedule</Button>
                  <Button variant="ghost" size="sm" class="text-red-600">Cancel</Button>
                </div>
              </div>
            </Card>
          </div>
        </section>

        <section>
          <h2 class="text-xl font-semibold text-rich-black mb-4">Past Appointments</h2>
          <div class="space-y-4">
            <Card class="p-6">
              <div class="flex items-center justify-between">
                <div class="space-y-1">
                  <h3 class="font-semibold text-rich-black">Luna - Dental Cleaning</h3>
                  <p class="text-sm text-gray-600">Dr. Sarah Johnson • Jan 15, 2024</p>
                  <p class="text-xs text-green-600">✓ Completed</p>
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
    </template>

    <!-- VETERINARIAN VIEW: Schedule Management -->
    <template v-if="userStore.isVeterinarian">
      <!-- Today's Schedule -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Today's Schedule - {{ currentDate }}</h2>
          <div class="flex space-x-2">
            <Button variant="outline" size="sm">
              <Calendar class="w-4 h-4 mr-2" />
              View Calendar
            </Button>
            <Button variant="outline" size="sm">
              <Clock class="w-4 h-4 mr-2" />
              Set Availability
            </Button>
          </div>
        </div>
        
        <div class="space-y-3">
          <div class="flex items-center justify-between p-4 bg-green-50 border border-green-200 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">10:00 AM - Buddy (Golden Retriever)</p>
              <p class="text-sm text-gray-600">Annual checkup • John Smith • +1 (555) 123-4567</p>
              <p class="text-xs text-green-600">✓ Patient checked in</p>
            </div>
            <div class="flex space-x-2">
              <Button size="sm">Start Visit</Button>
              <Button variant="outline" size="sm">View History</Button>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-4 bg-blue-50 border border-blue-200 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">11:30 AM - Max (German Shepherd)</p>
              <p class="text-sm text-gray-600">Follow-up • Sarah Johnson • +1 (555) 987-6543</p>
              <p class="text-xs text-blue-600">⏱ Next appointment</p>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm">Prepare</Button>
              <Button variant="ghost" size="sm">Contact Owner</Button>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">2:00 PM - Luna (Persian Cat)</p>
              <p class="text-sm text-gray-600">Vaccination • Mike Davis • +1 (555) 456-7890</p>
              <p class="text-xs text-gray-600">Scheduled</p>
            </div>
            <div class="flex space-x-2">
              <Button variant="outline" size="sm">View Details</Button>
              <Button variant="ghost" size="sm">Reschedule</Button>
            </div>
          </div>
        </div>
      </Card>

      <!-- Availability Management -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">This Week's Availability</h2>
        <div class="grid grid-cols-7 gap-2">
          <div v-for="day in weekDays" :key="day.name" class="text-center">
            <p class="text-sm font-medium text-gray-700 mb-2">{{ day.name }}</p>
            <div class="space-y-1">
              <div v-for="slot in day.slots" :key="slot" 
                   :class="slot.includes('Booked') ? 'bg-red-100 text-red-700' : 'bg-green-100 text-green-700'"
                   class="text-xs p-1 rounded">
                {{ slot }}
              </div>
            </div>
          </div>
        </div>
        <Button variant="outline" class="mt-4">Update Availability</Button>
      </Card>

      <!-- Recent Patients -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Recent Patients</h2>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Charlie - Routine Checkup</p>
              <p class="text-sm text-gray-600">Yesterday • All healthy ✓</p>
            </div>
            <Button variant="ghost" size="sm">Add Follow-up</Button>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">Bella - Vaccination</p>
              <p class="text-sm text-gray-600">2 days ago • Vaccines up to date</p>
            </div>
            <Button variant="ghost" size="sm">View Records</Button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { usePetsStore } from '@/stores/pets'
import { useAppointmentsStore, type CreateAppointmentRequest, type TimeSlot } from '@/stores/appointments'
import { 
  Plus, Settings, Calendar, Clock, Users, TrendingUp, AlertCircle 
} from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const userStore = useUserStore()
const petsStore = usePetsStore()
const apptStore = useAppointmentsStore()

// Booking UI state
const showBook = ref(false)
const veterinarians = ref<{ id: string; name: string }[]>([])
const slots = ref<TimeSlot[]>([])
const selectedSlot = ref<TimeSlot | null>(null)
const date = ref<string>('')
const form = ref<CreateAppointmentRequest>({ veterinarian_id: '', pet_id: '', appointment_date: '', reason: '' })

onMounted(async () => {
  if (userStore.isClient) {
    if (petsStore.pets.length === 0) await petsStore.fetchPets()
    await apptStore.fetchAppointments()
    const vets = await apptStore.listVeterinarians()
    veterinarians.value = vets.map(v => ({ id: v.id, name: v.name }))
  }
})

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

/**
 * Mock data for veterinarian schedule
 */
const weekDays = [
  {
    name: 'Mon',
    slots: ['9:00 AM', '10:00 AM - Booked', '11:00 AM', '2:00 PM - Booked']
  },
  {
    name: 'Tue',
    slots: ['9:00 AM - Booked', '10:00 AM', '11:00 AM - Booked', '2:00 PM']
  },
  {
    name: 'Wed',
    slots: ['9:00 AM', '10:00 AM', '11:00 AM', '2:00 PM']
  },
  {
    name: 'Thu',
    slots: ['9:00 AM - Booked', '10:00 AM - Booked', '11:00 AM', '2:00 PM']
  },
  {
    name: 'Fri',
    slots: ['9:00 AM', '10:00 AM', '11:00 AM - Booked', '2:00 PM - Booked']
  },
  {
    name: 'Sat',
    slots: ['9:00 AM - Booked', '10:00 AM', 'Closed', 'Closed']
  },
  {
    name: 'Sun',
    slots: ['Closed', 'Closed', 'Closed', 'Closed']
  }
]

// Derived lists
const upcomingAppointments = computed(() => apptStore.upcomingAppointments)

// Pets for select
const pets = computed(() => petsStore.pets)

/**
 * Load available time slots for the selected vet and date.
 */
const loadSlots = async () => {
  if (!form.value.veterinarian_id || !date.value) return
  const iso = new Date(date.value + 'T00:00:00').toISOString()
  slots.value = await apptStore.getAvailableSlots(form.value.veterinarian_id, iso)
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
    reason: form.value.reason,
  })
  // Reset UI
  showBook.value = false
  form.value = { veterinarian_id: '', pet_id: '', appointment_date: '', reason: '' }
  date.value = ''
  slots.value = []
  selectedSlot.value = null
}
</script> 

<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">Veterinarian Dashboard</h1>
        <p class="text-gray-600 mt-1">Manage your practice and patients</p>
      </div>
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
            <p class="text-2xl font-bold text-rich-black">{{ todaysCount }}</p>
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
            <p class="text-2xl font-bold text-rich-black">{{ activePatientsCount }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <FileText class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Visits This Week</p>
            <p class="text-2xl font-bold text-rich-black">{{ visitsThisWeek }}</p>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="flex items-center">
          <div class="w-12 h-12 bg-aquamarine-100 rounded-lg flex items-center justify-center mr-4">
            <Coins class="w-6 h-6 text-aquamarine" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-600">Revenue This Month</p>
            <p class="text-2xl font-bold text-rich-black">₱{{ monthlyRevenue.toFixed(2) }}</p>
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
          <Button variant="ghost" size="sm" @click="goToAvailability">Manage Availability</Button>
        </div>
        <div class="space-y-3">
          <div v-if="todaysList.length === 0" class="text-sm text-gray-600">No appointments today.</div>
          <div v-for="appt in todaysList" :key="appt.id" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="space-y-1">
              <p class="font-medium text-rich-black">{{ new Date(appt.appointment_date).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' }) }} — {{ appt.reason }} • <span class="text-gray-600 text-sm">{{ petLabel(appt.pet_id) }}</span></p>
              <p class="text-xs capitalize" :class="appt.status === 'confirmed' ? 'text-green-600' : 'text-gray-600'">{{ appt.status }}</p>
            </div>
            <Button variant="ghost" size="sm" @click="goToPet(appt.pet_id)">Start Visit</Button>
          </div>
        </div>
      </Card>

      <!-- Recent Patients -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Recent Patients</h2>
          <Button variant="ghost" size="sm" @click="goToPatients">View All Patients</Button>
        </div>
        <div class="space-y-3">
          <div v-if="recentPatients.length === 0" class="text-sm text-gray-600">No recent patients.</div>
          <div v-for="p in recentPatients" :key="p.petId" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-aquamarine-100 rounded-full flex items-center justify-center">
                <Heart class="w-5 h-5 text-aquamarine" />
              </div>
              <div>
                <p class="font-medium text-rich-black">{{ p.name || 'Pet' }}</p>
                <p class="text-sm text-gray-600">Last visit: {{ p.lastVisit }}</p>
              </div>
            </div>
            <Button variant="ghost" size="sm" @click="goToPet(p.petId)">View Records</Button>
          </div>
        </div>
      </Card>

      <!-- Product Sales -->
      <Card class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-lg font-semibold text-rich-black">Product Sales</h2>
          <Button variant="ghost" size="sm" @click="goToProducts">Manage Products</Button>
        </div>
        <div class="space-y-3">
          <div v-if="topSales.length === 0" class="text-sm text-gray-600">No sales in the last 7 days.</div>
          <div v-for="s in topSales" :key="s.productId" class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <p class="font-medium text-rich-black">{{ s.name || s.productId }}</p>
              <p class="text-sm text-gray-600">{{ s.quantity }} units sold this week</p>
              <p class="text-xs text-green-600">₱{{ s.revenue.toFixed(2) }} revenue</p>
            </div>
            <Button variant="ghost" size="sm" @click="goToProducts">Details</Button>
          </div>
        </div>
      </Card>

      <!-- Quick Actions -->
      <Card class="p-6">
        <h2 class="text-lg font-semibold text-rich-black mb-4">Quick Actions</h2>
        <div class="grid grid-cols-2 gap-3">
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToNewRecord">
            <FileText class="w-6 h-6 mb-2" />
            <span class="text-sm">New Record</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToAvailability">
            <Calendar class="w-6 h-6 mb-2" />
            <span class="text-sm">Set Availability</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToProducts">
            <Package class="w-6 h-6 mb-2" />
            <span class="text-sm">Add Product</span>
          </Button>
          <Button variant="outline" class="flex flex-col items-center p-4 h-auto" @click="goToPatients">
            <Users class="w-6 h-6 mb-2" />
            <span class="text-sm">Patients</span>
          </Button>
        </div>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, Heart, Calendar, FileText, Users, Coins, Package } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import { useAppointmentsStore } from '@/stores/appointments'
import { useOrdersStore } from '@/stores/orders'
import { useProductsStore } from '@/stores/products'

const router = useRouter()
const apptStore = useAppointmentsStore()
const ordersStore = useOrdersStore()
const productsStore = useProductsStore()

onMounted(async () => {
  await apptStore.fetchAppointments({ force: true })
  await ordersStore.fetchOrders({ force: true })
  // Ensure we have orders before aggregating sales
  await productsStore.fetchProducts({ force: true })
})

// Stats
const todaysCount = computed(() => apptStore.todaysAppointments.length)
const todaysList = computed(() => apptStore.todaysAppointments)
const visitsThisWeek = computed(() => {
  const now = new Date()
  const start = new Date(now)
  const day = start.getDay() // 0 Sun
  start.setDate(start.getDate() - day)
  start.setHours(0, 0, 0, 0)
  const end = new Date(start)
  end.setDate(start.getDate() + 7)
  return apptStore.upcomingAppointments.filter(a => {
    const t = new Date(a.appointment_date).getTime()
    return t >= start.getTime() && t < end.getTime()
  }).length
})
const activePatientsCount = computed(() => {
  // Count unique pet IDs from all appointments (non-cancelled)
  const set = new Set<string>()
  apptStore.appointments
    ?.filter(a => a.status !== 'cancelled')
    .forEach(a => set.add(a.pet_id))
  return set.size
})
const monthlyRevenue = computed(() => ordersStore.monthlyRevenue || 0)

// Recent patients: last 2 completed visits by distinct pet
const recentPatients = computed(() => {
  const seen = new Set<string>()
  const out: Array<{ petId: string; name?: string; lastVisit: string }> = []
  const past = [...apptStore.pastAppointments]
    .filter(a => a.status === 'completed')
    .sort((a, b) => new Date(b.appointment_date).getTime() - new Date(a.appointment_date).getTime())
  for (const a of past) {
    if (seen.has(a.pet_id)) continue
    seen.add(a.pet_id)
    out.push({ petId: a.pet_id, name: undefined, lastVisit: new Date(a.appointment_date).toLocaleDateString() })
    if (out.length >= 2) break
  }
  return out
})

// Product sales aggregate for the last 7 days
const topSales = ref<Array<{ productId: string; name?: string; quantity: number; revenue: number }>>([])
onMounted(async () => {
  const end = new Date()
  const start = new Date()
  start.setDate(end.getDate() - 7)
  // ensure orders are loaded; fetchOrders above already ran, keep non-forced to avoid extra network
  await ordersStore.fetchOrders()
  const agg = await ordersStore.aggregateItemsBetween(start, end)
  const rows: Array<{ productId: string; name?: string; quantity: number; revenue: number }> = Object.entries(agg).map(([productId, v]) => ({ productId, quantity: v.quantity, revenue: v.revenue }))
  const map: Record<string, string> = {}
  for (const p of productsStore.products) map[p.id] = p.name
  rows.forEach(r => (r.name = map[r.productId]))
  topSales.value = rows.sort((a, b) => b.revenue - a.revenue).slice(0, 5)
})

// Helpers
function petLabel(petId: string): string {
  // Minimal label; dashboard doesn't have pet store here, so just show id tail as placeholder
  return 'Pet ' + String(petId).slice(0, 4)
}

// Navigation helpers
const goToNewRecord = () => {
  router.push({ name: 'patients', query: { action: 'new-record' } })
}
/** Open schedule page to edit availability. */
const goToAvailability = () => router.push({ name: 'my-schedule' })
/** Navigate to vet product catalog management. */
const goToProducts = () => router.push({ name: 'manage-products' })
/** Navigate to patients grid. */
const goToPatients = () => router.push({ name: 'patients' })
/** Open a pet profile by its id. */
const goToPet = (id: string) => router.push({ name: 'pet-profile', params: { id } })
// (Replaced scan QR quick action with Patients)
</script>

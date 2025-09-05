<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">Sales History</h1>
        <p class="text-gray-600 mt-1">Your orders received over time</p>
      </div>
    </div>

    <Card class="p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-rich-black">{{ rangeLabel }}</h2>
        <div class="flex items-center space-x-2">
          <Button size="sm" :variant="range.type === '7' ? 'default' : 'outline'" @click="setQuickRange('7')">7 days</Button>
          <Button size="sm" :variant="range.type === '30' ? 'default' : 'outline'" @click="setQuickRange('30')">30 days</Button>
          <Button size="sm" :variant="range.type === 'all' ? 'default' : 'outline'" @click="setQuickRange('all')">Lifetime</Button>
          <div class="flex items-center space-x-2">
            <select class="border rounded px-2 pr-8 py-1 text-sm" v-model.number="range.monthIndex" @change="onMonthYearChange">
              <option :value="null" disabled>Select month</option>
              <option v-for="(m, i) in monthNames" :key="i" :value="i">{{ m }}</option>
            </select>
            <select class="border rounded px-2 pr-8 py-1 text-sm" v-model.number="range.year" @change="onMonthYearChange">
              <option :value="null" disabled>Select year</option>
              <option v-for="y in availableYears" :key="y" :value="y">{{ y }}</option>
            </select>
            <Button size="sm" variant="outline" @click="clearMonthRange" v-if="range.type === 'month'">Clear</Button>
          </div>
        </div>
      </div>
      <div class="flex items-center justify-end mb-2">
        <span class="text-sm text-gray-600">{{ ordersCount }} orders • {{ formatPHP(totalRevenue) }}</span>
      </div>

      <div class="space-y-4">
        <div
          v-for="o in orderedRecent"
          :key="o.id"
          class="rounded-lg border border-gray-200"
        >
          <div class="px-4 py-3 flex items-center justify-between bg-gray-50 rounded-t-lg">
            <div>
              <p class="font-medium text-rich-black">Order #{{ shortId(o.id) }}</p>
              <p class="text-xs text-gray-600">{{ formatDate(o.created_at) }} • Buyer: {{ ownerLabelById[o.client_id]?.name || 'Client' }} • Status: {{ o.status }} • Payment: {{ o.payment_status }}</p>
            </div>
            <div class="text-right">
              <p class="font-semibold text-rich-black">
                <span class="text-sm text-gray-600 mr-2">Total</span>
                <span class="text-2xl">{{ formatPHP(o.total_amount) }}</span>
              </p>
            </div>
          </div>
          <div class="px-4 py-3 divide-y">
            <div class="py-2 text-xs text-gray-700 space-y-1">
              <p v-if="ownerLabelById[o.client_id]?.email">Email: {{ ownerLabelById[o.client_id]?.email }}</p>
              <p v-if="ownerLabelById[o.client_id]?.phone">Phone: {{ ownerLabelById[o.client_id]?.phone }}</p>
              <p v-if="ownerLabelById[o.client_id]?.address">Address: {{ ownerLabelById[o.client_id]?.address }}</p>
            </div>
            <div
              v-for="it in (itemsByOrderId[o.id] || [])"
              :key="it.id"
              class="py-2 flex items-center justify-between"
            >
              <div class="truncate pr-4">
                <p class="text-sm text-rich-black truncate">{{ productNameById[it.product_id] || it.product_id }}</p>
                <p class="text-xs text-gray-600">Qty {{ it.quantity }} × {{ formatPHP(it.unit_price) }}</p>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium">{{ formatPHP(it.total_price) }}</p>
              </div>
            </div>
            <div v-if="(itemsByOrderId[o.id] || []).length === 0" class="text-sm text-gray-600 py-2">No items loaded.</div>
            
          </div>
        </div>

        <div v-if="orderedRecent.length === 0" class="text-sm text-gray-600 py-2">No sales in the last 30 days.</div>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'
import { useOrdersStore } from '@/stores/orders'
import { useProductsStore } from '@/stores/products'
import { useAuthStore } from '@/stores/auth'
import { formatPHP } from '@/lib/utils'

const ordersStore = useOrdersStore()
const productsStore = useProductsStore()
const auth = useAuthStore()

const itemsByOrderId = ref<Record<string, import('@/stores/orders').OrderItem[]>>({})
const productNameById = ref<Record<string, string>>({})
const ownerLabelById = ref<Record<string, { name: string; email?: string; phone?: string; address?: string }>>({})

/** Load orders then hydrate items, product names, and buyer labels. */
onMounted(async () => {
  await ordersStore.fetchOrders({ force: true })
  await preloadOrderDetails()
})

/** Selected range; default 7 days. Uses separate month/year controls. */
const range = ref<{ type: '7' | '30' | 'all' | 'month'; monthIndex: number | null; year: number | null }>({ type: '7', monthIndex: null, year: null })

/** Friendly label for the current range. */
const rangeLabel = computed(() => {
  if (range.value.type === '7') return 'Recent (7 days)'
  if (range.value.type === '30') return 'Recent (30 days)'
  if (range.value.type === 'month' && range.value.monthIndex != null && range.value.year != null) {
    return `${monthNames[range.value.monthIndex]} ${range.value.year}`
  }
  return 'All time'
})

/** Quick set 7/30/all; clears month. */
const setQuickRange = (type: '7' | '30' | 'all') => { range.value = { type, monthIndex: null, year: null } }
/**
 * When either month or year changes, switch to month mode when both are set.
 */
const onMonthYearChange = () => { if (range.value.monthIndex != null && range.value.year != null) range.value.type = 'month' }
/** Clear month/year selection and revert to 7 days. */
const clearMonthRange = () => { range.value = { type: '7', monthIndex: null, year: null } }

/** Orders filtered by range (incoming sales for this vet). */
const recentOrders = computed(() => {
  const now = Date.now()
  if (range.value.type === 'month' && range.value.monthIndex != null && range.value.year != null) {
    const start = new Date(range.value.year, range.value.monthIndex, 1)
    const end = new Date(range.value.year, range.value.monthIndex + 1, 1)
    return ordersStore.orders.filter(o => {
      const t = new Date(o.created_at).getTime()
      return t >= start.getTime() && t < end.getTime() && o.status !== 'cancelled'
    })
  }
  const cutoff = range.value.type === '7' ? now - 7 * 24 * 60 * 60 * 1000
    : range.value.type === '30' ? now - 30 * 24 * 60 * 60 * 1000
    : -Infinity
  return ordersStore.orders.filter(o => new Date(o.created_at).getTime() >= cutoff && o.status !== 'cancelled')
})

const orderedRecent = computed(() => [...recentOrders.value].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()))
const ordersCount = computed(() => recentOrders.value.length)
const totalRevenue = computed(() => recentOrders.value.reduce((s, o) => s + (o.total_amount || 0), 0))

/** Short id helper. */
const shortId = (id: string) => id.slice(0, 8)
/** Format date helper. */
const formatDate = (iso: string) => new Date(iso).toLocaleString()

/** Month names for month-year filter dropdown. */
const monthNames = ['January','February','March','April','May','June','July','August','September','October','November','December']

/**
 * Available years based on existing orders; falls back to current year.
 */
const availableYears = computed(() => {
  const set = new Set<number>()
  for (const o of ordersStore.orders) {
    const y = new Date(o.created_at).getFullYear()
    if (!Number.isNaN(y)) set.add(y)
  }
  const arr = Array.from(set)
  if (!arr.length) return [new Date().getFullYear()]
  return arr.sort((a, b) => b - a)
})

/**
 * Preload items for each order and resolve product names and buyer labels.
 */
const preloadOrderDetails = async () => {
  const productIds = new Set<string>()
  const ownerIds = new Set<string>()
  for (const o of recentOrders.value) {
    try {
      const items = await ordersStore.fetchOrderItems(o.id)
      itemsByOrderId.value[o.id] = items
      items.forEach(it => productIds.add(it.product_id))
      ownerIds.add(o.client_id)
    } catch (_) {}
  }
  // Seed product names from products store
  for (const p of productsStore.products) productNameById.value[p.id] = p.name
  // Fetch any missing product names
  await Promise.all(
    Array.from(productIds).filter(id => !productNameById.value[id]).map(id => ensureProductName(id)),
  )
  // Fetch buyer labels for owners
  await Promise.all(Array.from(ownerIds).map(id => ensureOwnerLabel(id)))
}

// Rehydrate details as range changes to include newly visible orders
watch(recentOrders, async () => {
  await preloadOrderDetails()
})

/** Ensure a product name is available; fetches /products/:id if missing. */
const ensureProductName = async (productId: string) => {
  if (productNameById.value[productId]) return productNameById.value[productId]
  try {
    if (!auth.session?.access_token) return productId
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products/${encodeURIComponent(productId)}`, {
      headers: { Authorization: `Bearer ${auth.session.access_token}`, 'Content-Type': 'application/json' },
    })
    if (!res.ok) return productId
    const body = await res.json()
    const product = body.data || body
    if (product?.name) productNameById.value[productId] = product.name
    return productNameById.value[productId] || productId
  } catch (_) {
    return productId
  }
}

/** Ensure we have buyer label for display; calls /owners/:id/label and caches fields. */
const ensureOwnerLabel = async (ownerId: string) => {
  const cached = ownerLabelById.value[ownerId]
  if (cached) return cached
  try {
    if (!auth.session?.access_token) return null
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/owners/${encodeURIComponent(ownerId)}/label`, {
      headers: { Authorization: `Bearer ${auth.session.access_token}`, 'Content-Type': 'application/json' },
    })
    if (!res.ok) return null
    const body = await res.json()
    const label = body.data || body
    if (label?.id) {
      ownerLabelById.value[ownerId] = {
        name: label.name,
        email: label.email,
        phone: label.phone,
        address: label.address,
      }
      return ownerLabelById.value[ownerId]
    }
  } catch (_) {}
  return null
}
</script>

<style scoped>
</style>



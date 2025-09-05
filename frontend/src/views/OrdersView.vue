<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">Order History</h1>
        <p class="text-gray-600 mt-1">Your orders over time</p>
      </div>
    </div>

    <Card class="p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-rich-black">Recent (30 days)</h2>
        <span class="text-sm text-gray-600">{{ orderedRecent.length }} orders • {{ formatMoney(recentRevenue) }}</span>
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
              <p class="text-xs text-gray-600">{{ formatDate(o.created_at) }} • Status: {{ o.status }} • Payment: {{ o.payment_status }}</p>
            </div>
            <div class="text-right">
              <p class="font-semibold">{{ formatMoney(o.total_amount) }}</p>
            </div>
          </div>
          <div class="px-4 py-3 divide-y">
            <div
              v-for="it in (itemsByOrderId[o.id] || [])"
              :key="it.id"
              class="py-2 flex items-center justify-between"
            >
              <div class="truncate pr-4">
                <p class="text-sm text-rich-black truncate">{{ productNameById[it.product_id] || it.product_id }}</p>
                <p class="text-xs text-gray-600">Qty {{ it.quantity }} × {{ formatMoney(it.unit_price) }}</p>
              </div>
              <div class="text-right">
                <p class="text-sm font-medium">{{ formatMoney(it.total_price) }}</p>
              </div>
            </div>
            <div v-if="(itemsByOrderId[o.id] || []).length === 0" class="text-sm text-gray-600 py-2">No items loaded.</div>
            <div class="pt-3 flex items-center justify-end">
              <div class="text-right">
                <p class="text-sm text-gray-600">Order total</p>
                <p class="text-base font-semibold">{{ formatMoney(o.total_amount) }}</p>
              </div>
            </div>
          </div>
        </div>

        <div v-if="orderedRecent.length === 0" class="text-sm text-gray-600 py-2">No orders in the last 30 days.</div>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Card from '@/components/ui/Card.vue'
import { useOrdersStore } from '@/stores/orders'
import { useProductsStore } from '@/stores/products'
import { useAuthStore } from '@/stores/auth'

const ordersStore = useOrdersStore()
const productsStore = useProductsStore()
const auth = useAuthStore()

const itemsByOrderId = ref<Record<string, import('@/stores/orders').OrderItem[]>>({})
const productNameById = ref<Record<string, string>>({})

/**
 * Ensure we have orders loaded to display history.
 */
onMounted(async () => {
  if (ordersStore.orders.length === 0) {
    await ordersStore.fetchOrders()
  }
  // Preload items for recent orders and resolve product names.
  await preloadRecentOrderDetails()
})

/**
 * Orders from the last 30 days.
 */
const recentOrders = computed(() => {
  const THIRTY_DAYS_MS = 30 * 24 * 60 * 60 * 1000
  const cutoffMs = Date.now() - THIRTY_DAYS_MS
  return ordersStore.orders.filter(o => {
    const createdMs = new Date(o.created_at).getTime()
    return !Number.isNaN(createdMs) && createdMs >= cutoffMs
  })
})

/**
 * Recent orders sorted by most recent first.
 */
const orderedRecent = computed(() =>
  [...recentOrders.value].sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()),
)

/**
 * Total revenue for recent orders (30 days).
 */
const recentRevenue = computed(() => recentOrders.value.reduce((sum, o) => sum + (o.total_amount || 0), 0))

/**
 * Format date into a compact readable label.
 */
const formatDate = (iso: string) => new Date(iso).toLocaleString()

/**
 * Format number as currency. Keep it chill.
 */
const formatMoney = (amount: number) => `$${(amount ?? 0).toFixed(2)}`

/**
 * Shorten IDs for display.
 */
const shortId = (id: string) => id.slice(0, 8)

/**
 * Preload items for all recent orders and resolve product names.
 */
const preloadRecentOrderDetails = async () => {
  const uniqueProductIds = new Set<string>()
  for (const o of recentOrders.value) {
    try {
      const items = await ordersStore.fetchOrderItems(o.id)
      itemsByOrderId.value[o.id] = items
      for (const it of items) uniqueProductIds.add(it.product_id)
    } catch (_) {}
  }
  // Seed names from any products already in store
  for (const p of productsStore.products) {
    productNameById.value[p.id] = p.name
  }
  // Fetch missing product names
  await Promise.all(
    Array.from(uniqueProductIds)
      .filter(pid => !productNameById.value[pid])
      .map(pid => ensureProductName(pid)),
  )
}

/**
 * Ensure a product name is cached; fetches product details if needed.
 */
const ensureProductName = async (productId: string) => {
  if (productNameById.value[productId]) return productNameById.value[productId]
  try {
    if (!auth.session?.access_token) return productId
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products/${encodeURIComponent(productId)}`, {
      headers: {
        Authorization: `Bearer ${auth.session.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    if (!res.ok) return productId
    const body = await res.json()
    const product = (body.data || body)
    if (product && product.name) {
      productNameById.value[productId] = product.name as string
      return product.name as string
    }
  } catch (_) {}
  return productId
}
</script>

<style scoped>
</style>



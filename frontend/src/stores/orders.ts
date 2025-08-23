import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

export interface OrderItem {
  id: string
  order_id: string
  product_id: string
  quantity: number
  unit_price: number
  total_price: number
}

export interface Order {
  id: string
  client_id: string
  veterinarian_id: string
  total_amount: number
  status: string
  payment_status: string
  created_at: string
  updated_at: string
}

/**
 * Store for vet/client orders. Fetches lightweight order lists and, on demand,
 * loads order details (items) to support simple revenue/sales summaries.
 */
export const useOrdersStore = defineStore('orders', () => {
  const orders = ref<Order[]>([])
  const orderItemsByOrderId = ref<Record<string, OrderItem[]>>({})
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAt = ref<number | null>(null)

  const auth = useAuthStore()

  /**
   * Load the authenticated user's orders from the API.
   * Skips the network within a TTL unless `force` is set.
   */
  const fetchOrders = async (options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 60 * 1000
    if (!force && lastFetchedAt.value && Date.now() - lastFetchedAt.value < ttlMs && orders.value.length) return
    if (!auth.session?.access_token) {
      error.value = 'No authentication token'
      return
    }
    loading.value = true
    error.value = null
    try {
      const res = await fetch('http://localhost:3000/api/v1/orders', {
        headers: {
          Authorization: `Bearer ${auth.session.access_token}`,
          'Content-Type': 'application/json',
        },
      })
      if (!res.ok) throw new Error(res.statusText)
      const body = await res.json()
      orders.value = (body.data || body) as Order[]
      lastFetchedAt.value = Date.now()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch orders'
      console.error('Error fetching orders', e)
    } finally {
      loading.value = false
    }
  }

  /**
   * Load line items for a specific order.
   * Returns cached items when available.
   */
  const fetchOrderItems = async (orderId: string) => {
    if (!auth.session?.access_token) throw new Error('No authentication token')
    if (orderItemsByOrderId.value[orderId]) return orderItemsByOrderId.value[orderId]
    const res = await fetch(`http://localhost:3000/api/v1/orders/${encodeURIComponent(orderId)}`, {
      headers: {
        Authorization: `Bearer ${auth.session.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    if (!res.ok) throw new Error(res.statusText)
    const body = await res.json()
    const items = (body.items || []) as OrderItem[]
    orderItemsByOrderId.value[orderId] = items
    return items
  }

  /**
   * Compute revenue for orders created in a given month (local time).
   * Cancelled orders are excluded.
   */
  const revenueForMonth = (dateInMonth: Date) => {
    const y = dateInMonth.getFullYear()
    const m = dateInMonth.getMonth()
    return orders.value
      .filter(o => {
        const d = new Date(o.created_at)
        return d.getFullYear() === y && d.getMonth() === m && o.status !== 'cancelled'
      })
      .reduce((sum, o) => sum + (o.total_amount || 0), 0)
  }

  /**
   * Aggregate items across orders within [start, end] time range.
   * Returns totals per product id with quantity and revenue.
   */
  const aggregateItemsBetween = async (start: Date, end: Date) => {
    const startMs = start.getTime()
    const endMs = end.getTime()
    const relevant = orders.value.filter(o => {
      const t = new Date(o.created_at).getTime()
      return t >= startMs && t <= endMs && o.status !== 'cancelled'
    })
    const productIdToTotals: Record<string, { quantity: number; revenue: number }> = {}
    for (const o of relevant) {
      try {
        const items = await fetchOrderItems(o.id)
        for (const it of items) {
          const bucket = productIdToTotals[it.product_id] || { quantity: 0, revenue: 0 }
          bucket.quantity += it.quantity
          bucket.revenue += it.total_price
          productIdToTotals[it.product_id] = bucket
        }
      } catch (_) {}
    }
    return productIdToTotals
  }

  const monthlyRevenue = computed(() => revenueForMonth(new Date()))

  return {
    // state
    orders,
    orderItemsByOrderId,
    loading,
    error,
    lastFetchedAt,
    // actions
    fetchOrders,
    fetchOrderItems,
    // helpers
    revenueForMonth,
    aggregateItemsBetween,
    // computed
    monthlyRevenue,
  }
})



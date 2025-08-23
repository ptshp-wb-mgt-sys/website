import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

export interface Product {
  id: string
  veterinarian_id: string
  name: string
  description?: string
  category?: string
  price: number
  stock_quantity: number
  is_active: boolean
  created_at: string
  updated_at: string
}

/** Simple product store to list vet products for dashboard widgets. */
export const useProductsStore = defineStore('products', () => {
  const products = ref<Product[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAt = ref<number | null>(null)

  const auth = useAuthStore()

  /**
   * Fetch products. Vets see their own catalog; optionally scope by vet via `veterinarianId`.
   */
  const fetchProducts = async (opts?: { veterinarianId?: string; force?: boolean; ttlMs?: number }) => {
    const ttlMs = opts?.ttlMs ?? 60 * 1000
    const force = opts?.force === true
    if (!force && lastFetchedAt.value && Date.now() - lastFetchedAt.value < ttlMs && products.value.length) return
    if (!auth.session?.access_token) {
      error.value = 'No authentication token'
      return
    }
    loading.value = true
    error.value = null
    try {
      const qs = opts?.veterinarianId ? `/veterinarians/${encodeURIComponent(opts.veterinarianId)}/products` : '/products'
      const res = await fetch(`http://localhost:3000/api/v1${qs}`, {
        headers: {
          Authorization: `Bearer ${auth.session.access_token}`,
          'Content-Type': 'application/json',
        },
      })
      if (!res.ok) throw new Error(res.statusText)
      const body = await res.json()
      products.value = (body.data || body) as Product[]
      lastFetchedAt.value = Date.now()
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch products'
      console.error('Error fetching products', e)
    } finally {
      loading.value = false
    }
  }

  const activeProducts = computed(() => products.value.filter(p => p.is_active))

  return { products, activeProducts, loading, error, lastFetchedAt, fetchProducts }
})



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
  images?: string[]
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
  const fetchProducts = async (opts?: { veterinarianId?: string; force?: boolean; ttlMs?: number; category?: string; search?: string }) => {
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
      const base = opts?.veterinarianId ? `/veterinarians/${encodeURIComponent(opts.veterinarianId)}/products` : '/products'
      const params = new URLSearchParams()
      if (opts?.category) params.set('category', opts.category)
      if (opts?.search) params.set('search', opts.search)
      const qs = params.toString() ? `${base}?${params.toString()}` : base
      const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1${qs}`, {
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

  /** Create a product (vet/admin). */
  const createProduct = async (payload: Partial<Product>) => {
    if (!auth.session?.access_token) throw new Error('No authentication token')
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${auth.session.access_token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    })
    if (!res.ok) throw new Error((await res.json().catch(() => ({}))).message || res.statusText)
    const data = await res.json()
    const created = (data.data || data) as Product
    products.value.unshift(created)
    return created
  }

  /** Update a product by id (vet/admin). */
  const updateProduct = async (id: string, updates: Partial<Product>) => {
    if (!auth.session?.access_token) throw new Error('No authentication token')
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products/${encodeURIComponent(id)}`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${auth.session.access_token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify(updates),
    })
    if (!res.ok) throw new Error((await res.json().catch(() => ({}))).message || res.statusText)
    const data = await res.json()
    const updated = (data.data || data) as Product
    const idx = products.value.findIndex(p => p.id === id)
    if (idx !== -1) products.value[idx] = updated
    return updated
  }

  /** Soft delete (deactivate) product. */
  const deactivateProduct = async (id: string) => {
    if (!auth.session?.access_token) throw new Error('No authentication token')
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products/${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: { 'Authorization': `Bearer ${auth.session.access_token}` },
    })
    if (!res.ok) throw new Error(res.statusText)
    const idx = products.value.findIndex(p => p.id === id)
    if (idx !== -1) products.value[idx].is_active = false
  }

  /** Update product stock quantity. */
  const updateStock = async (id: string, quantity: number) => {
    if (!auth.session?.access_token) throw new Error('No authentication token')
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/products/${encodeURIComponent(id)}/stock`, {
      method: 'PUT',
      headers: { 'Authorization': `Bearer ${auth.session.access_token}`, 'Content-Type': 'application/json' },
      body: JSON.stringify({ quantity }),
    })
    if (!res.ok) throw new Error(res.statusText)
    const item = products.value.find(p => p.id === id)
    if (item) item.stock_quantity = quantity
  }

  return { products, activeProducts, loading, error, lastFetchedAt, fetchProducts, createProduct, updateProduct, deactivateProduct, updateStock }
})



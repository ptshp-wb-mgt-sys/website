import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface CartItem {
  /** Unique product id */
  productId: string
  /** Display name for quick rendering */
  name: string
  /** Unit price in PHP */
  price: number
  /** Quantity selected */
  quantity: number
  /** Optional small image URL */
  imageUrl?: string
}

/**
 * Lightweight cart store for pet owners. Pure client-side for now.
 * Persists to localStorage so the cart survives reloads.
 */
export const useCartStore = defineStore('cart', () => {
  const items = ref<CartItem[]>([])

  function loadFromStorage(): void {
    try {
      const raw = localStorage.getItem('petmgt:cart')
      if (raw) items.value = JSON.parse(raw)
    } catch {
      items.value = []
    }
  }

  function persist(): void {
    localStorage.setItem('petmgt:cart', JSON.stringify(items.value))
  }

  /** Add or bump quantity for a product. */
  function addItem(payload: Omit<CartItem, 'quantity'> & { quantity?: number }): void {
    const qty = payload.quantity ?? 1
    const existing = items.value.find(i => i.productId === payload.productId)
    if (existing) {
      existing.quantity += qty
    } else {
      items.value.push({ ...payload, quantity: qty })
    }
    persist()
  }

  /** Update quantity; removes item when qty <= 0. */
  function updateQuantity(productId: string, quantity: number): void {
    const idx = items.value.findIndex(i => i.productId === productId)
    if (idx === -1) return
    if (quantity <= 0) {
      items.value.splice(idx, 1)
    } else {
      items.value[idx].quantity = quantity
    }
    persist()
  }

  /** Remove by product id. */
  function removeItem(productId: string): void {
    items.value = items.value.filter(i => i.productId !== productId)
    persist()
  }

  /** Clear everything. */
  function clear(): void {
    items.value = []
    persist()
  }

  const itemCount = computed(() => items.value.reduce((sum, i) => sum + i.quantity, 0))
  const distinctCount = computed(() => items.value.length)
  const subtotal = computed(() => items.value.reduce((sum, i) => sum + i.price * i.quantity, 0))

  // initialize from storage on first use
  loadFromStorage()

  return { items, itemCount, distinctCount, subtotal, addItem, updateQuantity, removeItem, clear }
})



<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">Pet Products</h1>
      <div class="flex space-x-2">
        <Button v-if="userStore.isClient" variant="outline" @click="openCart = true">View Cart ({{ cart.distinctCount }})</Button>
        <Button v-if="userStore.isVeterinarian" @click="startCreate">Add Product</Button>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-2 mb-2 items-center">
      <Button v-if="userStore.isVeterinarian" :variant="mineOnly ? 'secondary' : 'ghost'" size="sm" @click="toggleMine">{{ mineOnly ? 'My Products ✔️' : 'My Products ❌' }}</Button>
      <Button v-if="userStore.isVeterinarian" :variant="showUnlisted ? 'secondary' : 'ghost'" size="sm" @click="toggleUnlisted">{{ showUnlisted ? 'Unlisted ✔️' : 'Unlisted ❌' }}</Button>
      <span class="text-gray-300">|</span>
      <Button :variant="category === '' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('')">All</Button>
      <Button :variant="category === 'Food' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('Food')">Food</Button>
      <Button :variant="category === 'Medicine' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('Medicine')">Medicine</Button>
      <Button :variant="category === 'Accessories' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('Accessories')">Accessories</Button>
      <Button :variant="category === 'Toys' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('Toys')">Toys</Button>
      <Button :variant="category === 'Other' ? 'secondary' : 'ghost'" size="sm" @click="setCategory('Other')">Other</Button>
      <input v-model="search" placeholder="Search" class="ml-auto border rounded px-3 py-2 w-full md:w-64" @input="onSearch" />
    </div>

    <!-- Products Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card v-for="p in filtered" :key="p.id" class="p-6 relative" :class="{ 'opacity-75 border-dashed border-gray-400': !p.is_active }">
        <!-- Unlisted badge for vets -->
        <div v-if="!p.is_active && userStore.isVeterinarian" class="absolute top-2 right-2 bg-gray-500 text-white text-xs px-2 py-1 rounded">
          Unlisted
        </div>
        <div class="space-y-4">
          <div class="aspect-square rounded-lg mb-4 overflow-hidden bg-gray-100 border border-gray-200 flex items-center justify-center">
            <img v-if="p.images && p.images.length" :src="p.images[0]" alt="" class="w-full h-full object-cover" />
            <svg v-else xmlns="http://www.w3.org/2000/svg" viewBox="0 0 128 128" class="w-12 h-12 text-aquamarine fill-current opacity-90">
              <!-- clear, simple paw print: 4 toe pads + big pad -->
              <circle cx="36" cy="36" r="12"/>
              <circle cx="64" cy="28" r="12"/>
              <circle cx="92" cy="36" r="12"/>
              <circle cx="24" cy="58" r="11"/>
              <path d="M64 60c-20 0-38 14-38 30 0 8 6 12 12 12 9 0 13-4 26-4s17 4 26 4c6 0 12-4 12-12 0-16-18-30-38-30z"/>
            </svg>
          </div>

          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">{{ p.name }}</h3>
            <p class="text-sm text-gray-600 line-clamp-2">{{ p.description }}</p>
            <!-- Vet info -->
            <p v-if="getVetInfo(p.veterinarian_id)" class="text-xs text-gray-500">
              <span class="font-medium">Dr. {{ getVetInfo(p.veterinarian_id)?.name }}</span>
              <span v-if="getVetInfo(p.veterinarian_id)?.clinic_address"> • {{ getVetInfo(p.veterinarian_id)?.clinic_address }}</span>
            </p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">{{ formatPriceWithMeasurement(p) }}</span>
              <span class="text-sm text-gray-500">
                {{ p.stock_quantity > 0
                  ? (p.stock_quantity < 5
                    ? `Low Stock (${p.stock_quantity})`
                    : `In Stock (${p.stock_quantity})`)
                  : 'Out of Stock' }}
              </span>
            </div>
          </div>

          <!-- Client view: Add to cart -->
          <div v-if="userStore.isClient" class="flex items-center gap-2">
            <input
              type="number"
              class="w-16 border rounded px-2 py-1"
              min="1"
              :max="Math.max(1, p.stock_quantity)"
              :value="getQty(p.id)"
              @input="onQtyInput(p.id, $event)"
            />
            <Button size="sm" class="flex-1" :disabled="p.stock_quantity <= 0" @click="addToCart(p)">Add to Cart</Button>
          </div>

          <!-- Vet view: Edit, Unlist/Relist on left, Delete on right -->
          <div v-if="userStore.isVeterinarian" class="flex items-center justify-between">
            <div class="flex gap-2">
              <Button size="sm" @click="startEdit(p)">Edit</Button>
              <Button variant="outline" size="sm" @click="toggleActive(p)">{{ p.is_active ? 'Unlist' : 'Relist' }}</Button>
            </div>
            <Button variant="ghost" size="sm" class="text-red-600 hover:text-red-700" @click="confirmDelete(p)">Delete</Button>
          </div>
        </div>
      </Card>
    </div>

    <!-- Cart modal -->
    <CartModal :open="openCart" @close="openCart = false" @checkout="applyCheckout" />

    <!-- Product modal -->
    <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="cancelEdit" />
      <div class="relative z-10 w-full max-w-lg rounded-lg bg-white shadow-lg p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-xl font-semibold text-rich-black">{{ editing ? 'Edit Product' : 'Add Product' }}</h3>
          <button class="text-gray-500 hover:text-gray-700" @click="cancelEdit">✕</button>
        </div>
        <div class="grid gap-4">
          <div>
            <label class="block text-sm font-medium mb-1">Product name</label>
            <input v-model="form.name" placeholder="e.g., Dog Food" class="border rounded px-3 py-2 w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Description</label>
            <textarea v-model="form.description" placeholder="Short description" class="border rounded px-3 py-2 w-full" />
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Photo</label>
            <div class="flex items-center gap-3">
              <input type="file" accept="image/*" class="border rounded px-3 py-2 w-full" disabled title="Coming soon" />
              <span class="text-xs text-gray-500">Coming soon</span>
            </div>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Category</label>
              <select v-model="form.category" class="border rounded px-3 py-2 w-full">
                <option value="">Select category</option>
                <option>Food</option>
                <option>Medicine</option>
                <option>Accessories</option>
                <option>Toys</option>
                <option>Other</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Price (₱)</label>
              <div class="flex items-center">
                <span class="inline-flex items-center px-3 py-2 border border-r-0 rounded-l bg-gray-50 text-gray-600">₱</span>
                <input type="number" min="0" step="0.01" v-model.number="form.price" placeholder="0.00" class="border rounded-r px-3 py-2 w-full" />
              </div>
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium mb-1">Qty in stock</label>
            <input type="number" min="0" v-model.number="form.stock_quantity" placeholder="0" class="border rounded px-3 py-2 w-full" />
          </div>
          <!-- Measurement fields for Food and Medicine -->
          <div v-if="needsMeasurement(form.category)" class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium mb-1">Content/Size</label>
              <input type="number" min="0" step="0.01" v-model.number="form.measurement_value" placeholder="e.g., 500" class="border rounded px-3 py-2 w-full" />
            </div>
            <div>
              <label class="block text-sm font-medium mb-1">Unit</label>
              <select v-model="form.measurement_unit" class="border rounded px-3 py-2 w-full">
                <option value="">Select unit</option>
                <option v-for="unit in measurementUnits[form.category] || []" :key="unit" :value="unit">{{ unit }}</option>
              </select>
            </div>
          </div>
        </div>
        <div class="flex gap-3 justify-end">
          <Button variant="outline" @click="cancelEdit">Cancel</Button>
          <Button :disabled="saving || !form.name || !form.price" @click="saveProduct">{{ saving ? 'Saving...' : 'Save' }}</Button>
        </div>
      </div>
    </div>

    <!-- Delete confirmation modal -->
    <div v-if="showDeleteModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/40" @click="cancelDelete" />
      <div class="relative z-10 w-full max-w-md rounded-lg bg-white shadow-lg p-6 space-y-4">
        <div class="flex items-center justify-between">
          <h3 class="text-xl font-semibold text-rich-black">Delete Product</h3>
          <button class="text-gray-500 hover:text-gray-700" @click="cancelDelete">✕</button>
        </div>
        <p class="text-gray-600">
          Are you sure you want to permanently delete <strong>{{ deleteTarget?.name }}</strong>? This action cannot be undone.
        </p>
        <div class="flex gap-3 justify-end">
          <Button variant="outline" @click="cancelDelete">Cancel</Button>
          <Button variant="ghost" class="bg-red-600 text-white hover:bg-red-700" :disabled="deleting" @click="deleteProduct">
            {{ deleting ? 'Deleting...' : 'Delete' }}
          </Button>
        </div>
      </div>
    </div>
  </div>
  
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import CartModal from '@/components/CartModal.vue'
import { useProductsStore, type Product } from '@/stores/products'
import { useCartStore } from '@/stores/cart'
import { useUserStore } from '@/stores/user'
import { useAuthStore } from '@/stores/auth'
import { useAppointmentsStore } from '@/stores/appointments'
import { formatPHP } from '@/lib/utils'

/** Products page supporting browsing for clients and managing for vets. */
const productsStore = useProductsStore()
const cart = useCartStore()
const userStore = useUserStore()
const authStore = useAuthStore()
const appointmentsStore = useAppointmentsStore()

const category = ref('')
const search = ref('')
const openCart = ref(false)
// When true, vets only see their own products; toggles global vs mine scope
const mineOnly = ref(false)
// When true, vets see only unlisted products; when false, see only active
const showUnlisted = ref(false)

// Cache vet info for displaying on product cards
const vetCache = ref<Record<string, { name: string; clinic_address?: string }>>({})

/**
 * Get vet info by ID from cache.
 */
function getVetInfo(vetId: string): { name: string; clinic_address?: string } | null {
  return vetCache.value[vetId] || null
}

/**
 * Load vet info for all products.
 */
async function loadVetInfo(): Promise<void> {
  try {
    const vets = await appointmentsStore.listVeterinarians()
    for (const v of vets) {
      vetCache.value[v.id] = { name: v.name, clinic_address: v.clinic_address }
    }
  } catch (_) {}
}

// Delete confirmation modal state
const showDeleteModal = ref(false)
const deleteTarget = ref<Product | null>(null)
const deleting = ref(false)

const showModal = ref(false)
const editing = ref<Product | null>(null)
const saving = ref(false)
const form = reactive<{ name: string; description: string; category: string; price: number; stock_quantity: number; measurement_value: number | null; measurement_unit: string }>({
  name: '',
  description: '',
  category: '',
  price: 0,
  stock_quantity: 0,
  measurement_value: null,
  measurement_unit: '',
})

// Categories that require measurement units
const needsMeasurement = (cat: string) => ['Food', 'Medicine'].includes(cat)

// Measurement unit options per category
const measurementUnits: Record<string, string[]> = {
  Food: ['kg', 'g', 'L', 'mL', 'pcs'],
  Medicine: ['mL', 'mg', 'g', 'tablets', 'capsules'],
}

// local qty per product id for adding to cart
const qtyByProductId = reactive<Record<string, number>>({})
function getQty(id: string): number {
  return qtyByProductId[id] ?? 1
}
function onQtyInput(id: string, evt: Event): void {
  const v = parseInt((evt.target as HTMLInputElement).value, 10)
  qtyByProductId[id] = isFinite(v) && v > 0 ? v : 1
}

/** Update active category filter and refresh from API. */
function setCategory(cat: string): void {
  category.value = cat
  const opts: { veterinarianId?: string; force: boolean; category?: string; search?: string } = { force: true }
  if (mineOnly.value && userStore.isVeterinarian && userStore.profile) {
    opts.veterinarianId = userStore.profile.id
  }
  if (cat) opts.category = cat
  const q = search.value.trim()
  if (q) opts.search = q
  productsStore.fetchProducts(opts)
}

/** Sync server-side search with a light debounce (handled by browser input). */
function onSearch(): void {
  const q = search.value.trim()
  const opts: { veterinarianId?: string; force: boolean; category?: string; search?: string } = { force: true }
  if (mineOnly.value && userStore.isVeterinarian && userStore.profile) {
    opts.veterinarianId = userStore.profile.id
  }
  if (category.value) opts.category = category.value
  if (q) opts.search = q
  productsStore.fetchProducts(opts)
}

/** Toggle scope: when on, only the logged-in vet's products are shown. */
function toggleMine(): void {
  mineOnly.value = !mineOnly.value
  const opts: { veterinarianId?: string; force: boolean; category?: string; search?: string } = { force: true }
  if (mineOnly.value && userStore.isVeterinarian && userStore.profile) {
    opts.veterinarianId = userStore.profile.id
  }
  if (category.value) opts.category = category.value
  const q = search.value.trim()
  if (q) opts.search = q
  productsStore.fetchProducts(opts)
}

/** Toggle unlisted filter: shows only unlisted products when on. */
function toggleUnlisted(): void {
  showUnlisted.value = !showUnlisted.value
}

// Make ordering deterministic so cards don't jump around. We'll sort by name ASC,
// and fall back to created_at DESC when name is missing or equal.
const filtered = computed(() => {
  // Vets see all products, clients only see active ones
  let base = userStore.isVeterinarian ? productsStore.products : productsStore.activeProducts
  
  // For vets: filter by unlisted toggle
  if (userStore.isVeterinarian) {
    base = showUnlisted.value
      ? base.filter(p => !p.is_active)  // Show only unlisted
      : base.filter(p => p.is_active)   // Show only active (default)
  }
  
  const scoped = (mineOnly.value && userStore.isVeterinarian && userStore.profile)
    ? base.filter(p => p.veterinarian_id === userStore.profile!.id)
    : base
  const list = scoped.filter(p => (category.value ? p.category === category.value : true) && (search.value ? (p.name?.toLowerCase().includes(search.value.toLowerCase()) || p.description?.toLowerCase().includes(search.value.toLowerCase())) : true))
  return [...list].sort((a, b) => {
    const an = (a.name || '').toLowerCase()
    const bn = (b.name || '').toLowerCase()
    if (an && bn && an !== bn) return an < bn ? -1 : 1
    // tiebreaker: newest first
    const at = new Date(a.created_at).getTime()
    const bt = new Date(b.created_at).getTime()
    return bt - at
  })
})

/** Add a product to cart with a quantity of 1. */
function addToCart(p: Product): void {
  const qty = getQty(p.id)
  cart.addItem({ productId: p.id, name: p.name, price: p.price, quantity: qty })
}

/** Begin creating a new product (vet only). */
function startCreate(): void {
  editing.value = null
  Object.assign(form, { name: '', description: '', category: '', price: 0, stock_quantity: 0, measurement_value: null, measurement_unit: '' })
  showModal.value = true
}

/** Begin editing an existing product (vet only). */
function startEdit(p: Product): void {
  editing.value = p
  Object.assign(form, {
    name: p.name,
    description: p.description || '',
    category: p.category || '',
    price: p.price,
    stock_quantity: p.stock_quantity,
    measurement_value: p.dimensions?.measurement_value ?? null,
    measurement_unit: p.dimensions?.measurement_unit || '',
  })
  showModal.value = true
}

/** Close the modal without saving. */
function cancelEdit(): void {
  showModal.value = false
}

/** Persist new or edited product to the backend. */
async function saveProduct(): Promise<void> {
  saving.value = true
  try {
    // Build payload with dimensions for measurement
    const payload: Record<string, unknown> = {
      name: form.name,
      description: form.description,
      category: form.category,
      price: form.price,
      stock_quantity: form.stock_quantity,
    }
    // Include measurement in dimensions if applicable
    if (needsMeasurement(form.category) && form.measurement_value && form.measurement_unit) {
      payload.dimensions = {
        measurement_value: form.measurement_value,
        measurement_unit: form.measurement_unit,
      }
    }
    if (editing.value) {
      await productsStore.updateProduct(editing.value.id, payload)
    } else {
      await productsStore.createProduct(payload)
    }
    // Force-refresh list so newly saved item shows up immediately
    const opts: { veterinarianId?: string; force: boolean; category?: string; search?: string } = { force: true }
    if (mineOnly.value && userStore.isVeterinarian && userStore.profile) {
      opts.veterinarianId = userStore.profile.id
    }
    if (category.value) opts.category = category.value
    const q = search.value.trim()
    if (q) opts.search = q
    await productsStore.fetchProducts(opts)
    showModal.value = false
  } finally {
    saving.value = false
  }
}

/** Flip active state for a product (soft delete/restore). */
async function toggleActive(p: Product): Promise<void> {
  if (p.is_active) {
    await productsStore.deactivateProduct(p.id)
  } else {
    await productsStore.updateProduct(p.id, { is_active: true })
  }
}

/** Open delete confirmation modal. */
function confirmDelete(p: Product): void {
  deleteTarget.value = p
  showDeleteModal.value = true
}

/** Cancel delete. */
function cancelDelete(): void {
  showDeleteModal.value = false
  deleteTarget.value = null
}

/** Permanently delete a product. */
async function deleteProduct(): Promise<void> {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await productsStore.deleteProduct(deleteTarget.value.id)
    showDeleteModal.value = false
    deleteTarget.value = null
  } finally {
    deleting.value = false
  }
}

/**
 * Format price with measurement unit like "₱50 (12pcs)".
 */
function formatPriceWithMeasurement(p: Product): string {
  const priceStr = formatPHP(p.price)
  if (p.dimensions?.measurement_value && p.dimensions?.measurement_unit) {
    return `${priceStr} (${p.dimensions.measurement_value}${p.dimensions.measurement_unit})`
  }
  return priceStr
}

// Load once and then rely on store TTL to avoid reshuffling on every tab switch
onMounted(async () => {
  if (userStore.isVeterinarian && userStore.profile) {
    mineOnly.value = true
    productsStore.fetchProducts({ veterinarianId: userStore.profile.id, force: true })
  } else {
    mineOnly.value = false
    productsStore.fetchProducts({ force: true })
  }
  // Load vet info for product cards
  await loadVetInfo()
})

watch(() => userStore.isVeterinarian, (isVet) => {
  if (isVet && userStore.profile) {
    mineOnly.value = true
    productsStore.fetchProducts({ veterinarianId: userStore.profile.id, force: true })
  } else {
    mineOnly.value = false
    productsStore.fetchProducts({ force: true })
  }
})

async function applyCheckout(items: Array<{ productId: string; quantity: number }>): Promise<void> {
  // Group items by veterinarian for order creation
  const vetToItems: Record<string, Array<{ product_id: string; quantity: number }>> = {}
  for (const it of items) {
    const prod = productsStore.products.find(p => p.id === it.productId)
    if (!prod) continue
    const vetId = prod.veterinarian_id
    if (!vetToItems[vetId]) vetToItems[vetId] = []
    vetToItems[vetId].push({ product_id: it.productId, quantity: it.quantity })
  }

  // Fire one order per veterinarian
  try {
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
    }
    if (authStore.session?.access_token) headers['Authorization'] = `Bearer ${authStore.session.access_token}`

    await Promise.all(
      Object.entries(vetToItems).map(async ([vetId, orderItems]) => {
        const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/orders`, {
          method: 'POST',
          headers,
          body: JSON.stringify({ veterinarian_id: vetId, items: orderItems }),
        })
        const body = await res.json().catch(() => ({}))
        if (!res.ok) throw new Error(body.error || res.statusText)
      })
    )
  } catch (_) {
    // Ignore for now; UI will still update below
  }

  // Update local stock optimistically
  for (const it of items) {
    const prod = productsStore.products.find(p => p.id === it.productId)
    if (prod) prod.stock_quantity = Math.max(0, prod.stock_quantity - it.quantity)
  }
}
</script>

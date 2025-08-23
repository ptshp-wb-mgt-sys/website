<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">{{ pageTitle }}</h1>
        <p class="text-gray-600 mt-1">{{ pageSubtitle }}</p>
      </div>
      <Button v-if="userStore.isClient || userStore.isAdmin" @click="openAddPetModal">
        <Plus class="w-4 h-4 mr-2" />
        Add New Pet
      </Button>
    </div>

    <!-- Search and Filters -->
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
          <input
            type="text"
            placeholder="Search pets..."
            class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-aquamarine focus:border-transparent"
            v-model="searchQuery"
          />
        </div>
        <select v-if="userStore.isVeterinarian || userStore.isAdmin" v-model="vetSpecies" class="px-3 py-2 border border-gray-300 rounded-lg w-48 md:w-56">
          <option value="">All Species</option>
          <option value="dog">Dogs</option>
          <option value="cat">Cats</option>
          <option value="bird">Birds</option>
          <option value="rabbit">Rabbits</option>
          <option value="hamster">Hamsters</option>
          <option value="fish">Fish</option>
          <option value="other">Other</option>
        </select>
        <select v-if="userStore.isVeterinarian || userStore.isAdmin" v-model="vetOwnerId" class="px-3 py-2 border border-gray-300 rounded-lg w-56 md:w-64">
          <option value="">All Owners</option>
          <option v-for="o in vetOwnerOptions" :key="o.id" :value="o.id">{{ o.name }}</option>
        </select>
      </div>
      <div class="text-sm text-gray-600">
        {{ totalPetsText }}
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="petsStore.loading" class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="w-8 h-8 animate-spin text-aquamarine mx-auto mb-4" />
        <p class="text-gray-600">Loading pets...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="petsStore.error" class="text-center py-12">
      <div class="text-red-600 mb-4">
        <AlertCircle class="w-12 h-12 mx-auto mb-2" />
        <p>{{ petsStore.error }}</p>
      </div>
      <Button @click="petsStore.fetchPets()" variant="outline">
        Try Again
      </Button>
    </div>

    <!-- Pet Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Client's Own Pets -->
      <template v-if="userStore.isClient">
        <Card 
          v-for="pet in filteredClientPets" 
          :key="pet.id" 
          class="p-6 cursor-pointer hover:shadow-md transition"
          @click="$router.push({ name: 'pet-profile', params: { id: pet.id } })"
        >
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
                {{ pet.type }}
              </span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ formatAge(pet.date_of_birth) }}</p>
              <p><span class="font-medium">Weight:</span> {{ pet.weight }} kg</p>
            </div>
            
            <div class="flex space-x-2" @click.stop>
              <Button variant="ghost" size="sm" @click="openQRModal(pet.id)" title="QR Code">
                <QrCode class="w-4 h-4" />
              </Button>
              <Button variant="outline" size="sm" class="flex-1" @click="$router.push({ name: 'pet-profile', params: { id: pet.id } })">
                <FileText class="w-4 h-4 mr-1" />
                Profile
              </Button>
              <Button variant="ghost" size="sm" class="flex-1" @click="openEditPetModal(pet)">Edit</Button>
              <Button 
                variant="ghost" 
                size="sm" 
                class="text-red-600 hover:text-red-700"
                @click="deletePet(pet.id)"
              >
                <Trash2 class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      </template>

      <!-- Veterinarian's Patient View -->
      <template v-if="userStore.isVeterinarian">
        <Card v-for="pet in filteredVetPets" :key="pet.id" class="p-6 cursor-pointer hover:shadow-md transition" @click="handleVetCardClick(pet.id)">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">{{ pet.type }}</span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ formatAge(pet.date_of_birth) }}</p>
              <p><span class="font-medium">Owner:</span> {{ ownerNames[pet.owner_id] || '—' }}</p>
            </div>
            
            <div class="flex space-x-2" @click.stop>
              <Button variant="outline" size="sm" class="flex-1" @click="$router.push({ name: 'pet-profile', params: { id: pet.id }, query: newRecordQuery })">
                <FileText class="w-4 h-4 mr-1" />
                Records
              </Button>
              <Button variant="ghost" size="sm" class="flex-1" @click="openQRModal(pet.id)">
                <QrCode class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      </template>

      <!-- Admin View -->
      <template v-if="userStore.isAdmin">
        <Card v-for="pet in filteredVetPets" :key="pet.id" class="p-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">{{ pet.type }}</span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ formatAge(pet.date_of_birth) }}</p>
              <p><span class="font-medium">Created:</span> {{ (pet as any).created_at?.slice ? (pet as any).created_at.slice(0,10) : '' }}</p>
            </div>
            
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" class="flex-1" @click="$router.push({ name: 'pet-profile', params: { id: pet.id }, query: newRecordQuery })">
                <FileText class="w-4 h-4 mr-1" />
                View
              </Button>
              <Button variant="ghost" size="sm" class="flex-1" @click="openEditPetModal(pet as any)">Edit</Button>
              <Button variant="ghost" size="sm" class="text-red-600 hover:text-red-700" @click="deletePet(pet.id)">
                <Trash2 class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      </template>
    </div>

    <!-- Empty State -->
    <Card v-if="shouldShowEmptyState" class="p-12 text-center">
      <Heart class="w-16 h-16 text-gray-300 mx-auto mb-4" />
      <h3 class="text-xl font-semibold text-gray-600 mb-2">{{ emptyStateTitle }}</h3>
      <p class="text-gray-500 mb-6">{{ emptyStateMessage }}</p>
      <Button v-if="userStore.isClient || userStore.isAdmin" @click="openAddPetModal">
        <Plus class="w-4 h-4 mr-2" />
        Add Your First Pet
      </Button>
    </Card>

    <!-- Add Pet Modal -->
    <AddPetModal
      :is-open="showAddPetModal"
      @close="closeAddPetModal"
      @pet-added="handlePetAdded"
    />

    <!-- Edit Pet Modal -->
    <EditPetModal
      :is-open="showEditPetModal"
      :pet="selectedPet"
      @close="closeEditPetModal"
      @pet-updated="handlePetUpdated"
    />

    <!-- QR Preview Modal -->
    <QRPreviewModal :is-open="showQRModal" :pet-id="qrPetId" @close="closeQRModal" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePetsStore, type Pet } from '@/stores/pets'
import { Plus, Search, QrCode, FileText, Heart, Trash2, Loader2, AlertCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import AddPetModal from '@/components/AddPetModal.vue'
import EditPetModal from '@/components/EditPetModal.vue'
import QRPreviewModal from '@/components/QRPreviewModal.vue'
import { useQRCodesStore } from '@/stores/qrcodes'
import { useAppointmentsStore } from '@/stores/appointments'

const userStore = useUserStore()
const petsStore = usePetsStore()
const qrStore = useQRCodesStore()
const apptStore = useAppointmentsStore()
const route = useRoute()
const router = useRouter()

// Search/filter state
const searchQuery = ref('')
const vetSpecies = ref('')
const vetOwnerId = ref('')

// Owner name cache for vet cards
const ownerNames = ref<Record<string, string>>({})
const vetOwnerOptions = computed(() => {
  // Build owner options from currently known vet patients
  const map = new Map<string, string>()
  for (const p of vetPets.value) {
    if (p.owner_id) {
      const name = ownerNames.value[p.owner_id] || ''
      if (name) map.set(p.owner_id, name)
    }
  }
  return Array.from(map.entries()).map(([id, name]) => ({ id, name })).sort((a, b) => a.name.localeCompare(b.name))
})

// Modal state
const showAddPetModal = ref(false)
const showEditPetModal = ref(false)
const selectedPet = ref<Pet | null>(null)
const showQRModal = ref(false)
const qrPetId = ref<string>('')

/**
 * Format age from date of birth
 */
const formatAge = (dateOfBirth: string) => {
  // Handle both ISO date strings and simple date strings
  const birthDate = new Date(dateOfBirth + 'T00:00:00')
  const today = new Date()
  const ageInYears = today.getFullYear() - birthDate.getFullYear()
  const monthDiff = today.getMonth() - birthDate.getMonth()
  
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    return `${ageInYears - 1} years`
  }
  
  if (ageInYears === 0) {
    const ageInMonths = monthDiff + (today.getDate() >= birthDate.getDate() ? 0 : -1)
    return ageInMonths <= 0 ? 'Less than 1 month' : `${ageInMonths} months`
  }
  
  return `${ageInYears} years`
}

/**
 * Resolve and cache owner names for a list of user IDs.
 */
const warmOwnerNames = async (ownerIds: string[]) => {
  const unique = Array.from(new Set(ownerIds)).filter(id => !!id && !ownerNames.value[id])
  if (unique.length === 0) return
  await Promise.all(unique.map(async (id) => {
    try {
      const { useAuthStore } = await import('@/stores/auth')
      const auth = useAuthStore()
      const res = await fetch(`http://localhost:3000/api/v1/owners/${encodeURIComponent(id)}/label`, {
        headers: {
          'Authorization': `Bearer ${auth.session?.access_token || ''}`,
          'Content-Type': 'application/json',
        },
      })
      if (!res.ok) return
      const body = await res.json()
      const user = body.data || body
      if (user && user.name) {
        ownerNames.value[id] = user.name
      }
    } catch (_) {}
  }))
}

/**
 * Open add pet modal
 */
const openAddPetModal = () => {
  showAddPetModal.value = true
}

/**
 * Close add pet modal
 */
const closeAddPetModal = () => {
  showAddPetModal.value = false
}

/**
 * Handle pet added successfully
 */
const handlePetAdded = () => {
  // Pet is already added to store by the modal
  // Just close the modal
  closeAddPetModal()
}

/**
 * Open edit pet modal
 */
const openEditPetModal = (pet: Pet) => {
  selectedPet.value = pet
  showEditPetModal.value = true
}
/**
 * Open QR preview modal (fetch or generate QR first)
 */
const openQRModal = async (petId: string) => {
  try {
    await qrStore.getOrCreateForPet(petId)
    qrPetId.value = petId
    showQRModal.value = true
  } catch (e) {
    console.error('Failed to load QR', e)
  }
}

/** Close QR modal */
const closeQRModal = () => {
  showQRModal.value = false
  qrPetId.value = ''
}

/**
 * Close edit pet modal
 */
const closeEditPetModal = () => {
  showEditPetModal.value = false
  selectedPet.value = null
}

/**
 * Handle pet updated successfully
 */
const handlePetUpdated = () => {
  // Pet is already updated in store by the modal
  // Just close the modal
  closeEditPetModal()
}

/**
 * Delete a pet
 */
const deletePet = async (petId: string) => {
  if (confirm('Are you sure you want to delete this pet? This action cannot be undone.')) {
    try {
      await petsStore.deletePet(petId)
      // Pet is already removed from store by the store
    } catch (error) {
      console.error('Error deleting pet:', error)
      // Error is handled by the store
    }
  }
}

/**
 * Initialize component
 */
onMounted(async () => {
  if (userStore.isClient) {
    // Check if user has a profile first
    if (!userStore.hasProfile) {
      await userStore.fetchProfile()
      if (!userStore.hasProfile) {
        // Redirect to profile setup if user doesn't have a profile
        window.location.href = '/profile-setup'
        return
      }
    }
    await petsStore.fetchPets()
  }
  if (userStore.isVeterinarian || userStore.isAdmin) {
    // Prefer store-cached vet patients; load if missing
    if (petsStore.vetPatients.length === 0) {
      await petsStore.loadVetPatients()
    }
    vetPets.value = petsStore.vetPatients
    // Warm owner names for vet cards
    try {
      await warmOwnerNames(vetPets.value.map(p => p.owner_id))
    } catch (_) {}
  }
})

const vetPets = ref<Pet[]>([])

/**
 * Build a vet-facing patients list by hydrating unique pet ids
 * from the vet's appointment history.
 */
const loadVetPatients = async () => {
  await petsStore.loadVetPatients()
  vetPets.value = petsStore.vetPatients
  try {
    await warmOwnerNames(vetPets.value.map(p => p.owner_id))
  } catch (_) {}
}

/**
 * Computed properties for dynamic content
 */
const pageTitle = computed(() => {
  if (userStore.isClient) return 'My Pets'
  if (userStore.isVeterinarian) return 'Patients'
  if (userStore.isAdmin) return 'All Pets'
  return 'Pets'
})

const pageSubtitle = computed(() => {
  if (userStore.isClient) return 'Manage your beloved pets and their information'
  if (userStore.isVeterinarian) return 'View patient information and medical history'
  if (userStore.isAdmin) return 'System-wide pet management and administration'
  return ''
})

const totalPetsText = computed(() => {
  const count = userStore.isClient ? filteredClientPets.value.length : filteredVetPets.value.length
  const petWord = count === 1 ? 'pet' : 'pets'
  
  if (userStore.isClient) return `${count} ${petWord}`
  if (userStore.isVeterinarian) return `${count} patients`
  if (userStore.isAdmin) return `${count} total ${petWord}`
  return `${count} ${petWord}`
})

const shouldShowEmptyState = computed(() => {
  if (userStore.isClient) return filteredClientPets.value.length === 0
  return filteredVetPets.value.length === 0
})

/**
 * Determine if we are in the quick add-record flow (from dashboard).
 */
const isQuickNewRecord = computed(() => route.query.action === 'new-record')

/**
 * Provide query to propagate the new-record intent when navigating deeper.
 */
const newRecordQuery = computed(() => (isQuickNewRecord.value ? { action: 'new-record' } : {}))

/**
 * Navigate to a pet profile; preserve new-record intent if present.
 */
const handleVetCardClick = (petId: string) => {
  const q = isQuickNewRecord.value ? { action: 'new-record' } : {}
  router.push({ name: 'pet-profile', params: { id: petId }, query: q })
}

const emptyStateTitle = computed(() => {
  if (userStore.isClient) return 'No pets yet'
  if (userStore.isVeterinarian) return 'No patients found'
  if (userStore.isAdmin) return 'No pets in system'
  return 'No pets found'
})

const emptyStateMessage = computed(() => {
  if (userStore.isClient) return 'Start by adding your first pet to track their health and appointments.'
  if (userStore.isVeterinarian) return 'When pet owners register their pets, they will appear here.'
  if (userStore.isAdmin) return 'Pets will appear here as users register them in the system.'
  return 'Get started by adding a pet.'
})

/**
 * Filter client's pets by search text (name, type, breed)
 */
const filteredClientPets = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return petsStore.pets
  return petsStore.pets.filter(p =>
    p.name.toLowerCase().includes(q) ||
    p.type.toLowerCase().includes(q) ||
    p.breed.toLowerCase().includes(q)
  )
})

/**
 * Filter vet/admin pets by search text + species dropdown
 */
const filteredVetPets = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  const species = vetSpecies.value
  const ownerId = vetOwnerId.value
  return vetPets.value.filter(p => {
    const matchesSpecies = !species || p.type.toLowerCase() === species
    const matchesOwner = !ownerId || p.owner_id === ownerId
    if (!q) return matchesSpecies && matchesOwner
    const matchesText = p.name.toLowerCase().includes(q) || p.type.toLowerCase().includes(q) || p.breed.toLowerCase().includes(q)
    return matchesSpecies && matchesOwner && matchesText
  })
})
</script> 

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
          />
        </div>
        <select v-if="userStore.isVeterinarian || userStore.isAdmin" class="px-3 py-2 border border-gray-300 rounded-lg">
          <option value="">All Species</option>
          <option value="dog">Dogs</option>
          <option value="cat">Cats</option>
          <option value="bird">Birds</option>
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
          v-for="pet in petsStore.pets" 
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
              <Button variant="outline" size="sm" class="flex-1" @click="openQRModal(pet.id)">
                <QrCode class="w-4 h-4 mr-1" />
                QR Code
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
        <Card v-for="pet in allPets" :key="pet.id" class="p-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
                {{ pet.species }}
              </span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Owner:</span> {{ pet.ownerName }}</p>
              <p><span class="font-medium">Phone:</span> {{ pet.ownerPhone }}</p>
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ pet.age }}</p>
              <p><span class="font-medium">Last Visit:</span> {{ pet.lastVisit || 'Never' }}</p>
            </div>
            
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" class="flex-1">
                <FileText class="w-4 h-4 mr-1" />
                Records
              </Button>
              <Button variant="ghost" size="sm" class="flex-1">
                <QrCode class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </Card>
      </template>

      <!-- Admin View -->
      <template v-if="userStore.isAdmin">
        <Card v-for="pet in allPets" :key="pet.id" class="p-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
                {{ pet.species }}
              </span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Owner:</span> {{ pet.ownerName }}</p>
              <p><span class="font-medium">Email:</span> {{ pet.ownerEmail }}</p>
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ pet.age }}</p>
              <p><span class="font-medium">Created:</span> {{ pet.createdDate }}</p>
            </div>
            
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" class="flex-1">
                <FileText class="w-4 h-4 mr-1" />
                View
              </Button>
              <Button variant="ghost" size="sm" class="flex-1">Edit</Button>
              <Button variant="ghost" size="sm" class="text-red-600 hover:text-red-700">
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
import { useUserStore } from '@/stores/user'
import { usePetsStore, type Pet } from '@/stores/pets'
import { Plus, Search, QrCode, FileText, Heart, Trash2, Loader2, AlertCircle } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import AddPetModal from '@/components/AddPetModal.vue'
import EditPetModal from '@/components/EditPetModal.vue'
import QRPreviewModal from '@/components/QRPreviewModal.vue'
import { useQRCodesStore } from '@/stores/qrcodes'

const userStore = useUserStore()
const petsStore = usePetsStore()
const qrStore = useQRCodesStore()

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
})

const allPets = [
  {
    id: '1',
    name: 'Buddy',
    species: 'Dog',
    breed: 'Golden Retriever',
    age: '3 years',
    weight: '30 kg',
    ownerName: 'John Smith',
    ownerPhone: '+1 (555) 123-4567',
    ownerEmail: 'john@example.com',
    lastVisit: '2 days ago',
    createdDate: '2023-06-15'
  },
  {
    id: '2',
    name: 'Luna',
    species: 'Cat',
    breed: 'Persian Cat',
    age: '2 years',
    weight: '4 kg',
    ownerName: 'John Smith',
    ownerPhone: '+1 (555) 123-4567',
    ownerEmail: 'john@example.com',
    lastVisit: '2 days ago',
    createdDate: '2023-06-15'
  },
  {
    id: '3',
    name: 'Max',
    species: 'Dog',
    breed: 'German Shepherd',
    age: '5 years',
    weight: '35 kg',
    ownerName: 'Sarah Johnson',
    ownerPhone: '+1 (555) 987-6543',
    ownerEmail: 'sarah@example.com',
    lastVisit: '1 week ago',
    createdDate: '2023-05-20'
  },
  {
    id: '4',
    name: 'Whiskers',
    species: 'Cat',
    breed: 'Maine Coon',
    age: '4 years',
    weight: '6 kg',
    ownerName: 'Mike Davis',
    ownerPhone: '+1 (555) 456-7890',
    ownerEmail: 'mike@example.com',
    lastVisit: '3 days ago',
    createdDate: '2023-07-01'
  }
]

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
  const count = userStore.isClient ? petsStore.petsCount : allPets.length
  const petWord = count === 1 ? 'pet' : 'pets'
  
  if (userStore.isClient) return `${count} ${petWord}`
  if (userStore.isVeterinarian) return `${count} patients`
  if (userStore.isAdmin) return `${count} total ${petWord}`
  return `${count} ${petWord}`
})

const shouldShowEmptyState = computed(() => {
  if (userStore.isClient) return !petsStore.hasPets
  return allPets.length === 0
})

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
</script> 

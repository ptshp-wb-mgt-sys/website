<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">{{ pageTitle }}</h1>
        <p class="text-gray-600 mt-1">{{ pageSubtitle }}</p>
      </div>
      <Button v-if="userStore.isClient || userStore.isAdmin">
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

    <!-- Pet Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Client's Own Pets -->
      <template v-if="userStore.isClient">
        <Card v-for="pet in clientPets" :key="pet.id" class="p-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-rich-black">{{ pet.name }}</h3>
              <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
                {{ pet.species }}
              </span>
            </div>
            
            <div class="space-y-2 text-sm">
              <p><span class="font-medium">Breed:</span> {{ pet.breed }}</p>
              <p><span class="font-medium">Age:</span> {{ pet.age }}</p>
              <p><span class="font-medium">Weight:</span> {{ pet.weight }}</p>
            </div>
            
            <div class="flex space-x-2">
              <Button variant="outline" size="sm" class="flex-1">
                <QrCode class="w-4 h-4 mr-1" />
                QR Code
              </Button>
              <Button variant="ghost" size="sm" class="flex-1">Edit</Button>
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
      <Button v-if="userStore.isClient || userStore.isAdmin">
        <Plus class="w-4 h-4 mr-2" />
        Add Your First Pet
      </Button>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { Plus, Search, QrCode, FileText, Heart, Trash2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const userStore = useUserStore()

/**
 * Mock data - in real app this would come from API
 */
const clientPets = [
  {
    id: '1',
    name: 'Buddy',
    species: 'Dog',
    breed: 'Golden Retriever',
    age: '3 years',
    weight: '30 kg'
  },
  {
    id: '2',
    name: 'Luna',
    species: 'Cat',
    breed: 'Persian Cat',
    age: '2 years',
    weight: '4 kg'
  }
]

const allPets = [
  ...clientPets.map(pet => ({
    ...pet,
    ownerName: 'John Smith',
    ownerPhone: '+1 (555) 123-4567',
    ownerEmail: 'john@example.com',
    lastVisit: '2 days ago',
    createdDate: '2023-06-15'
  })),
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
  const count = userStore.isClient ? clientPets.length : allPets.length
  const petWord = count === 1 ? 'pet' : 'pets'
  
  if (userStore.isClient) return `${count} ${petWord}`
  if (userStore.isVeterinarian) return `${count} patients`
  if (userStore.isAdmin) return `${count} total ${petWord}`
  return `${count} ${petWord}`
})

const shouldShowEmptyState = computed(() => {
  if (userStore.isClient) return clientPets.length === 0
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

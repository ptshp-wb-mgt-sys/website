import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

// Types for pets
export interface Pet {
  id: string
  name: string
  type: string
  breed: string
  date_of_birth: string
  weight: number
  owner_id: string
  created_at: string
  updated_at: string
}

export interface CreatePetRequest {
  name: string
  type: string
  breed: string
  date_of_birth: string
  weight: number
  owner_id?: string
}

export interface UpdatePetRequest {
  name?: string
  type?: string
  breed?: string
  date_of_birth?: string
  weight?: number
}

export const usePetsStore = defineStore('pets', () => {
  const pets = ref<Pet[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAt = ref<number | null>(null)

  const authStore = useAuthStore()

  // Computed properties
  const petsCount = computed(() => pets.value.length)
  const hasPets = computed(() => pets.value.length > 0)

  /**
   * Fetch pets for the current user.
   * Skips network if data is fresh unless `force` is true.
   */
  const fetchPets = async (options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 2 * 60 * 1000
    if (!force && pets.value.length > 0 && lastFetchedAt.value && Date.now() - lastFetchedAt.value < ttlMs) {
      return
    }
    if (!authStore.session?.access_token) {
      error.value = 'No authentication token'
      return
    }

    loading.value = true
    error.value = null

    try {
      // Get user ID from auth store
      const userId = authStore.user?.id
      if (!userId) {
        throw new Error('User ID not found')
      }

      const response = await fetch(`http://localhost:3000/api/v1/clients/${userId}/pets`, {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        if (response.status === 404) {
          // User doesn't have a profile yet, return empty array
          pets.value = []
          lastFetchedAt.value = Date.now()
          return
        }
        throw new Error(`Failed to fetch pets: ${response.statusText}`)
      }

      const data = await response.json()
      pets.value = data.data || data
      lastFetchedAt.value = Date.now()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch pets'
      console.error('Error fetching pets:', err)
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new pet
   */
  const createPet = async (petData: CreatePetRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch('http://localhost:3000/api/v1/pets', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(petData)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to create pet')
      }

      const data = await response.json()
      const newPet = data.data || data
      
      // Add to local state
      pets.value.push(newPet)
      return newPet
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create pet'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update a pet
   */
  const updatePet = async (petId: string, updates: UpdatePetRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/pets/${petId}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updates)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to update pet')
      }

      const data = await response.json()
      const updatedPet = data.data || data
      
      // Update in local state
      const index = pets.value.findIndex(pet => pet.id === petId)
      if (index !== -1) {
        pets.value[index] = updatedPet
      }
      
      return updatedPet
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update pet'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Delete a pet
   */
  const deletePet = async (petId: string) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/pets/${petId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to delete pet')
      }

      // Remove from local state
      pets.value = pets.value.filter(pet => pet.id !== petId)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete pet'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Get a pet by ID
   */
  const getPet = async (petId: string) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    try {
      const response = await fetch(`http://localhost:3000/api/v1/pets/${petId}`, {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch pet: ${response.statusText}`)
      }

      const data = await response.json()
      return data.data || data
    } catch (err) {
      console.error('Error fetching pet:', err)
      throw err
    }
  }

  /**
   * Clear pets data
   */
  const clearPets = () => {
    pets.value = []
    error.value = null
    lastFetchedAt.value = null
  }

  /**
   * Initialize pets after authentication
   */
  const initialize = async () => {
    if (authStore.isAuthenticated) {
      await fetchPets()
    }
  }

  return {
    // State
    pets,
    loading,
    error,
    lastFetchedAt,
    
    // Computed
    petsCount,
    hasPets,
    
    // Actions
    fetchPets,
    createPet,
    updatePet,
    deletePet,
    getPet,
    clearPets,
    initialize
  }
}) 

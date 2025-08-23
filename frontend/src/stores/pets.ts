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
  // Simple in-memory caches to speed up vet flows
  const petByIdCache = ref<Record<string, { pet: Pet; fetchedAt: number }>>({})
  const vetPatients = ref<Pet[]>([])
  const vetPatientsFetchedAt = ref<number | null>(null)
  // Owner label cache to support vet views (owner names on patient cards)
  const ownerNameById = ref<Record<string, { name: string; fetchedAt: number }>>({})

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

      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/clients/${userId}/pets`, {
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
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets`, {
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
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets/${petId}`, {
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
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets/${petId}`, {
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
   * Get a pet by ID with a tiny TTL cache.
   * Returns instantly from memory when available.
   */
  const getPet = async (petId: string, options?: { ttlMs?: number }) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    try {
      // Try the client-owned list first
      const fromList = pets.value.find(p => p.id === petId)
      if (fromList) return fromList

      // Serve from cache when fresh
      const ttlMs = options?.ttlMs ?? 5 * 60 * 1000
      const cached = petByIdCache.value[petId]
      if (cached && Date.now() - cached.fetchedAt < ttlMs) {
        return cached.pet
      }

      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets/${petId}`, {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch pet: ${response.statusText}`)
      }

      const data = await response.json()
      const pet = (data.data || data) as Pet
      petByIdCache.value[petId] = { pet, fetchedAt: Date.now() }
      return pet
    } catch (err) {
      console.error('Error fetching pet:', err)
      throw err
    }
  }

  /**
   * Build and cache a vet's patients list from appointments.
   * Hydrates in parallel and caches for a short TTL.
   */
  const loadVetPatients = async (options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 5 * 60 * 1000
    if (!force && vetPatients.value.length > 0 && vetPatientsFetchedAt.value && Date.now() - vetPatientsFetchedAt.value < ttlMs) {
      return
    }
    const { useAppointmentsStore } = await import('./appointments')
    const apptStore = useAppointmentsStore()
    if (apptStore.appointments.length === 0) {
      await apptStore.fetchAppointments()
    }
    const uniqueIds = Array.from(new Set(apptStore.appointments.map(a => a.pet_id)))
    const loaded: Pet[] = []
    await Promise.all(uniqueIds.map(async (id) => {
      try {
        const p = await getPet(id)
        if (p) loaded.push(p)
      } catch (_) {}
    }))
    vetPatients.value = loaded
    vetPatientsFetchedAt.value = Date.now()
  }

  /**
   * Synchronous label accessor for rendering without flicker.
   * Falls back to empty string if unknown to avoid placeholders.
   */
  const getPetLabelSync = (petId: string): string => {
    const fromList = pets.value.find(p => p.id === petId)
    if (fromList) return `${fromList.name} (${fromList.type})`
    const cached = petByIdCache.value[petId]?.pet
    if (cached) return `${cached.name} (${cached.type})`
    const cachedVet = vetPatients.value.find(p => p.id === petId)
    if (cachedVet) return `${cachedVet.name} (${cachedVet.type})`
    return ''
  }

  /**
   * Warm up labels for a set of pet IDs in parallel, using by-id TTL cache.
   */
  const warmPetLabels = async (ids: string[]) => {
    const unique = Array.from(new Set(ids)).filter(id => !petByIdCache.value[id])
    await Promise.all(unique.map(async (id) => {
      try {
        const p = await getPet(id)
        petByIdCache.value[id] = { pet: p, fetchedAt: Date.now() }
      } catch (_) {}
    }))
  }

  /**
   * Synchronously get an owner display name from cache.
   */
  const getOwnerNameSync = (ownerId: string): string => {
    return ownerNameById.value[ownerId]?.name || ''
  }

  /**
   * Warm owner names for the given ids using a small TTL cache.
   */
  const warmOwnerNames = async (ids: string[], options?: { ttlMs?: number }) => {
    const ttlMs = options?.ttlMs ?? 5 * 60 * 1000
    const pending = Array.from(new Set(ids)).filter((id) => {
      if (!id) return false
      const cached = ownerNameById.value[id]
      return !cached || (Date.now() - cached.fetchedAt > ttlMs)
    })
    if (pending.length === 0) return
    for (const id of pending) {
      try {
        const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/owners/${encodeURIComponent(id)}/label`, {
          headers: {
            'Authorization': authStore.session?.access_token ? `Bearer ${authStore.session.access_token}` : '',
            'Content-Type': 'application/json',
          },
        })
        if (!res.ok) continue
        const body = await res.json().catch(() => ({}))
        const data = (body.data || body) as { name?: string }
        if (data?.name) {
          ownerNameById.value[id] = { name: data.name, fetchedAt: Date.now() }
        }
      } catch (_) {
        // ignore
      }
    }
  }

  /**
   * Clear pets data
   */
  const clearPets = () => {
    pets.value = []
    error.value = null
    lastFetchedAt.value = null
    petByIdCache.value = {}
    vetPatients.value = []
    vetPatientsFetchedAt.value = null
    ownerNameById.value = {}
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
    vetPatients,
    vetPatientsFetchedAt,
    ownerNameById,
    
    // Computed
    petsCount,
    hasPets,
    
    // Actions
    fetchPets,
    createPet,
    updatePet,
    deletePet,
    getPet,
    getPetLabelSync,
    warmPetLabels,
    getOwnerNameSync,
    warmOwnerNames,
    loadVetPatients,
    clearPets,
    initialize
  }
}) 

import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'

export interface MedicalRecord {
  id: string
  pet_id: string
  veterinarian_id: string
  date_of_visit: string
  reason_for_visit: string
  diagnosis: string
  medication_prescribed: string[]
  notes: string
  created_at: string
  updated_at: string
}

export interface CreateMedicalRecordRequest {
  date_of_visit?: string
  reason_for_visit: string
  diagnosis: string
  medication_prescribed: string[]
  notes?: string
}

export interface UpdateMedicalRecordRequest {
  date_of_visit?: string
  reason_for_visit?: string
  diagnosis?: string
  medication_prescribed?: string[]
  notes?: string
}

export const useMedicalRecordsStore = defineStore('medicalRecords', () => {
  const recordsByPetId = ref<Record<string, MedicalRecord[]>>({})
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAtByPetId = ref<Record<string, number>>({})

  const authStore = useAuthStore()

  /**
   * Load medical records for a given pet from the API and cache them by `petId`.
   * Uses a per-pet TTL cache unless `force` is set.
   */
  const fetchByPetId = async (petId: string, options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 2 * 60 * 1000
    const last = lastFetchedAtByPetId.value[petId]
    if (!force && recordsByPetId.value[petId] && last && Date.now() - last < ttlMs) {
      return recordsByPetId.value[petId]!
    }
    if (!authStore.session?.access_token) {
      error.value = 'No authentication token'
      return [] as MedicalRecord[]
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/pets/${encodeURIComponent(petId)}/medical-records`, {
        headers: {
          Authorization: `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch medical records: ${response.statusText}`)
      }

      const data = await response.json()
      const records = (data.data || data) as MedicalRecord[]
      recordsByPetId.value[petId] = records
      lastFetchedAtByPetId.value[petId] = Date.now()
      return records
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch medical records'
      console.error('Error fetching medical records:', err)
      return [] as MedicalRecord[]
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new medical record for a pet.
   */
  const createRecord = async (petId: string, payload: CreateMedicalRecordRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/pets/${encodeURIComponent(petId)}/medical-records`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to create medical record')
      }

      const data = await response.json()
      const newRecord = (data.data || data) as MedicalRecord
      const current = recordsByPetId.value[petId] || []
      recordsByPetId.value[petId] = [newRecord, ...current]
      return newRecord
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create medical record'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update an existing medical record by its ID.
   */
  const updateRecord = async (recordId: string, updates: UpdateMedicalRecordRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/medical-records/${encodeURIComponent(recordId)}`, {
        method: 'PUT',
        headers: {
          Authorization: `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updates),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to update medical record')
      }

      const data = await response.json()
      const updated = (data.data || data) as MedicalRecord

      // Update local cache if we can find the pet bucket
      const petId = updated.pet_id
      const bucket = recordsByPetId.value[petId]
      if (bucket) {
        recordsByPetId.value[petId] = bucket.map(r => (r.id === updated.id ? updated : r))
      }
      return updated
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update medical record'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Delete a medical record by ID and remove it from cache.
   */
  const deleteRecord = async (recordId: string, petId: string) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/medical-records/${encodeURIComponent(recordId)}`, {
        method: 'DELETE',
        headers: {
          Authorization: `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to delete medical record')
      }

      const current = recordsByPetId.value[petId] || []
      recordsByPetId.value[petId] = current.filter(r => r.id !== recordId)
      lastFetchedAtByPetId.value[petId] = Date.now()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to delete medical record'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Return cached records for a pet, or an empty array if not loaded.
   */
  const getCachedForPet = (petId: string) => {
    return recordsByPetId.value[petId] || []
  }

  return {
    // State
    recordsByPetId,
    loading,
    error,
    lastFetchedAtByPetId,
    // Actions
    fetchByPetId,
    createRecord,
    updateRecord,
    deleteRecord,
    getCachedForPet,
  }
})



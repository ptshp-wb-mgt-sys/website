import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

// Types for appointments
export interface Appointment {
  id: string
  client_id: string
  veterinarian_id: string
  pet_id: string
  appointment_date: string
  duration_minutes: number
  reason: string
  status: string
  notes?: string
  created_at: string
  updated_at: string
}

export interface CreateAppointmentRequest {
  veterinarian_id: string
  pet_id: string
  appointment_date: string
  duration_minutes?: number
  reason: string
  notes?: string
}

export interface UpdateAppointmentRequest {
  appointment_date?: string
  duration_minutes?: number
  reason?: string
  status?: 'scheduled' | 'completed' | 'cancelled' | 'rescheduled'
  notes?: string
}

export interface TimeSlot {
  start_time: string
  end_time: string
  available: boolean
}

export interface VeterinarianListItem {
  id: string
  name: string
  email: string
  phone?: string
  clinic_address?: string
}

export interface SetAvailabilityRequest {
  available_hours: Array<{
    day_of_week: string
    start: string
    end: string
  }>
  clinic_address?: string
}

export const useAppointmentsStore = defineStore('appointments', () => {
  const appointments = ref<Appointment[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAt = ref<number | null>(null)
  const lastFetchedUserId = ref<string | null>(null)

  const authStore = useAuthStore()

  /**
   * Load current user's appointments from the API.
   * Skips network if cached within ttl, unless `force`.
   */
  const fetchAppointments = async (options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 60 * 1000
    const currentUserId = authStore.user?.id || null
    const userChanged = currentUserId !== lastFetchedUserId.value
    if (!force && !userChanged && appointments.value.length > 0 && lastFetchedAt.value && Date.now() - lastFetchedAt.value < ttlMs) {
      return
    }
    if (!authStore.session?.access_token) {
      error.value = 'No authentication token'
      return
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch('http://localhost:3000/api/v1/appointments', {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
      })

      if (!response.ok) {
        throw new Error(`Failed to fetch appointments: ${response.statusText}`)
      }

      const data = await response.json()
      appointments.value = (data.data || data) as Appointment[]
      lastFetchedAt.value = Date.now()
      lastFetchedUserId.value = currentUserId
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch appointments'
      console.error('Error fetching appointments:', err)
    } finally {
      loading.value = false
    }
  }

  /**
   * Create a new appointment for the current user.
   */
  const createAppointment = async (request: CreateAppointmentRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch('http://localhost:3000/api/v1/appointments', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to create appointment')
      }

      const data = await response.json()
      const newAppointment = (data.data || data) as Appointment
      appointments.value.push(newAppointment)
      return newAppointment
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create appointment'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update an existing appointment.
   */
  const updateAppointment = async (id: string, updates: UpdateAppointmentRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`http://localhost:3000/api/v1/appointments/${encodeURIComponent(id)}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(updates),
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to update appointment')
      }

      const data = await response.json()
      const updated = (data.data || data) as Appointment
      const idx = appointments.value.findIndex(a => a.id === id)
      if (idx !== -1) {
        appointments.value[idx] = updated
      }
      return updated
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update appointment'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Cancel appointment (soft-cancel on backend).
   */
  const cancelAppointment = async (id: string) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    try {
      const response = await fetch(`http://localhost:3000/api/v1/appointments/${encodeURIComponent(id)}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
        },
      })
      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to cancel appointment')
      }

      // Optimistic update: set status to cancelled
      const idx = appointments.value.findIndex(a => a.id === id)
      if (idx !== -1) {
        appointments.value[idx] = { ...appointments.value[idx], status: 'cancelled', updated_at: new Date().toISOString() }
      }
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to cancel appointment'
      throw err
    }
  }

  /**
   * Get available time slots for a veterinarian on a given date.
   */
  const getAvailableSlots = async (veterinarianId: string, dateISO: string) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    // Backend expects YYYY-MM-DD
    const dateOnly = new Date(dateISO)
    const y = dateOnly.getFullYear()
    const m = String(dateOnly.getMonth() + 1).padStart(2, '0')
    const d = String(dateOnly.getDate()).padStart(2, '0')
    const dateParam = `${y}-${m}-${d}`

    const response = await fetch(
      `http://localhost:3000/api/v1/veterinarians/${encodeURIComponent(veterinarianId)}/availability?date=${encodeURIComponent(dateParam)}`,
      {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
      },
    )
    if (!response.ok) {
      throw new Error(`Failed to get availability: ${response.statusText}`)
    }
    const data = await response.json()
    return (data.data || data) as TimeSlot[]
  }

  /**
   * List veterinarians for booking UI.
   */
  const listVeterinarians = async (limit = 10, offset = 0) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    const response = await fetch(
      `http://localhost:3000/api/v1/veterinarians?limit=${limit}&offset=${offset}`,
      {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json',
        },
      },
    )
    if (!response.ok) {
      throw new Error(`Failed to list veterinarians: ${response.statusText}`)
    }
    const data = await response.json()
    return (data.data || data) as VeterinarianListItem[]
  }

  /**
   * Veterinarian: set availability windows
   */
  const setAvailability = async (veterinarianId: string, req: SetAvailabilityRequest) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    const response = await fetch(`http://localhost:3000/api/v1/veterinarians/${encodeURIComponent(veterinarianId)}/availability`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.session.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(req),
    })
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: response.statusText }))
      throw new Error(errorData.message || 'Failed to set availability')
    }
    return await response.json()
  }

  // Helpers
  const upcomingAppointments = computed(() => {
    const now = new Date()
    // Start of today in local timezone
    const startOfToday = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
    return [...appointments.value]
      .filter(a => {
        const t = new Date(a.appointment_date).getTime()
        return t >= startOfToday && a.status !== 'cancelled'
      })
      .sort((a, b) => new Date(a.appointment_date).getTime() - new Date(b.appointment_date).getTime())
  })

  const pastAppointments = computed(() => {
    const now = Date.now()
    return [...appointments.value]
      .filter(a => new Date(a.appointment_date).getTime() < now)
      .sort((a, b) => new Date(b.appointment_date).getTime() - new Date(a.appointment_date).getTime())
  })

  const todaysAppointments = computed(() => {
    const today = new Date()
    const y = today.getFullYear()
    const m = today.getMonth()
    const d = today.getDate()
    return appointments.value.filter(a => {
      const t = new Date(a.appointment_date)
      return t.getFullYear() === y && t.getMonth() === m && t.getDate() === d
    })
  })

  /**
   * Clear in-memory state.
   */
  const clearAppointments = () => {
    appointments.value = []
    error.value = null
    lastFetchedAt.value = null
  }

  /**
   * Initialize after auth.
   */
  const initialize = async () => {
    if (authStore.isAuthenticated) {
      await fetchAppointments()
    }
  }

  return {
    // State
    appointments,
    loading,
    error,
    lastFetchedAt,
    // Computed
    upcomingAppointments,
    pastAppointments,
    todaysAppointments,
    // Actions
    fetchAppointments,
    createAppointment,
    updateAppointment,
    cancelAppointment,
    getAvailableSlots,
    listVeterinarians,
    setAvailability,
    clearAppointments,
    initialize,
  }
})



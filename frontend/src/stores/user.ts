import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { useAuthStore } from './auth'

// Types for user profiles
export interface UserProfile {
  id: string
  name: string
  email: string
  phone?: string
  address?: string
  role: 'client' | 'veterinarian' | 'admin'
}

export interface ClientProfile extends UserProfile {
  role: 'client'
  address: string
}

export interface VeterinarianProfile extends UserProfile {
  role: 'veterinarian'
  clinic_address?: string
  available_hours?: Array<{
    day_of_week: string
    start: string
    end: string
  }>
}

export const useUserStore = defineStore('user', () => {
  const profile = ref<UserProfile | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)
  const lastFetchedAt = ref<number | null>(null)

  const authStore = useAuthStore()

  // Computed properties
  const isClient = computed(() => profile.value?.role === 'client')
  const isVeterinarian = computed(() => profile.value?.role === 'veterinarian')
  const isAdmin = computed(() => profile.value?.role === 'admin')
  const hasProfile = computed(() => !!profile.value)
  
  const displayName = computed(() => {
    if (!profile.value) return 'User'
    return profile.value.name || 'User'
  })

  const roleDisplayName = computed(() => {
    if (!profile.value) return ''
    
    switch (profile.value.role) {
      case 'client':
        return 'Pet Owner'
      case 'veterinarian':
        return 'Doctor'
      case 'admin':
        return 'Admin'
      default:
        return ''
    }
  })

  const fullDisplayName = computed(() => {
    if (!profile.value) return 'User'
    const role = roleDisplayName.value
    const name = displayName.value
    return role ? `${role} ${name}` : name
  })

  /**
   * Fetch user profile from backend with simple cache TTL.
   */
  const fetchProfile = async (options?: { force?: boolean; ttlMs?: number }) => {
    const force = options?.force === true
    const ttlMs = options?.ttlMs ?? 5 * 60 * 1000
    if (!force && profile.value && lastFetchedAt.value && Date.now() - lastFetchedAt.value < ttlMs) {
      return
    }
    if (!authStore.session?.access_token) {
      error.value = 'No authentication token'
      return
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/profile`, {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (!response.ok) {
        if (response.status === 404) {
          // User hasn't created profile yet
          profile.value = null
          lastFetchedAt.value = Date.now()
          return
        }
        throw new Error(`Failed to fetch profile: ${response.statusText}`)
      }

      const data = await response.json()
      
      // If we only get basic info, we need to fetch full profile
      if (data.user_id && !data.name) {
        await fetchFullProfile(data.user_id)
      } else {
        profile.value = data
      }
      lastFetchedAt.value = Date.now()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch profile'
      console.error('Error fetching profile:', err)
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch full profile by user ID
   */
  const fetchFullProfile = async (userId: string) => {
    if (!authStore.session?.access_token) return

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/users/${userId}`, {
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        }
      })

      if (response.ok) {
        const responseData = await response.json()
        
        // Extract the actual user data from the response
        const userData = responseData.data || responseData
        profile.value = userData
      }
    } catch (err) {
      console.error('Error fetching full profile:', err)
    }
  }

  /**
   * Create user profile
   */
  const createProfile = async (profileData: Omit<UserProfile, 'id'>) => {
    if (!authStore.session?.access_token) {
      throw new Error('No authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/users`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(profileData)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to create profile')
      }

      const data = await response.json()
      const created = data.data || data
      profile.value = created
      return created
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to create profile'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Update user profile
   */
  const updateProfile = async (updates: Partial<UserProfile>) => {
    if (!profile.value || !authStore.session?.access_token) {
      throw new Error('No profile or authentication token')
    }

    loading.value = true
    error.value = null

    try {
      const response = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/users/${profile.value.id}`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${authStore.session.access_token}`,
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(updates)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({ message: response.statusText }))
        throw new Error(errorData.message || 'Failed to update profile')
      }

      const data = await response.json()
      const updated = data.data || data
      profile.value = updated
      return updated
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to update profile'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Clear profile data
   */
  const clearProfile = () => {
    profile.value = null
    error.value = null
    lastFetchedAt.value = null
  }

  /**
   * Initialize profile after authentication
   */
  const initialize = async () => {
    if (authStore.isAuthenticated) {
      await fetchProfile()
    }
  }

  return {
    // State
    profile,
    loading,
    error,
    lastFetchedAt,
    
    // Computed
    isClient,
    isVeterinarian,
    isAdmin,
    hasProfile,
    displayName,
    roleDisplayName,
    fullDisplayName,
    
    // Actions
    fetchProfile,
    createProfile,
    updateProfile,
    clearProfile,
    initialize
  }
}) 

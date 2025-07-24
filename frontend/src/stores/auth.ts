import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { supabase } from '@/lib/supabase'
import type { User, Session } from '@supabase/supabase-js'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const session = ref<Session | null>(null)
  const loading = ref(true)

  const isAuthenticated = computed(() => !!user.value)

  /**
   * Initialize auth state and listen for auth changes
   */
  const initialize = async () => {
    try {
      const { data: { session: currentSession } } = await supabase.auth.getSession()
      session.value = currentSession
      user.value = currentSession?.user ?? null

      supabase.auth.onAuthStateChange(async (event, newSession) => {
        session.value = newSession
        user.value = newSession?.user ?? null
        
        // Initialize user profile when signed in
        if (event === 'SIGNED_IN' && newSession) {
          // Dynamically import to avoid circular dependency
          const { useUserStore } = await import('./user')
          const userStore = useUserStore()
          await userStore.initialize()
        }
        
        // Clear profile when signed out
        if (event === 'SIGNED_OUT') {
          const { useUserStore } = await import('./user')
          const userStore = useUserStore()
          userStore.clearProfile()
        }
      })
    } catch (error) {
      console.error('Error initializing auth:', error)
    } finally {
      loading.value = false
    }
  }

  /**
   * Sign up a new user
   */
  const signUp = async (email: string, password: string) => {
    const { data, error } = await supabase.auth.signUp({
      email,
      password,
    })
    return { data, error }
  }

  /**
   * Sign in an existing user
   */
  const signIn = async (email: string, password: string) => {
    const { data, error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })
    return { data, error }
  }

  /**
   * Sign out the current user
   */
  const signOut = async () => {
    const { error } = await supabase.auth.signOut()
    if (!error) {
      user.value = null
      session.value = null
    }
    return { error }
  }

  return {
    user,
    session,
    loading,
    isAuthenticated,
    initialize,
    signUp,
    signIn,
    signOut,
  }
}) 

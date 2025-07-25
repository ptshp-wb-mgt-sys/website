<template>
  <nav class="bg-white border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <RouterLink to="/" class="text-xl font-bold text-rich-black">
            PetCare
          </RouterLink>
        </div>
        
        <!-- Authenticated Navigation -->
        <div v-if="authStore.isAuthenticated" class="flex items-center space-x-8">
          <RouterLink 
            to="/dashboard" 
            class="text-gray-700 hover:text-aquamarine transition-colors"
            :class="{ 'text-aquamarine': $route.name === 'dashboard' }"
          >
            Dashboard
          </RouterLink>
          <RouterLink 
            to="/pets" 
            class="text-gray-700 hover:text-aquamarine transition-colors"
            :class="{ 'text-aquamarine': $route.name === 'pets' }"
          >
            Pets
          </RouterLink>
          <RouterLink 
            to="/appointments" 
            class="text-gray-700 hover:text-aquamarine transition-colors"
            :class="{ 'text-aquamarine': $route.name === 'appointments' }"
          >
            Appointments
          </RouterLink>
          <RouterLink 
            to="/products" 
            class="text-gray-700 hover:text-aquamarine transition-colors"
            :class="{ 'text-aquamarine': $route.name === 'products' }"
          >
            Products
          </RouterLink>
          
          <!-- User Menu -->
          <div class="flex items-center space-x-4">
            <RouterLink 
              to="/profile" 
              class="text-sm text-gray-600 hover:text-aquamarine transition-colors"
              :class="{ 'text-aquamarine': $route.name === 'profile' }"
            >
              <span v-if="userStore.hasProfile">{{ userStore.fullDisplayName }}</span>
              <span v-else-if="authStore.user?.email">{{ authStore.user.email }}</span>
              <span v-else>User</span>
            </RouterLink>
            <Button variant="ghost" size="sm" @click="handleSignOut">
              Sign out
            </Button>
          </div>
        </div>

        <!-- Unauthenticated Navigation -->
        <div v-else class="flex items-center space-x-4">
          <RouterLink to="/login">
            <Button variant="ghost">Sign In</Button>
          </RouterLink>
          <RouterLink to="/signup">
            <Button>Get Started</Button>
          </RouterLink>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const authStore = useAuthStore()
const userStore = useUserStore()

/**
 * Initialize user profile when component mounts
 */
onMounted(async () => {
  if (authStore.isAuthenticated && !userStore.hasProfile) {
    await userStore.fetchProfile()
  }
})

/**
 * Watch for auth state changes
 */
watch(() => authStore.isAuthenticated, async (isAuth) => {
  if (isAuth && !userStore.hasProfile) {
    await userStore.fetchProfile()
  }
})

/**
 * Handle user sign out
 */
const handleSignOut = async () => {
  await authStore.signOut()
  router.push('/login')
}
</script> 

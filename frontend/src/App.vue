<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink, RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const authStore = useAuthStore()

/**
 * Handle user sign out
 */
const handleSignOut = async () => {
  await authStore.signOut()
  router.push('/login')
}

/**
 * Initialize auth on app mount
 */
onMounted(() => {
  authStore.initialize()
})
</script>

<template>
  <div class="min-h-screen bg-seasalt">
    <!-- Navigation (only show when authenticated) -->
    <nav v-if="authStore.isAuthenticated" class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <RouterLink to="/" class="text-xl font-bold text-rich-black">
              PetCare
            </RouterLink>
          </div>
          
          <div class="flex items-center space-x-8">
            <RouterLink 
              to="/" 
              class="text-gray-700 hover:text-aquamarine transition-colors"
              :class="{ 'text-aquamarine': $route.name === 'home' }"
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
              <span class="text-sm text-gray-600">
                {{ authStore.user?.email }}
              </span>
              <Button variant="ghost" size="sm" @click="handleSignOut">
                Sign out
              </Button>
            </div>
          </div>
        </div>
      </div>
    </nav>

    <!-- Main Content -->
    <main :class="authStore.isAuthenticated ? 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8' : ''">
      <!-- Loading state -->
      <div v-if="authStore.loading" class="min-h-screen flex items-center justify-center">
        <div class="text-center">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-aquamarine mx-auto"></div>
          <p class="mt-2 text-gray-600">Loading...</p>
        </div>
      </div>
      
      <!-- Router content -->
      <RouterView v-else />
    </main>
  </div>
</template>



<script setup lang="ts">
import { onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import Navigation from '@/components/Navigation.vue'

const authStore = useAuthStore()
const userStore = useUserStore()

/**
 * Initialize auth on app mount
 */
onMounted(async () => {
  await authStore.initialize()
  // Initialize user profile if already authenticated
  if (authStore.isAuthenticated) {
    await userStore.initialize()
  }
})
</script>

<template>
  <div class="min-h-screen bg-seasalt">
    <!-- Navigation -->
    <Navigation />

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



<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useUserStore } from '@/stores/user'
import ClientDashboard from './ClientDashboard.vue'
import VeterinarianDashboard from './VeterinarianDashboard.vue'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'

const userStore = useUserStore()

/**
 * Determine which dashboard to show based on user role
 */
const dashboardComponent = computed(() => {
  if (userStore.isClient) {
    return ClientDashboard
  } else if (userStore.isVeterinarian) {
    return VeterinarianDashboard
  } else if (userStore.isAdmin) {
    // For now, show veterinarian dashboard for admin - we can create AdminDashboard later
    return VeterinarianDashboard
  }
  return null
})
</script>

<template>
  <div>
    <!-- Show role-specific dashboard if user has a profile -->
    <component 
      v-if="userStore.hasProfile && dashboardComponent" 
      :is="dashboardComponent" 
    />
    
    <!-- Fallback for users without profiles or unknown roles -->
    <div v-else-if="!userStore.hasProfile" class="space-y-6">
      <Card class="p-8 text-center">
        <h2 class="text-2xl font-semibold text-rich-black mb-4">Welcome to PetCare!</h2>
        <p class="text-gray-600 mb-6">Please complete your profile setup to access your dashboard.</p>
        <RouterLink to="/profile-setup">
          <Button>Complete Profile Setup</Button>
        </RouterLink>
      </Card>
    </div>
    
    <!-- Fallback for unknown roles -->
    <div v-else class="space-y-6">
      <Card class="p-8 text-center">
        <h2 class="text-2xl font-semibold text-rich-black mb-4">Dashboard Loading...</h2>
        <p class="text-gray-600">Setting up your personalized dashboard...</p>
      </Card>
    </div>
  </div>
</template>

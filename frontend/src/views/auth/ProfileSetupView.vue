<template>
  <div class="min-h-screen flex items-center justify-center bg-seasalt py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold text-rich-black">
          Complete Your Profile
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Tell us a bit about yourself to get started
        </p>
      </div>

      <Card class="p-8">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <div>
            <Label for="name">Full Name</Label>
            <Input
              id="name"
              name="name"
              type="text"
              required
              placeholder="Enter your full name"
              v-model="form.name"
              class="mt-1"
            />
          </div>

          <div>
            <Label for="role">I am a</Label>
            <select
              id="role"
              name="role"
              required
              v-model="form.role"
              class="mt-1 flex h-10 w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aquamarine focus-visible:ring-offset-2"
            >
              <option value="">Select your role</option>
              <option value="client">Pet Owner</option>
              <option value="veterinarian">Veterinarian</option>
            </select>
          </div>

          <div>
            <Label for="phone">Phone Number</Label>
            <Input
              id="phone"
              name="phone"
              type="tel"
              placeholder="Enter your phone number"
              v-model="form.phone"
              class="mt-1"
            />
          </div>

          <div v-if="form.role === 'client'">
            <Label for="address">Address</Label>
            <textarea
              id="address"
              name="address"
              rows="3"
              placeholder="Enter your address"
              v-model="form.address"
              class="mt-1 flex w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aquamarine focus-visible:ring-offset-2"
            />
          </div>

          <div v-if="form.role === 'veterinarian'">
            <Label for="clinic_address">Clinic Address</Label>
            <textarea
              id="clinic_address"
              name="clinic_address"
              rows="3"
              placeholder="Enter your clinic address"
              v-model="form.clinic_address"
              class="mt-1 flex w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aquamarine focus-visible:ring-offset-2"
            />
          </div>

          <div v-if="error" class="text-red-600 text-sm">
            {{ error }}
          </div>

          <div>
            <Button
              type="submit"
              :disabled="loading"
              class="w-full"
            >
              {{ loading ? 'Creating Profile...' : 'Complete Setup' }}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'

const router = useRouter()
const authStore = useAuthStore()
const userStore = useUserStore()

const loading = ref(false)
const error = ref('')

const form = ref({
  name: '',
  role: '',
  phone: '',
  address: '',
  clinic_address: ''
})

/**
 * Handle form submission
 */
const handleSubmit = async () => {
  if (!form.value.name || !form.value.role) {
    error.value = 'Please fill in all required fields'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const profileData = {
      name: form.value.name,
      email: authStore.user?.email || '',
      phone: form.value.phone,
      role: form.value.role as 'client' | 'veterinarian'
    }

    // Add role-specific fields
    if (form.value.role === 'client' && form.value.address) {
      (profileData as any).address = form.value.address
    } else if (form.value.role === 'veterinarian' && form.value.clinic_address) {
      (profileData as any).clinic_address = form.value.clinic_address
    }

    await userStore.createProfile(profileData)
    
    // Refresh the user profile to update the navigation bar immediately
    await userStore.fetchProfile()
    
    // Redirect to dashboard
    router.push('/dashboard')
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to create profile'
  } finally {
    loading.value = false
  }
}
</script> 

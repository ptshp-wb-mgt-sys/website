<template>
  <div class="min-h-screen flex items-center justify-center bg-seasalt py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold text-rich-black">
          Sign in to PetCare
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or
          <RouterLink to="/signup" class="font-medium text-aquamarine hover:text-aquamarine-600">
            create a new account
          </RouterLink>
        </p>
      </div>

      <Card class="p-8">
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <Label for="email">Email address</Label>
            <Input
              id="email"
              name="email"
              type="email"
              autocomplete="email"
              required
              placeholder="Enter your email"
              v-model="email"
              class="mt-1"
            />
          </div>

          <div>
            <Label for="password">Password</Label>
            <Input
              id="password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              placeholder="Enter your password"
              v-model="password"
              class="mt-1"
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
              {{ loading ? 'Signing in...' : 'Sign in' }}
            </Button>
          </div>
        </form>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import { RouterLink } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

/**
 * Handle login form submission
 */
const handleLogin = async () => {
  if (!email.value || !password.value) {
    error.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const { error: loginError } = await authStore.signIn(email.value, password.value)
    
    if (loginError) {
      error.value = loginError.message
    } else {
      router.push('/')
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}
</script> 

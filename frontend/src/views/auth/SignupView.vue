<template>
  <div class="min-h-screen flex items-center justify-center bg-seasalt py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div>
        <h2 class="mt-6 text-center text-3xl font-bold text-rich-black">
          Create your PetCare account
        </h2>
        <p class="mt-2 text-center text-sm text-gray-600">
          Or
          <RouterLink to="/login" class="font-medium text-aquamarine hover:text-aquamarine-600">
            sign in to your account
          </RouterLink>
        </p>
      </div>

      <Card class="p-8">
        <form @submit.prevent="handleSignup" class="space-y-6">
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
              autocomplete="new-password"
              required
              placeholder="Create a password"
              v-model="password"
              class="mt-1"
            />
          </div>

          <div>
            <Label for="confirmPassword">Confirm Password</Label>
            <Input
              id="confirmPassword"
              name="confirmPassword"
              type="password"
              autocomplete="new-password"
              required
              placeholder="Confirm your password"
              v-model="confirmPassword"
              class="mt-1"
            />
          </div>

          <div v-if="error" class="text-red-600 text-sm">
            {{ error }}
          </div>

          <div v-if="success" class="text-green-600 text-sm">
            {{ success }}
          </div>

          <div>
            <Button
              type="submit"
              :disabled="loading"
              class="w-full"
            >
              {{ loading ? 'Creating account...' : 'Create account' }}
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
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')
const success = ref('')

/**
 * Handle signup form submission
 */
const handleSignup = async () => {
  if (!email.value || !password.value || !confirmPassword.value) {
    error.value = 'Please fill in all fields'
    return
  }

  if (password.value !== confirmPassword.value) {
    error.value = 'Passwords do not match'
    return
  }

  if (password.value.length < 6) {
    error.value = 'Password must be at least 6 characters'
    return
  }

  loading.value = true
  error.value = ''
  success.value = ''

  try {
    const { error: signupError } = await authStore.signUp(email.value, password.value)
    
    if (signupError) {
      error.value = signupError.message
    } else {
      success.value = 'Account created! Please check your email to verify your account.'
      setTimeout(() => {
        router.push('/login')
      }, 2000)
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}
</script> 

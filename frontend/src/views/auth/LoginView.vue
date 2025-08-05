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
            <div class="relative mt-1">
              <Input
                id="password"
                name="password"
                :type="showPassword ? 'text' : 'password'"
                autocomplete="current-password"
                required
                placeholder="Enter your password"
                v-model="password"
                class="pr-10"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600"
              >
                <Eye v-if="!showPassword" class="w-5 h-5" />
                <EyeOff v-else class="w-5 h-5" />
              </button>
            </div>
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
import { Eye, EyeOff } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import { RouterLink } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

/**
 * Handle login form submission
 */
const handleLogin = async () => {
  // Trim whitespace and check values
  const emailValue = email.value?.trim() || ''
  const passwordValue = password.value?.trim() || ''
  
  if (!emailValue || !passwordValue) {
    error.value = 'Please fill in all fields'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const { error: loginError } = await authStore.signIn(emailValue, passwordValue)
    
    if (loginError) {
      error.value = loginError.message
    } else {
      router.push('/dashboard')
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}
</script> 

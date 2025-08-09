<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">Profile</h1>
      <Button v-if="!editing" @click="editing = true">Edit Profile</Button>
      <div v-else class="space-x-2">
        <Button variant="outline" @click="cancelEdit">Cancel</Button>
        <Button @click="saveProfile" :disabled="saving">
          {{ saving ? 'Saving...' : 'Save' }}
        </Button>
      </div>
    </div>

    <!-- Profile Content -->
    <Card class="p-6">
      <div v-if="userStore.loading" class="text-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-aquamarine mx-auto"></div>
        <p class="mt-2 text-gray-600">Loading profile...</p>
      </div>

      <div v-else-if="userStore.error" class="text-center py-8">
        <p class="text-red-600">{{ userStore.error }}</p>
        <Button variant="outline" @click="userStore.fetchProfile" class="mt-4">
          Retry
        </Button>
      </div>

      <div v-else-if="userStore.profile" class="space-y-6">
        <!-- Display Mode -->
        <div v-if="!editing" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <h3 class="text-sm font-medium text-gray-500">Full Name</h3>
              <p class="mt-1 text-lg text-rich-black">{{ userStore.profile.name }}</p>
            </div>
            
            <div>
              <h3 class="text-sm font-medium text-gray-500">Role</h3>
              <p class="mt-1 text-lg text-rich-black">{{ userStore.roleDisplayName }}</p>
            </div>
            
            <div>
              <h3 class="text-sm font-medium text-gray-500">Email</h3>
              <p class="mt-1 text-lg text-rich-black">{{ userStore.profile.email }}</p>
            </div>
            
            <div v-if="userStore.profile.phone">
              <h3 class="text-sm font-medium text-gray-500">Phone</h3>
              <p class="mt-1 text-lg text-rich-black">{{ userStore.profile.phone }}</p>
            </div>
          </div>

          <div v-if="userStore.profile.address && userStore.isClient" class="mt-6">
            <h3 class="text-sm font-medium text-gray-500">Address</h3>
            <p class="mt-1 text-lg text-rich-black">{{ userStore.profile.address }}</p>
          </div>

          <div v-if="(userStore.profile as any).clinic_address && userStore.isVeterinarian" class="mt-6">
            <h3 class="text-sm font-medium text-gray-500">Clinic Address</h3>
            <p class="mt-1 text-lg text-rich-black">{{ (userStore.profile as any).clinic_address }}</p>
          </div>
        </div>

        <!-- Edit Mode -->
        <form v-else @submit.prevent="saveProfile" class="space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <Label for="name">Full Name</Label>
              <Input
                id="name"
                v-model="editForm.name"
                required
                class="mt-1"
              />
            </div>
            
            <div>
              <Label for="email">Email</Label>
              <Input
                id="email"
                v-model="editForm.email"
                type="email"
                required
                class="mt-1"
              />
            </div>

            <div>
              <Label for="phone">Phone Number</Label>
              <Input
                id="phone"
                v-model="editForm.phone"
                type="tel"
                class="mt-1"
              />
            </div>
          </div>

          <div v-if="userStore.isClient">
            <Label for="address">Address</Label>
            <textarea
              id="address"
              v-model="editForm.address"
              rows="3"
              class="mt-1 flex w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aquamarine focus-visible:ring-offset-2"
            />
          </div>

          <div v-if="userStore.isVeterinarian">
            <Label for="clinic_address">Clinic Address</Label>
            <textarea
              id="clinic_address"
              v-model="editForm.clinic_address"
              rows="3"
              class="mt-1 flex w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-aquamarine focus-visible:ring-offset-2"
            />
          </div>

          <div v-if="error" class="text-red-600 text-sm">
            {{ error }}
          </div>
        </form>
      </div>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'

const userStore = useUserStore()

const editing = ref(false)
const saving = ref(false)
const error = ref('')

const editForm = ref({
  name: '',
  email: '',
  phone: '',
  address: '',
  clinic_address: ''
})

/**
 * Initialize edit form with current profile data
 */
const initEditForm = () => {
  if (userStore.profile) {
    editForm.value = {
      name: userStore.profile.name || '',
      email: userStore.profile.email || '',
      phone: userStore.profile.phone || '',
      address: userStore.profile.address || '',
      clinic_address: (userStore.profile as any).clinic_address || ''
    }
  }
}

/**
 * Watch for profile changes to update edit form
 */
watch(() => userStore.profile, initEditForm, { immediate: true })

/**
 * Cancel editing
 */
const cancelEdit = () => {
  editing.value = false
  error.value = ''
  initEditForm()
}

/**
 * Save profile changes
 */
const saveProfile = async () => {
  if (!editForm.value.name) {
    error.value = 'Name is required'
    return
  }

  saving.value = true
  error.value = ''

  try {
    const updates: any = {
      name: editForm.value.name,
      email: editForm.value.email,
      phone: editForm.value.phone
    }

    if (userStore.isClient) {
      updates.address = editForm.value.address
    } else if (userStore.isVeterinarian) {
      updates.clinic_address = editForm.value.clinic_address
    }

    await userStore.updateProfile(updates)
    editing.value = false
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to save profile'
  } finally {
    saving.value = false
  }
}
</script> 

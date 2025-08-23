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
        <!-- Avatar + Upload (coming soon) -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-20 h-20 rounded-full bg-gray-100 border flex items-center justify-center text-2xl text-gray-400 select-none">
              {{ getInitial(userStore.profile.name) }}
            </div>
            <div>
              <div class="text-xl font-semibold text-rich-black">{{ userStore.profile.name }}</div>
              <div class="text-gray-500">{{ userStore.roleDisplayName }}</div>
            </div>
          </div>
          <div class="text-right">
            <Button variant="outline" size="sm" disabled>
              <Camera class="w-4 h-4 mr-1" /> Upload Photo
            </Button>
            <p class="text-xs text-gray-400 mt-1">Coming soon</p>
          </div>
        </div>
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

          <div v-if="userStore.isClient" class="space-y-2">
            <Label>Address</Label>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <Input
                  id="address_line1"
                  v-model="clientAddress.line1"
                  placeholder="Street address"
                  class="mt-1"
                />
              </div>
              <Input
                id="address_city"
                v-model="clientAddress.city"
                placeholder="City"
                class="mt-1"
              />
              <Input
                id="address_state"
                v-model="clientAddress.state"
                placeholder="State/Province"
                class="mt-1"
              />
              <Input
                id="address_postal"
                v-model="clientAddress.postalCode"
                placeholder="Postal code"
                class="mt-1"
              />
              <Input
                id="address_country"
                v-model="clientAddress.country"
                placeholder="Country"
                class="mt-1"
              />
            </div>
          </div>

          <div v-if="userStore.isVeterinarian" class="space-y-2">
            <Label>Clinic Address</Label>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="md:col-span-2">
                <Input
                  id="clinic_address_line1"
                  v-model="clinicAddress.line1"
                  placeholder="Street address"
                  class="mt-1"
                />
              </div>
              <Input
                id="clinic_address_city"
                v-model="clinicAddress.city"
                placeholder="City"
                class="mt-1"
              />
              <Input
                id="clinic_address_state"
                v-model="clinicAddress.state"
                placeholder="State/Province"
                class="mt-1"
              />
              <Input
                id="clinic_address_postal"
                v-model="clinicAddress.postalCode"
                placeholder="Postal code"
                class="mt-1"
              />
              <Input
                id="clinic_address_country"
                v-model="clinicAddress.country"
                placeholder="Country"
                class="mt-1"
              />
            </div>
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
import { Camera } from 'lucide-vue-next'
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

// Structured address state for better UX
const clientAddress = ref({ line1: '', city: '', state: '', postalCode: '', country: '' })
const clinicAddress = ref({ line1: '', city: '', state: '', postalCode: '', country: '' })

/**
 * formatAddress builds a simple single-line address from parts.
 */
const formatAddress = (addr: { line1?: string; city?: string; state?: string; postalCode?: string; country?: string }) => {
  const parts = [addr.line1, addr.city, addr.state, addr.postalCode, addr.country]
  return parts.filter(Boolean).join(', ')
}

/**
 * parseAddress splits a single string into basic parts. Best effort.
 */
const parseAddress = (value: string) => {
  const parts = (value || '').split(',').map(p => p.trim()).filter(Boolean)
  return {
    line1: parts[0] || '',
    city: parts[1] || '',
    state: parts[2] || '',
    postalCode: parts[3] || '',
    country: parts[4] || ''
  }
}

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

    // hydrate structured address inputs
    clientAddress.value = parseAddress(editForm.value.address)
    clinicAddress.value = parseAddress(editForm.value.clinic_address)
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
      updates.address = formatAddress(clientAddress.value)
    } else if (userStore.isVeterinarian) {
      updates.clinic_address = formatAddress(clinicAddress.value)
    }

    await userStore.updateProfile(updates)
    editing.value = false
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to save profile'
  } finally {
    saving.value = false
  }
}

/**
 * getInitial returns the first letter of a name in uppercase.
 */
const getInitial = (name: string): string => {
  const n = (name || '').trim()
  return n ? n.charAt(0).toUpperCase() : '?'
}
</script> 

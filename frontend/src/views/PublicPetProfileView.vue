<template>
  <div class="min-h-screen bg-gray-50 py-12">
    <div class="max-w-xl mx-auto bg-white shadow rounded-lg p-6">
      <div class="mb-6 text-center">
        <h1 class="text-2xl font-bold text-rich-black">Pet Profile</h1>
        <p class="text-gray-600">Public information for identification</p>
      </div>

      <div v-if="loading" class="py-10 text-center text-gray-600">
        <Loader2 class="w-6 h-6 animate-spin mx-auto mb-2" />
        Loading...
      </div>

      <div v-else-if="error" class="py-6 text-center text-red-600">{{ error }}</div>

      <div v-else-if="!profile" class="py-6 text-center text-gray-600">
        No profile found.
      </div>

      <div v-else class="space-y-6">
        <!-- Header with avatar placeholder and upload button -->
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-4">
            <div class="w-20 h-20 rounded-full bg-gray-100 border flex items-center justify-center text-2xl text-gray-400 select-none">
              {{ getInitial(profile.pet_name) }}
            </div>
            <div>
              <h2 class="text-xl font-semibold text-rich-black">{{ profile.pet_name }}</h2>
              <p class="text-gray-500 capitalize">{{ profile.pet_type }}</p>
            </div>
          </div>
          <div class="text-right">
            <Button variant="outline" size="sm" disabled>
              <Camera class="w-4 h-4 mr-1" /> Upload Photo
            </Button>
            <p class="text-xs text-gray-400 mt-1">Coming soon</p>
          </div>
        </div>

        <div class="grid grid-cols-1 gap-3">
          <div>
            <span class="font-medium">Pet:</span>
            <span class="ml-1">{{ profile.pet_name }}</span>
          </div>
          <div>
            <span class="font-medium">Type:</span>
            <span class="ml-1">{{ profile.pet_type }}</span>
          </div>
          <div v-if="profile.breed">
            <span class="font-medium">Breed:</span>
            <span class="ml-1">{{ profile.breed }}</span>
          </div>
          <div v-if="profile.date_of_birth">
            <span class="font-medium">DOB:</span>
            <span class="ml-1">{{ profile.date_of_birth }}</span>
          </div>
          <div v-if="profile.weight !== undefined">
            <span class="font-medium">Weight:</span>
            <span class="ml-1">{{ profile.weight }} kg</span>
          </div>
          <div>
            <span class="font-medium">Owner:</span>
            <span class="ml-1">{{ profile.owner_name }}</span>
          </div>
          <div>
            <span class="font-medium">Phone:</span>
            <a class="ml-1 text-aquamarine hover:underline" :href="`tel:${profile.owner_phone}`">{{ profile.owner_phone }}</a>
          </div>
          <div>
            <span class="font-medium">Email:</span>
            <a class="ml-1 text-aquamarine hover:underline" :href="`mailto:${profile.owner_email}`">{{ profile.owner_email }}</a>
          </div>
          <div>
            <span class="font-medium">Address:</span>
            <span class="ml-1">{{ profile.owner_address }}</span>
          </div>
          <div v-if="profile.emergency_contact">
            <span class="font-medium">Emergency Contact:</span>
            <span class="ml-1">{{ profile.emergency_contact }}</span>
          </div>
          <div v-if="(profile.medical_alerts || []).length">
            <span class="font-medium">Medical Alerts:</span>
            <ul class="list-disc ml-6 mt-1">
              <li v-for="(a, i) in profile.medical_alerts" :key="i">{{ a }}</li>
            </ul>
          </div>
          <div v-if="(profile.medical_records || []).length" class="pt-2">
            <h3 class="text-md font-semibold text-rich-black mb-1">Recent Medical Records</h3>
            <div class="space-y-2">
              <div v-for="(rec, i) in profile.medical_records" :key="i" class="border rounded p-3 bg-gray-50">
                <div class="text-sm text-gray-700"><span class="font-medium">Date:</span> {{ rec.date_of_visit }}</div>
                <div class="text-sm text-gray-700"><span class="font-medium">Reason:</span> {{ rec.reason_for_visit }}</div>
                <div v-if="rec.diagnosis" class="text-sm text-gray-700"><span class="font-medium">Diagnosis:</span> {{ rec.diagnosis }}</div>
                <div v-if="(rec.medication_prescribed || []).length" class="text-sm text-gray-700">
                  <span class="font-medium">Meds:</span>
                  <span class="ml-1">{{ (rec.medication_prescribed || []).join(', ') }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Loader2, Camera } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'

type PublicPetProfile = {
  pet_name: string
  pet_type: string
  breed?: string
  date_of_birth?: string
  weight?: number
  owner_name: string
  owner_phone: string
  owner_email: string
  owner_address: string
  emergency_contact?: string
  medical_alerts?: string[]
  medical_records?: Array<{
    date_of_visit: string
    reason_for_visit: string
    diagnosis?: string
    medication_prescribed?: string[]
  }>
}

const route = useRoute()
const loading = ref(false)
const error = ref<string | null>(null)
const profile = ref<PublicPetProfile | null>(null)

/**
 * Fetch the public profile by token. Works without auth.
 */
const fetchProfile = async () => {
  loading.value = true
  error.value = null
  try {
    const token = route.params.publicUrl as string
    // The backend accepts both full path and token; we send token only
    const res = await fetch(`http://localhost:3000/api/v1/public/pets/${encodeURIComponent(token)}`)
    if (!res.ok) {
      const body = await res.json().catch(() => ({} as any))
      throw new Error(body?.error || body?.message || res.statusText)
    }
    const body = await res.json()
    profile.value = (body.data || body) as PublicPetProfile
  } catch (e: any) {
    error.value = e?.message || 'Failed to load profile'
  } finally {
    loading.value = false
  }
}

onMounted(fetchProfile)

/**
 * Get first initial of the pet name for the avatar placeholder.
 */
const getInitial = (name: string): string => {
  const n = (name || '').trim()
  return n ? n.charAt(0).toUpperCase() : '?'
}
</script>

<style scoped>
</style>



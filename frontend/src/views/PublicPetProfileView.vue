<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-50 to-white py-10">
    <div class="max-w-5xl mx-auto px-4">
      <!-- Hero -->
      <div class="bg-white shadow-lg rounded-2xl p-8 mb-8">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-6">
          <div class="flex items-center gap-5">
            <div class="w-24 h-24 rounded-full bg-gray-100 border flex items-center justify-center text-3xl text-gray-400 select-none">
              {{ getInitial(profile?.pet_name || '') }}
            </div>
            <div>
              <div class="flex items-center gap-3">
                <h1 class="text-2xl font-bold text-rich-black">Pet Profile</h1>
                <span class="text-xs px-2 py-1 rounded-full bg-aquamarine/10 text-aquamarine">Public</span>
              </div>
              <div class="mt-2">
                <h2 class="text-xl font-semibold text-rich-black">{{ profile?.pet_name }}</h2>
                <p class="text-gray-500 capitalize">{{ profile?.pet_type }}</p>
              </div>
            </div>
          </div>
          <div class="text-right">
            <Button variant="outline" size="sm" disabled>
              <Camera class="w-4 h-4 mr-1" /> Upload Photo
            </Button>
            <p class="text-xs text-gray-400 mt-1">Coming soon</p>
          </div>
        </div>
      </div>

      <!-- Content -->
      <div class="bg-white shadow rounded-2xl p-6">
        <div v-if="loading" class="py-16 text-center text-gray-600">
          <Loader2 class="w-8 h-8 animate-spin mx-auto mb-3" />
          Loading profile...
        </div>

        <div v-else-if="error" class="py-10 text-center text-red-600">{{ error }}</div>

        <div v-else-if="!profile" class="py-10 text-center text-gray-600">
          No profile found.
        </div>

        <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Left: Owner & Contact -->
          <div class="space-y-6">
            <div class="rounded-xl border bg-white p-5">
              <h3 class="text-sm font-semibold text-gray-500 mb-3">Owner</h3>
              <div class="space-y-2 text-rich-black">
                <div class="text-lg font-medium">{{ profile.owner_name }}</div>
                <div class="flex items-center gap-2 text-gray-600">
                  <Phone class="w-4 h-4" />
                  <a :href="`tel:${profile.owner_phone}`" class="text-aquamarine hover:underline">{{ profile.owner_phone }}</a>
                </div>
                <div class="flex items-center gap-2 text-gray-600">
                  <Mail class="w-4 h-4" />
                  <a :href="`mailto:${profile.owner_email}`" class="text-aquamarine hover:underline">{{ profile.owner_email }}</a>
                </div>
                <div class="flex items-start gap-2 text-gray-600">
                  <MapPin class="w-4 h-4 mt-0.5" />
                  <span>{{ profile.owner_address }}</span>
                </div>
              </div>
            </div>

            <div v-if="profile.emergency_contact" class="rounded-xl border bg-white p-5">
              <h3 class="text-sm font-semibold text-gray-500 mb-3">Emergency Contact</h3>
              <div class="text-rich-black">{{ profile.emergency_contact }}</div>
            </div>

            <div v-if="(profile.medical_alerts || []).length" class="rounded-xl border bg-white p-5">
              <h3 class="text-sm font-semibold text-gray-500 mb-3 flex items-center gap-2">
                <ShieldAlert class="w-4 h-4" /> Alerts
              </h3>
              <div class="flex flex-wrap gap-2">
                <span
                  v-for="(a, i) in profile.medical_alerts"
                  :key="i"
                  class="px-2 py-1 text-xs rounded-full bg-red-50 text-red-700 border border-red-100"
                >{{ a }}</span>
              </div>
            </div>
          </div>

          <!-- Right: Pet details & Records -->
          <div class="lg:col-span-2 space-y-6">
            <div class="rounded-xl border bg-white p-5">
              <h3 class="text-sm font-semibold text-gray-500 mb-4">Pet Details</h3>
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-rich-black">
                <div>
                  <div class="text-sm text-gray-500">Name</div>
                  <div class="text-base">{{ profile.pet_name }}</div>
                </div>
                <div>
                  <div class="text-sm text-gray-500">Type</div>
                  <div class="text-base capitalize">{{ profile.pet_type }}</div>
                </div>
                <div v-if="profile.breed">
                  <div class="text-sm text-gray-500">Breed</div>
                  <div class="text-base">{{ profile.breed }}</div>
                </div>
                <div v-if="profile.date_of_birth">
                  <div class="text-sm text-gray-500">DOB</div>
                  <div class="text-base">{{ profile.date_of_birth }}</div>
                </div>
                <div v-if="profile.weight !== undefined">
                  <div class="text-sm text-gray-500">Weight</div>
                  <div class="text-base">{{ profile.weight }} kg</div>
                </div>
              </div>
            </div>

            <div v-if="(profile.medical_records || []).length" class="rounded-xl border bg-white p-5">
              <h3 class="text-sm font-semibold text-gray-500 mb-4 flex items-center gap-2">
                <Calendar class="w-4 h-4" /> Recent Medical Records
              </h3>
              <div class="relative">
                <div class="absolute left-3 top-0 bottom-0 w-px bg-gray-200"></div>
                <div class="space-y-5">
                  <div
                    v-for="(rec, i) in profile.medical_records"
                    :key="i"
                    class="relative pl-10"
                  >
                    <div class="absolute left-1.5 top-1 w-3 h-3 rounded-full bg-aquamarine border-2 border-white shadow"></div>
                    <div class="rounded-lg border bg-gray-50 p-4">
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
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { Loader2, Camera, Phone, Mail, MapPin, ShieldAlert, Calendar } from 'lucide-vue-next'
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



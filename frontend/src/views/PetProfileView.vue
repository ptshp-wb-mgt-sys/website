<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-rich-black">{{ pet?.name || 'Pet' }}</h1>
        <p class="text-gray-600 mt-1">Profile and medical history</p>
      </div>
      <div class="flex space-x-2" v-if="isClient">
        <Button variant="outline" size="sm" @click="goBack">Back</Button>
      </div>
    </div>

    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="text-center">
        <Loader2 class="w-8 h-8 animate-spin text-aquamarine mx-auto mb-4" />
        <p class="text-gray-600">Loading profile...</p>
      </div>
    </div>

    <div v-else-if="error" class="text-center py-12">
      <div class="text-red-600 mb-4">
        <AlertCircle class="w-12 h-12 mx-auto mb-2" />
        <p>{{ error }}</p>
      </div>
      <Button variant="outline" @click="initialize">Try Again</Button>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
      <!-- Left: Pet details -->
      <Card class="p-6 lg:col-span-1">
        <div class="space-y-3">
          <h2 class="text-lg font-semibold text-rich-black">Details</h2>
          <div class="text-sm space-y-1">
            <p><span class="font-medium">Type:</span> {{ pet?.type }}</p>
            <p><span class="font-medium">Breed:</span> {{ pet?.breed }}</p>
            <p><span class="font-medium">DOB:</span> {{ pet?.date_of_birth }}</p>
            <p><span class="font-medium">Weight:</span> {{ pet?.weight }} kg</p>
          </div>
          <div class="pt-3 border-t flex items-center gap-2">
            <Button variant="outline" size="sm" @click="openQRModal"><QrCode class="w-4 h-4 mr-1" /> QR Code</Button>
            <Button variant="ghost" size="sm" @click="goToPublic">Public Profile</Button>
            <Button v-if="isClient || isVet" variant="ghost" size="sm" @click="openEdit">Edit</Button>
          </div>
        </div>
      </Card>

      <!-- Right: Medical records -->
      <div class="space-y-4 lg:col-span-2">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold text-rich-black">Medical Records</h2>
          <Button v-if="isVet" size="sm" @click="showAdd = true">Add Record</Button>
        </div>

        <Card v-if="recordsLoading" class="p-6">
          <div class="flex items-center space-x-3 text-gray-600">
            <Loader2 class="w-5 h-5 animate-spin" />
            <span>Loading records...</span>
          </div>
        </Card>

        <template v-else>
          <Card v-if="records.length === 0" class="p-6 text-gray-600">No medical records yet.</Card>
          <div v-else class="space-y-3">
            <Card v-for="rec in recordsSorted" :key="rec.id" class="p-4">
              <div class="flex items-start justify-between">
                <div>
                  <p class="font-medium text-rich-black">{{ formatDate(rec.date_of_visit) }} — {{ rec.reason_for_visit }}</p>
                  <p class="text-sm text-gray-700">Diagnosis: {{ rec.diagnosis || '—' }}</p>
                  <p class="text-sm text-gray-700" v-if="rec.medication_prescribed?.length">
                    Meds: {{ rec.medication_prescribed.join(', ') }}
                  </p>
                  <p class="text-sm text-gray-600 mt-1" v-if="rec.notes">{{ rec.notes }}</p>
                </div>
                <div class="flex items-center space-x-2" v-if="isVet">
                  <Button variant="outline" size="sm" @click="startEdit(rec)">Edit</Button>
                  <Button variant="ghost" size="sm" class="text-red-600" @click="remove(rec)">Delete</Button>
                </div>
              </div>
            </Card>
          </div>
        </template>

        <!-- Add/Edit Modal -->
        <Card v-if="showAdd || editing" class="p-6">
          <h3 class="text-md font-semibold text-rich-black mb-3">{{ editing ? 'Edit Record' : 'Add Record' }}</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Date of visit</label>
              <input v-model="form.date_of_visit" type="date" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Reason</label>
              <input v-model="form.reason_for_visit" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">Diagnosis</label>
              <input v-model="form.diagnosis" type="text" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">Medication (comma separated)</label>
              <input v-model="medicationsText" type="text" placeholder="e.g. Amoxicillin, Prednisone" class="w-full px-3 py-2 border border-gray-300 rounded-lg" />
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">Notes</label>
              <textarea v-model="form.notes" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg"></textarea>
            </div>
          </div>
          <div class="mt-4 flex items-center space-x-2">
            <Button @click="save" :disabled="recordsLoading">Save</Button>
            <Button variant="ghost" @click="cancelEdit">Cancel</Button>
          </div>
        </Card>
      </div>
    </div>
  </div>

  <!-- QR Preview Modal -->
  <QRPreviewModal :is-open="showQRModal" :pet-id="id" @close="() => (showQRModal = false)" />
  
  <!-- Edit Pet Modal -->
  <EditPetModal :is-open="showEditPet" :pet="pet" @close="() => (showEditPet = false)" @pet-updated="initialize" />
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { usePetsStore, type Pet } from '@/stores/pets'
import { useMedicalRecordsStore, type MedicalRecord, type CreateMedicalRecordRequest, type UpdateMedicalRecordRequest } from '@/stores/medicalRecords'
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
import { Loader2, AlertCircle, QrCode } from 'lucide-vue-next'
import QRPreviewModal from '@/components/QRPreviewModal.vue'
import { useQRCodesStore } from '@/stores/qrcodes'
import EditPetModal from '@/components/EditPetModal.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const petsStore = usePetsStore()
const recordsStore = useMedicalRecordsStore()
const qrStore = useQRCodesStore()

const id = route.params.id as string

const pet = ref<Pet | null>(null)
const loading = ref(false)
const error = ref<string | null>(null)

const recordsLoading = computed(() => recordsStore.loading)
const records = computed<MedicalRecord[]>(() => recordsStore.getCachedForPet(id))
const recordsSorted = computed(() => {
  return [...records.value].sort((a, b) => new Date(b.date_of_visit).getTime() - new Date(a.date_of_visit).getTime())
})

const isClient = computed(() => userStore.isClient)
const isVet = computed(() => userStore.isVeterinarian)

const showAdd = ref(false)
const editing = ref<MedicalRecord | null>(null)
const medicationsText = ref('')
const showQRModal = ref(false)
const showEditPet = ref(false)

const form = ref<CreateMedicalRecordRequest | UpdateMedicalRecordRequest>({
  date_of_visit: '',
  reason_for_visit: '',
  diagnosis: '',
  medication_prescribed: [],
  notes: '',
})

/**
 * Fetch pet details and medical records.
 */
const initialize = async () => {
  loading.value = true
  error.value = null
  try {
    if (petsStore.pets.length === 0) {
      await petsStore.fetchPets()
    }
    // Prefer cache, fallback to API
    pet.value = petsStore.pets.find(p => p.id === id) || (await petsStore.getPet(id))
    await recordsStore.fetchByPetId(id)
  } catch (err) {
    error.value = err instanceof Error ? err.message : 'Failed to load pet profile'
  } finally {
    loading.value = false
  }
}

/**
 * Format a date string to a nice short format.
 */
const formatDate = (isoOrDate: string) => {
  const d = new Date(isoOrDate)
  return d.toLocaleDateString()
}

/**
 * Start editing an existing record.
 */
const startEdit = (rec: MedicalRecord) => {
  editing.value = rec
  showAdd.value = false
  form.value = {
    date_of_visit: rec.date_of_visit?.slice(0, 10),
    reason_for_visit: rec.reason_for_visit,
    diagnosis: rec.diagnosis,
    medication_prescribed: rec.medication_prescribed || [],
    notes: rec.notes,
  }
  medicationsText.value = (rec.medication_prescribed || []).join(', ')
}

/**
 * Save new or edited record.
 */
const save = async () => {
  if (!isVet.value) return
  const meds = medicationsText.value
    .split(',')
    .map(m => m.trim())
    .filter(Boolean)

  // Normalize date to ISO if provided as YYYY-MM-DD
  const normalizeDate = (d?: string) => {
    if (!d) return undefined
    if (d.includes('T')) return d
    try {
      return new Date(d + 'T00:00:00').toISOString()
    } catch {
      return undefined
    }
  }

  if (editing.value) {
    const updates: UpdateMedicalRecordRequest = {
      date_of_visit: normalizeDate(form.value.date_of_visit as string | undefined),
      reason_for_visit: form.value.reason_for_visit,
      diagnosis: form.value.diagnosis,
      medication_prescribed: meds,
      notes: form.value.notes,
    }
    await recordsStore.updateRecord(editing.value.id, updates)
  } else {
    const payload: CreateMedicalRecordRequest = {
      date_of_visit: normalizeDate(form.value.date_of_visit as string | undefined),
      reason_for_visit: form.value.reason_for_visit || '',
      diagnosis: form.value.diagnosis || '',
      medication_prescribed: meds,
      notes: form.value.notes || '',
    }
    await recordsStore.createRecord(id, payload)
  }

  cancelEdit()
}

/**
 * Cancel add/edit mode and reset state.
 */
const cancelEdit = () => {
  editing.value = null
  showAdd.value = false
  form.value = { date_of_visit: '', reason_for_visit: '', diagnosis: '', medication_prescribed: [], notes: '' }
  medicationsText.value = ''
}

/**
 * Delete a record after confirmation.
 */
const remove = async (rec: MedicalRecord) => {
  if (!isVet.value) return
  if (!confirm('Delete this medical record?')) return
  await recordsStore.deleteRecord(rec.id, id)
}

/**
 * Navigate back to previous page.
 */
const goBack = () => {
  router.back()
}

onMounted(async () => {
  await initialize()
})

/** Open QR modal (ensure QR exists) */
const openQRModal = async () => {
  try {
    await qrStore.getOrCreateForPet(id)
    showQRModal.value = true
  } catch (e) {
    console.error('Failed to load QR', e)
  }
}

/** Open Edit Pet modal */
const openEdit = () => {
  showEditPet.value = true
}

/**
 * Navigate to the public profile page (ensures a QR exists to get the URL).
 */
const goToPublic = async () => {
  try {
    const rec = await qrStore.getOrCreateForPet(id)
    const url = rec.encoded_content?.public_profile_url
    if (!url) return
    // Extract token from stored public path and route within SPA
    const token = url.split('/').pop() as string
    router.push(`/public/pets/${token}`)
  } catch (e) {
    console.error('Failed to open public profile', e)
  }
}
</script>

<style scoped>
</style>



<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-6 w-full max-w-md mx-4">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-xl font-semibold text-rich-black">Edit Pet</h2>
        <button
          @click="closeModal"
          class="text-gray-400 hover:text-gray-600 transition-colors"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <Label for="name">Pet Name</Label>
          <Input
            id="name"
            v-model="form.name"
            type="text"
            placeholder="Enter pet name"
            required
          />
        </div>

        <div>
          <Label for="type">Species</Label>
          <select
            id="type"
            v-model="form.type"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-aquamarine focus:border-transparent"
            required
          >
            <option value="">Select species</option>
            <option value="dog">Dog</option>
            <option value="cat">Cat</option>
            <option value="bird">Bird</option>
            <option value="rabbit">Rabbit</option>
            <option value="hamster">Hamster</option>
            <option value="fish">Fish</option>
            <option value="other">Other</option>
          </select>
        </div>

        <div>
          <Label for="breed">Breed</Label>
          <Input
            id="breed"
            v-model="form.breed"
            type="text"
            placeholder="Enter breed"
            required
          />
        </div>

        <div>
          <Label for="date_of_birth">Date of Birth</Label>
          <Input
            id="date_of_birth"
            v-model="form.date_of_birth"
            type="date"
            required
          />
        </div>

        <div>
          <Label for="weight">Weight (kg)</Label>
          <Input
            id="weight"
            :model-value="form.weight.toString()"
            @update:model-value="(value) => form.weight = parseFloat(value) || 0"
            type="number"
            step="0.1"
            min="0"
            placeholder="Enter weight in kg"
            required
          />
        </div>

        <div class="flex space-x-3 pt-4">
          <Button
            type="button"
            variant="outline"
            @click="closeModal"
            class="flex-1"
            :disabled="loading"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            class="flex-1"
            :disabled="loading"
          >
            <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
            {{ loading ? 'Updating...' : 'Update Pet' }}
          </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { X, Loader2 } from 'lucide-vue-next'
import { usePetsStore, type UpdatePetRequest, type Pet } from '@/stores/pets'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'

interface Props {
  isOpen: boolean
  pet: Pet | null
}

interface Emits {
  (e: 'close'): void
  (e: 'pet-updated'): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const petsStore = usePetsStore()
const loading = ref(false)

const form = reactive<UpdatePetRequest>({
  name: '',
  type: '',
  breed: '',
  date_of_birth: '',
  weight: 0
})

/**
 * Watch for pet changes and populate form
 */
watch(() => props.pet, (pet) => {
  if (pet) {
    form.name = pet.name
    form.type = pet.type
    form.breed = pet.breed
    form.date_of_birth = pet.date_of_birth // Already in YYYY-MM-DD format
    form.weight = pet.weight
  }
}, { immediate: true })

/**
 * Close the modal and reset form
 */
const closeModal = () => {
  resetForm()
  emit('close')
}

/**
 * Reset form to initial state
 */
const resetForm = () => {
  form.name = ''
  form.type = ''
  form.breed = ''
  form.date_of_birth = ''
  form.weight = 0
}

/**
 * Handle form submission
 */
const handleSubmit = async () => {
  if (loading.value || !props.pet) return

  loading.value = true

  try {
    // Convert date string to RFC3339 format for backend
    const petData: UpdatePetRequest = {
      ...form,
      date_of_birth: new Date(form.date_of_birth).toISOString()
    }
    
    await petsStore.updatePet(props.pet.id, petData)
    emit('pet-updated')
    closeModal()
  } catch (error) {
    console.error('Error updating pet:', error)
    // Error is handled by the store
  } finally {
    loading.value = false
  }
}
</script> 

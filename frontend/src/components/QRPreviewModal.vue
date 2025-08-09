<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg w-full max-w-md mx-4 p-6">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-lg font-semibold text-rich-black">QR Code</h2>
        <button @click="$emit('close')" class="text-gray-400 hover:text-gray-600">
          <X class="w-5 h-5" />
        </button>
      </div>

      <div v-if="loading" class="py-10 text-center text-gray-600">
        <Loader2 class="w-6 h-6 animate-spin mx-auto mb-2" />
        Generating...
      </div>

      <div v-else-if="error" class="py-6 text-center text-red-600">{{ error }}</div>

      <div v-else class="space-y-4">
        <div class="flex flex-col items-center">
          <img :src="qrSrc" alt="QR Code" class="w-56 h-56" />
          <p class="text-sm text-gray-600 mt-3 whitespace-pre-line text-center">
            {{ previewText }}
          </p>
        </div>

        <div class="flex items-center gap-2 pt-2">
          <Button class="flex-1" @click="downloadPng">Save PNG</Button>
          <Button variant="outline" class="flex-1" @click="copyText">Copy Text</Button>
        </div>
      </div>
    </div>
  </div>
  
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { X, Loader2 } from 'lucide-vue-next'
import Button from '@/components/ui/Button.vue'
import { useQRCodesStore } from '@/stores/qrcodes'

interface Props {
  isOpen: boolean
  petId: string
}

const props = defineProps<Props>()
defineEmits<{ (e: 'close'): void }>()

const qr = useQRCodesStore()

const loading = computed(() => qr.loading)
const error = computed(() => qr.error)

const record = computed(() => qr.cacheByPetId[props.petId])
const qrSrc = computed(() => record.value ? qr.toImageSrc(record.value.qr_code_data) : '')

// Build a human-friendly preview from encoded content
const previewText = computed(() => {
  const c = record.value?.encoded_content
  if (!c) return ''
  return `Pet: ${c.pet_name}\nOwner: ${c.owner_name}\nPhone: ${c.owner_phone}\nAddress: ${c.owner_address}\nProfile: ${window.location.origin}${c.public_profile_url}`
})

/**
 * Download the QR image as .png
 */
const downloadPng = () => {
  const src = qrSrc.value
  if (!src) return
  const link = document.createElement('a')
  link.href = src
  link.download = `pet-qr-${props.petId}.png`
  document.body.appendChild(link)
  link.click()
  link.remove()
}

/**
 * Copy the embedded text to clipboard
 */
const copyText = async () => {
  try {
    await navigator.clipboard.writeText(previewText.value)
  } catch (e) {
    console.error('Copy failed', e)
  }
}

</script>



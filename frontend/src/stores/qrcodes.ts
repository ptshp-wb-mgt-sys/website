import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useAuthStore } from './auth'

export interface QRCodeEncodedContent {
  pet_name: string
  pet_type: string
  owner_name: string
  owner_phone: string
  owner_email: string
  owner_address: string
  emergency_contact?: string
  medical_alerts?: string[]
  public_profile_url: string
}

export interface QRCodeRecord {
  id: string
  pet_id: string
  qr_code_data: string // base64 encoded PNG
  public_url: string
  encoded_content: QRCodeEncodedContent
  is_active: boolean
  created_at: string
  updated_at: string
}

export const useQRCodesStore = defineStore('qrcodes', () => {
  const cacheByPetId = ref<Record<string, QRCodeRecord | undefined>>({})
  const loading = ref(false)
  const error = ref<string | null>(null)

  const auth = useAuthStore()

  /**
   * Get QR code for a pet; if none exists, generate one.
   */
  const getOrCreateForPet = async (petId: string): Promise<QRCodeRecord> => {
    loading.value = true
    error.value = null
    try {
      // Return from cache if present
      const cached = cacheByPetId.value[petId]
      if (cached) return cached

      // Try to fetch existing
      const existing = await getForPet(petId)
      cacheByPetId.value[petId] = existing
      return existing
    } catch (err: any) {
      // If not found, generate
      if (err && typeof err === 'object' && 'status' in err && (err as any).status === 404) {
        const created = await generateForPet(petId)
        cacheByPetId.value[petId] = created
        return created
      }
      // Unknown error
      error.value = err instanceof Error ? err.message : 'Failed to get QR code'
      throw err
    } finally {
      loading.value = false
    }
  }

  /**
   * Fetch existing QR code for the given pet ID.
   */
  const getForPet = async (petId: string): Promise<QRCodeRecord> => {
    if (!auth.session?.access_token) {
      throw new Error('No authentication token')
    }
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets/${petId}/qr-code`, {
      headers: {
        Authorization: `Bearer ${auth.session.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    if (!res.ok) {
      const err: any = new Error(await safeErrorMessage(res))
      ;(err as any).status = res.status
      throw err
    }
    const body = await res.json()
    return (body.data || body) as QRCodeRecord
  }

  /**
   * Generate a new QR code for the given pet ID.
   */
  const generateForPet = async (petId: string): Promise<QRCodeRecord> => {
    if (!auth.session?.access_token) {
      throw new Error('No authentication token')
    }
    const res = await fetch(`${import.meta.env.VITE_API_URL ?? 'http://localhost:3000'}/api/v1/pets/${petId}/qr-code`, {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${auth.session.access_token}`,
        'Content-Type': 'application/json',
      },
    })
    if (!res.ok) {
      const errMsg = await safeErrorMessage(res)
      throw new Error(errMsg)
    }
    const body = await res.json()
    return (body.data || body) as QRCodeRecord
  }

  /**
   * Build an `img` src for the QR image from a base64 PNG string.
   */
  const toImageSrc = (base64Png: string) => `data:image/png;base64,${base64Png}`

  /**
   * Extract error message from a Response if possible.
   */
  const safeErrorMessage = async (res: Response): Promise<string> => {
    try {
      const j = await res.json()
      return j?.error || j?.message || res.statusText
    } catch {
      return res.statusText
    }
  }

  return {
    // state
    loading,
    error,
    cacheByPetId,

    // actions
    getOrCreateForPet,
    getForPet,
    generateForPet,

    // utils
    toImageSrc,
  }
})



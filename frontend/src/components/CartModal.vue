<template>
  <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
    <div class="absolute inset-0 bg-black/40" @click="onClose" />
    <div class="relative z-10 w-full max-w-lg rounded-lg bg-white shadow-lg p-6 space-y-4">
      <div class="flex items-center justify-between">
        <h3 class="text-xl font-semibold text-rich-black">My Cart</h3>
        <button class="text-gray-500 hover:text-gray-700" @click="onClose">✕</button>
      </div>

      <div v-if="cart.items.length === 0" class="text-center text-gray-600 py-8">
        Your cart is empty.
      </div>

      <div v-else class="space-y-3 max-h-80 overflow-auto pr-1">
        <div v-for="it in cart.items" :key="it.productId" class="flex items-center gap-3">
          <div class="w-14 h-14 bg-gray-100 rounded-md overflow-hidden flex items-center justify-center">
            <img v-if="it.imageUrl" :src="it.imageUrl" alt="" class="w-full h-full object-cover" />
            <span v-else class="text-xs text-gray-400">No Image</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="truncate font-medium text-rich-black">{{ it.name }}</div>
            <div class="text-sm text-gray-600">{{ formatPHP(it.price) }} ×</div>
          </div>
          <input type="number" min="0" class="w-20 border rounded px-2 py-1" :value="it.quantity" @input="onQtyInput(it.productId, $event)" />
          <div class="w-24 text-right font-medium">{{ formatPHP(it.price * it.quantity) }}</div>
          <button class="ml-2 text-red-500 hover:text-red-600" @click="cart.removeItem(it.productId)">Remove</button>
        </div>
      </div>

      <div class="flex items-center justify-between pt-4 border-t">
        <div class="text-gray-700">Subtotal</div>
        <div class="text-lg font-semibold">{{ formatPHP(cart.subtotal) }}</div>
      </div>

      <div class="grid grid-cols-2 gap-3">
        <Button variant="outline" @click="cart.clear" :disabled="cart.items.length === 0">Clear</Button>
        <Button :disabled="cart.items.length === 0" @click="checkout">Checkout</Button>
      </div>

      <p class="text-xs text-gray-500 text-center">
        Note: Checkout is for demo purposes only; no real payment or fulfillment.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import { useCartStore } from '@/stores/cart'
import { formatPHP } from '@/lib/utils'
import { defineProps, defineEmits } from 'vue'

/** Simple cart modal. Emits `checkout` with items on confirm. */
const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{ (e: 'close'): void; (e: 'checkout', items: Array<{ productId: string; quantity: number }> ): void }>()

const cart = useCartStore()

function onClose(): void {
  emit('close')
}

function onQtyInput(productId: string, evt: Event): void {
  const target = evt.target as HTMLInputElement
  const value = parseInt(target.value, 10)
  cart.updateQuantity(productId, isFinite(value) ? value : 0)
}

function checkout(): void {
  const payload = cart.items.map(i => ({ productId: i.productId, quantity: i.quantity }))
  emit('checkout', payload)
  cart.clear()
  emit('close')
}
</script>



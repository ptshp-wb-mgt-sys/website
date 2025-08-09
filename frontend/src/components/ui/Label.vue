<template>
  <label
    :class="cn('text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70', extraClass)"
    v-bind="$attrs"
  >
    <slot />
  </label>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
import { computed, useAttrs } from 'vue'
import type { ClassValue } from 'clsx'

const attrs = useAttrs()
const extraClass = computed(() => attrs.class as ClassValue | undefined)
</script> 

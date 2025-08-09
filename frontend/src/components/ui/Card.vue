<template>
  <div 
    :class="cn('rounded-lg border bg-white shadow-sm', extraClass)"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
import { computed, useAttrs } from 'vue'
import type { ClassValue } from 'clsx'

const attrs = useAttrs()
const extraClass = computed(() => attrs.class as ClassValue | undefined)
</script> 

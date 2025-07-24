<template>
  <div 
    :class="cn('rounded-lg border bg-white shadow-sm', $attrs.class)"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>

<script setup lang="ts">
import { cn } from '@/lib/utils'
</script> 

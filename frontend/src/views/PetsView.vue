<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">My Pets</h1>
      <Button>Add New Pet</Button>
    </div>

    <!-- Pet Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card class="p-6">
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-rich-black">Buddy</h3>
            <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
              Dog
            </span>
          </div>
          
          <div class="space-y-2 text-sm">
            <p><span class="font-medium">Breed:</span> Golden Retriever</p>
            <p><span class="font-medium">Age:</span> 3 years</p>
            <p><span class="font-medium">Weight:</span> 30 kg</p>
          </div>
          
          <div class="flex space-x-2">
            <Button variant="outline" size="sm" class="flex-1">View QR</Button>
            <Button variant="ghost" size="sm" class="flex-1">Edit</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-rich-black">Luna</h3>
            <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
              Cat
            </span>
          </div>
          
          <div class="space-y-2 text-sm">
            <p><span class="font-medium">Breed:</span> Persian Cat</p>
            <p><span class="font-medium">Age:</span> 2 years</p>
            <p><span class="font-medium">Weight:</span> 4 kg</p>
          </div>
          
          <div class="flex space-x-2">
            <Button variant="outline" size="sm" class="flex-1">View QR</Button>
            <Button variant="ghost" size="sm" class="flex-1">Edit</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-rich-black">Max</h3>
            <span class="px-2 py-1 bg-aquamarine-100 text-aquamarine-800 text-xs rounded-full">
              Dog
            </span>
          </div>
          
          <div class="space-y-2 text-sm">
            <p><span class="font-medium">Breed:</span> German Shepherd</p>
            <p><span class="font-medium">Age:</span> 5 years</p>
            <p><span class="font-medium">Weight:</span> 35 kg</p>
          </div>
          
          <div class="flex space-x-2">
            <Button variant="outline" size="sm" class="flex-1">View QR</Button>
            <Button variant="ghost" size="sm" class="flex-1">Edit</Button>
          </div>
        </div>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import Button from '@/components/ui/Button.vue'
import Card from '@/components/ui/Card.vue'
</script> 

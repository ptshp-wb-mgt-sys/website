<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold text-rich-black">Pet Products</h1>
      <div class="flex space-x-2">
        <Button variant="outline">View Cart (3)</Button>
        <Button>Add Product</Button>
      </div>
    </div>

    <!-- Product Categories -->
    <div class="flex space-x-4 mb-6">
      <Button variant="ghost" size="sm" class="text-aquamarine">All</Button>
      <Button variant="ghost" size="sm">Food</Button>
      <Button variant="ghost" size="sm">Medicine</Button>
      <Button variant="ghost" size="sm">Accessories</Button>
      <Button variant="ghost" size="sm">Toys</Button>
    </div>

    <!-- Products Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Premium Dog Food</h3>
            <p class="text-sm text-gray-600">High-quality nutrition for adult dogs</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$45.99</span>
              <span class="text-sm text-gray-500">In Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Flea & Tick Treatment</h3>
            <p class="text-sm text-gray-600">Monthly protection for dogs and cats</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$29.99</span>
              <span class="text-sm text-gray-500">In Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Interactive Cat Toy</h3>
            <p class="text-sm text-gray-600">Keep your cat entertained and active</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$15.99</span>
              <span class="text-sm text-gray-500">In Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Dog Leash & Collar Set</h3>
            <p class="text-sm text-gray-600">Durable and comfortable walking gear</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$24.99</span>
              <span class="text-sm text-gray-500">In Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Pet Vitamins</h3>
            <p class="text-sm text-gray-600">Daily supplements for optimal health</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$19.99</span>
              <span class="text-sm text-orange-500">Low Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
          </div>
        </div>
      </Card>

      <Card class="p-6">
        <div class="space-y-4">
          <div class="aspect-square bg-gray-100 rounded-lg mb-4"></div>
          
          <div class="space-y-2">
            <h3 class="font-semibold text-rich-black">Pet Carrier</h3>
            <p class="text-sm text-gray-600">Safe and comfortable travel carrier</p>
            <div class="flex items-center justify-between">
              <span class="text-lg font-bold text-aquamarine">$59.99</span>
              <span class="text-sm text-gray-500">In Stock</span>
            </div>
          </div>
          
          <div class="flex space-x-2">
            <Button size="sm" class="flex-1">Add to Cart</Button>
            <Button variant="outline" size="sm">Details</Button>
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

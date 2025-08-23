<template>
  <section class="bg-seasalt py-10 w-full">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <Card class="p-8 md:p-10">
        <div class="text-center space-y-8">
          <!-- Intro -->
          <div class="space-y-2">
            <span class="inline-block text-xs md:text-sm uppercase tracking-wide text-rich-black/70 bg-gray-50 border border-gray-200 rounded px-2 py-0.5">About</span>
            <h1 class="text-4xl md:text-5xl font-bold text-rich-black">PetCare</h1>
            <p class="text-base md:text-lg text-black">The complete pet management solution for modern pet care.</p>
          </div>
          <div class="h-px bg-gray-100" aria-hidden="true"></div>

          <!-- Project Details -->
          <div class="space-y-3">
            <div class="flex flex-wrap items-center justify-center gap-2 md:gap-3">
              <span class="uppercase tracking-wide text-xs md:text-sm font-medium text-rich-black bg-gray-50 border border-gray-200 rounded px-3 py-1">Petshop Web‑Based Management System</span>
              <span class="text-sm md:text-base text-black">Zamboanga del Sur Provincial Government College (ZDSPGC)</span>
              <br>
              <span class="text-sm md:text-base text-black">BSIS — Section B</span>
            </div>
          </div>
          <div class="h-px bg-gray-100" aria-hidden="true"></div>

          <!-- Team -->
          <div class="space-y-4">
            <h2 class="text-lg font-semibold text-rich-black">Team Members</h2>
            <ul class="grid grid-cols-1 md:grid-cols-2 gap-2 mx-auto text-black text-base">
              <li class="py-1">Bañez Jerald</li>
              <li class="py-1">Palahang Renato</li>
              <li class="py-1">Cabasura Neneth</li>
              <li class="py-1">Dumam-ag Willymae</li>
            </ul>
          </div>

          <!-- Credits -->
          <div class="space-y-3 text-black">
            <p class="text-base md:text-lg text-black">
              fueled by <span aria-hidden="true">☕</span> made with <span aria-hidden="true">🤍</span> by 
              <a 
                href="https://github.com/mrjxtr" 
                target="_blank" 
                rel="noopener noreferrer"
                class="text-aquamarine hover:text-aquamarine-600 font-medium transition-colors"
              >
                @mrjxtr
              </a>
            </p>
            <p class="text-sm text-black">
              Developer: 
              <a 
                href="https://github.com/mrjxtr" 
                target="_blank" 
                rel="noopener noreferrer"
                class="text-aquamarine hover:text-aquamarine-600 font-medium transition-colors"
              >
                @mrjxtr
              </a>
            </p>
            <a 
              href="https://github.com/ptshp-wb-mgt-sys/website"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="View source code on GitHub"
              class="inline-flex items-center gap-1 h-10 px-4 py-2 rounded-md border border-gray-300 bg-transparent text-rich-black hover:bg-gray-50 transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2"
            >
              View Source Code →
            </a>
          </div>
        </div>
      </Card>
    </div>
  </section>
  
</template>

<script setup lang="ts">
import Card from '@/components/ui/Card.vue'

/**
 * About page: footer-like layout, richer details, and cleaner visual grouping.
 */
</script>



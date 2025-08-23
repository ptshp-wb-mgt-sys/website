<template>
  <nav class="bg-white border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex items-center">
          <RouterLink to="/" class="text-3xl font-bold text-rich-black">
            PetCare
          </RouterLink>
        </div>

        <!-- Mobile: Hamburger -->
        <div class="flex items-center md:hidden">
          <button
            type="button"
            class="inline-flex items-center justify-center rounded-md p-2 text-black hover:text-aquamarine focus:outline-none focus:ring-2 focus:ring-aquamarine"
            :aria-expanded="isMobileMenuOpen ? 'true' : 'false'"
            aria-controls="mobile-menu"
            @click="toggleMobileMenu"
          >
            <X v-if="isMobileMenuOpen" class="w-6 h-6" />
            <Menu v-else class="w-6 h-6" />
          </button>
        </div>

        <!-- Authenticated Navigation -->
        <div v-if="authStore.isAuthenticated" class="hidden md:flex items-center space-x-8">
          <RouterLink
            to="/dashboard"
            class="text-black hover:text-aquamarine transition-colors"
            :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'dashboard' }"
          >
            Dashboard
          </RouterLink>

          <!-- Client Navigation -->
          <template v-if="userStore.isClient">
            <RouterLink
              to="/my-pets"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'my-pets' }"
            >
              My Pets
            </RouterLink>
            <RouterLink
              to="/book-appointment"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'book-appointment' }"
            >
              Book Appointment
            </RouterLink>
            <RouterLink
              to="/browse-products"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'browse-products' }"
            >
              Browse Products
            </RouterLink>
            <RouterLink
              to="/about"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'about' }"
            >
              About
            </RouterLink>
          </template>

          <!-- Veterinarian Navigation -->
          <template v-if="userStore.isVeterinarian">
            <RouterLink
              to="/my-schedule"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'my-schedule' }"
            >
              My Schedule
            </RouterLink>
            <RouterLink
              to="/patients"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'patients' }"
            >
              Patients
            </RouterLink>
            <RouterLink
              to="/manage-products"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'manage-products' }"
            >
              Manage Products
            </RouterLink>
            <RouterLink
              to="/about"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'about' }"
            >
              About
            </RouterLink>
          </template>

          <!-- Admin Navigation -->
          <template v-if="userStore.isAdmin">
            <RouterLink
              to="/users"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'users' }"
            >
              Users
            </RouterLink>
            <RouterLink
              to="/all-pets"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'all-pets' }"
            >
              All Pets
            </RouterLink>
            <RouterLink
              to="/all-appointments"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'all-appointments' }"
            >
              All Appointments
            </RouterLink>
            <RouterLink
              to="/analytics"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'analytics' }"
            >
              Analytics
            </RouterLink>
            <RouterLink
              to="/about"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'about' }"
            >
              About
            </RouterLink>
          </template>

          <!-- Trailing: Icons (rightmost) -->
          <div class="flex items-center space-x-4">
            <!-- <RouterLink
              to="/about"
              class="text-black hover:text-aquamarine transition-colors"
              :class="{ 'underline decoration-aquamarine decoration-2 underline-offset-8': $route.name === 'about' }"
            >
              About
            </RouterLink> -->
            <a
              href="https://github.com/ptshp-wb-mgt-sys/website"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Project GitHub Repository"
              title="Project GitHub Repository"
              class="text-black hover:text-aquamarine transition-colors inline-flex items-center"
            >
              <Github class="w-5 h-5" />
            </a>
            <a
              href="https://github.com/mrjxtr"
              target="_blank"
              rel="noopener noreferrer"
              aria-label="Developer (@mrjxtr)"
              title="Developer (@mrjxtr)"
              class="text-black hover:text-aquamarine transition-colors inline-flex items-center"
            >
              <Code2 class="w-5 h-5" />
            </a>
          </div>

          <!-- User Menu -->
          <div class="flex items-center space-x-4">
            <RouterLink
              to="/profile"
              class="text-sm text-gray-400 hover:text-aquamarine transition-colors"
              :class="{ 'text-aquamarine': $route.name === 'profile' }"
            >
              <span v-if="userStore.hasProfile">{{ userStore.fullDisplayName }}</span>
              <span v-else-if="authStore.user?.email">{{ authStore.user.email }}</span>
              <span v-else>User</span>
            </RouterLink>
            <Button variant="ghost" size="sm" @click="handleSignOut"> Sign out </Button>
          </div>
        </div>

        <!-- Unauthenticated Navigation -->
        <div v-else class="hidden md:flex items-center space-x-4">
          <RouterLink to="/login">
            <Button variant="ghost">Sign In</Button>
          </RouterLink>
          <RouterLink to="/signup">
            <Button>Get Started</Button>
          </RouterLink>
          <RouterLink to="/about" class="text-black hover:text-aquamarine transition-colors">About</RouterLink>
          <a
            href="https://github.com/ptshp-wb-mgt-sys/website"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Project GitHub Repository"
            title="Project GitHub Repository"
            class="text-black hover:text-aquamarine transition-colors inline-flex items-center"
          >
            <Github class="w-5 h-5" />
          </a>
          <a
            href="https://github.com/mrjxtr"
            target="_blank"
            rel="noopener noreferrer"
            aria-label="Developer (@mrjxtr)"
            title="Developer (@mrjxtr)"
            class="text-black hover:text-aquamarine transition-colors inline-flex items-center"
          >
            <Code2 class="w-5 h-5" />
          </a>
        </div>
      </div>
      
      <!-- Mobile Menu Panel -->
      <div
        id="mobile-menu"
        class="md:hidden"
        v-if="isMobileMenuOpen"
      >
        <div class="space-y-1 pb-3 pt-2">
          <template v-if="authStore.isAuthenticated">
            <RouterLink
              to="/dashboard"
              class="block px-3 py-2 text-black hover:text-aquamarine"
              @click="closeMobileMenu"
            >
              Dashboard
            </RouterLink>

            <!-- Client Links -->
            <template v-if="userStore.isClient">
              <RouterLink to="/my-pets" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">My Pets</RouterLink>
              <RouterLink to="/book-appointment" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Book Appointment</RouterLink>
              <RouterLink to="/browse-products" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Browse Products</RouterLink>
              <RouterLink to="/about" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">About</RouterLink>
            </template>

            <!-- Veterinarian Links -->
            <template v-if="userStore.isVeterinarian">
              <RouterLink to="/my-schedule" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">My Schedule</RouterLink>
              <RouterLink to="/patients" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Patients</RouterLink>
              <RouterLink to="/manage-products" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Manage Products</RouterLink>
              <RouterLink to="/about" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">About</RouterLink>
            </template>

            <!-- Admin Links -->
            <template v-if="userStore.isAdmin">
              <RouterLink to="/users" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Users</RouterLink>
              <RouterLink to="/all-pets" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">All Pets</RouterLink>
              <RouterLink to="/all-appointments" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">All Appointments</RouterLink>
              <RouterLink to="/analytics" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">Analytics</RouterLink>
              <RouterLink to="/about" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">About</RouterLink>
            </template>

            <!-- Profile + Sign out -->
            <div class="border-t border-gray-200 mt-2 pt-2 px-3 space-y-2">
              <RouterLink to="/profile" class="block text-sm text-gray-600 hover:text-aquamarine" @click="closeMobileMenu">
                <span v-if="userStore.hasProfile">{{ userStore.fullDisplayName }}</span>
                <span v-else-if="authStore.user?.email">{{ authStore.user.email }}</span>
                <span v-else>User</span>
              </RouterLink>
              <Button variant="ghost" size="sm" class="w-full justify-start" @click="handleSignOut">Sign out</Button>
            </div>
          </template>

          <template v-else>
            <div class="px-3 space-y-2">
              <RouterLink to="/login" @click="closeMobileMenu">
                <Button variant="ghost" class="w-full justify-center">Sign In</Button>
              </RouterLink>
              <RouterLink to="/signup" @click="closeMobileMenu">
                <Button class="w-full justify-center">Get Started</Button>
              </RouterLink>
            </div>
            <RouterLink to="/about" class="block px-3 py-2 text-black hover:text-aquamarine" @click="closeMobileMenu">About</RouterLink>
            <div class="flex items-center gap-4 px-3 py-2">
              <a
                href="https://github.com/ptshp-wb-mgt-sys/website"
                target="_blank"
                rel="noopener noreferrer"
                aria-label="Project GitHub Repository"
                title="Project GitHub Repository"
                class="text-black hover:text-aquamarine inline-flex items-center"
              >
                <Github class="w-5 h-5" />
              </a>
              <a
                href="https://github.com/mrjxtr"
                target="_blank"
                rel="noopener noreferrer"
                aria-label="Developer (@mrjxtr)"
                title="Developer (@mrjxtr)"
                class="text-black hover:text-aquamarine inline-flex items-center"
              >
                <Code2 class="w-5 h-5" />
              </a>
            </div>
          </template>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { onMounted, watch, ref } from 'vue'
import { useRouter } from 'vue-router'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useUserStore } from '@/stores/user'
import Button from '@/components/ui/Button.vue'
import { Github, Code2, Menu, X } from 'lucide-vue-next'

const router = useRouter()
const authStore = useAuthStore()
const userStore = useUserStore()
const isMobileMenuOpen = ref(false)

/**
 * Initialize user profile when component mounts
 */
onMounted(async () => {
  if (authStore.isAuthenticated && !userStore.hasProfile) {
    await userStore.fetchProfile()
  }
})

/**
 * Watch for auth state changes
 */
watch(
  () => authStore.isAuthenticated,
  async (isAuth) => {
    if (isAuth && !userStore.hasProfile) {
      await userStore.fetchProfile()
    }
  },
)

/**
 * Toggle the mobile menu open/close — simple and snappy.
 */
const toggleMobileMenu = () => {
  isMobileMenuOpen.value = !isMobileMenuOpen.value
}

/**
 * Close the mobile menu, e.g., after navigating.
 */
const closeMobileMenu = () => {
  isMobileMenuOpen.value = false
}

/**
 * Auto-close mobile menu on route change to keep UI tidy.
 */
watch(
  () => router.currentRoute.value.fullPath,
  () => {
    isMobileMenuOpen.value = false
  },
)

/**
 * Handle user sign out
 */
const handleSignOut = async () => {
  await authStore.signOut()
  router.push('/login')
}
</script>

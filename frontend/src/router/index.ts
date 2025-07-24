import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: () => import('../views/auth/LoginView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('../views/auth/SignupView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/profile-setup',
      name: 'profile-setup',
      component: () => import('../views/auth/ProfileSetupView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/',
      name: 'landing',
      component: () => import('../views/LandingView.vue')
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/pets',
      name: 'pets',
      component: () => import('../views/PetsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/appointments',
      name: 'appointments',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/products',
      name: 'products',
      component: () => import('../views/ProductsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true }
    },
  ],
})

/**
 * Route guard to protect authenticated routes
 */
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Wait for auth to initialize if it hasn't yet
  if (authStore.loading) {
    await new Promise(resolve => {
      const unwatch = authStore.$subscribe(() => {
        if (!authStore.loading) {
          unwatch()
          resolve(true)
        }
      })
    })
  }

  const requiresAuth = to.meta.requiresAuth
  const requiresGuest = to.meta.requiresGuest
  const isAuthenticated = authStore.isAuthenticated

  if (requiresAuth && !isAuthenticated) {
    next('/login')
  } else if (requiresGuest && isAuthenticated) {
    next('/dashboard')
  } else if (isAuthenticated && to.name === 'landing') {
    // Redirect authenticated users away from landing page to dashboard
    next('/dashboard')
  } else if (isAuthenticated && requiresAuth && to.name !== 'profile-setup') {
    // Check if user has completed their profile
    const { useUserStore } = await import('@/stores/user')
    const userStore = useUserStore()
    
    if (!userStore.hasProfile) {
      await userStore.fetchProfile()
      if (!userStore.hasProfile) {
        next('/profile-setup')
        return
      }
    }
    next()
  } else {
    next()
  }
})

export default router

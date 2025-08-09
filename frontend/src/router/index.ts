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
      meta: { requiresGuest: true },
    },
    {
      path: '/signup',
      name: 'signup',
      component: () => import('../views/auth/SignupView.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/profile-setup',
      name: 'profile-setup',
      component: () => import('../views/auth/ProfileSetupView.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/',
      name: 'landing',
      component: () => import('../views/LandingView.vue'),
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: HomeView,
      meta: { requiresAuth: true },
    },
    // Client Routes
    {
      path: '/my-pets',
      name: 'my-pets',
      component: () => import('../views/PetsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['client'] },
    },
    {
      path: '/book-appointment',
      name: 'book-appointment',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['client'] },
    },
    {
      path: '/pet/:id',
      name: 'pet-profile',
      component: () => import('../views/PetProfileView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['client', 'veterinarian'] },
    },
    {
      path: '/browse-products',
      name: 'browse-products',
      component: () => import('../views/ProductsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['client'] },
    },

    // Veterinarian Routes
    {
      path: '/my-schedule',
      name: 'my-schedule',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['veterinarian'] },
    },
    {
      path: '/patients',
      name: 'patients',
      component: () => import('../views/PetsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['veterinarian'] },
    },
    {
      path: '/manage-products',
      name: 'manage-products',
      component: () => import('../views/ProductsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['veterinarian'] },
    },
    {
      path: '/medical-records',
      name: 'medical-records',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['veterinarian'] },
    },

    // Admin Routes
    {
      path: '/users',
      name: 'users',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['admin'] },
    },
    {
      path: '/all-pets',
      name: 'all-pets',
      component: () => import('../views/PetsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['admin'] },
    },
    {
      path: '/all-appointments',
      name: 'all-appointments',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['admin'] },
    },
    {
      path: '/analytics',
      name: 'analytics',
      component: () => import('../views/AppointmentsView.vue'),
      meta: { requiresAuth: true, allowedRoles: ['admin'] },
    },

    // Legacy redirects
    {
      path: '/pets',
      redirect: (to) => {
        // This will be handled by the route guard to redirect based on role
        return '/dashboard'
      },
    },
    {
      path: '/appointments',
      redirect: (to) => {
        return '/dashboard'
      },
    },
    {
      path: '/products',
      redirect: (to) => {
        return '/dashboard'
      },
    },
    {
      path: '/profile',
      name: 'profile',
      component: () => import('../views/ProfileView.vue'),
      meta: { requiresAuth: true },
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
    await new Promise((resolve) => {
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

    // Check role-based permissions
    const allowedRoles = to.meta.allowedRoles as string[] | undefined
    if (allowedRoles && userStore.profile?.role) {
      if (!allowedRoles.includes(userStore.profile.role)) {
        // User doesn't have permission for this route, redirect to dashboard
        console.warn(
          `User with role '${userStore.profile.role}' tried to access route requiring roles: ${allowedRoles.join(', ')}`,
        )
        next('/dashboard')
        return
      }
    }

    next()
  } else {
    next()
  }
})

export default router

# Authentication Setup Guide

## Supabase Configuration

To get authentication working, you need to:

### 1. Create a Supabase Project
1. Go to [supabase.com](https://supabase.com)
2. Sign up/Sign in and create a new project
3. Wait for the project to be ready

### 2. Get Your Credentials
1. Go to your project dashboard
2. Click on "Settings" in the sidebar
3. Click on "API"
4. Copy the following:
   - **Project URL** (under "Project URL")
   - **Anon (public) key** (under "Project API keys")

### 3. Set Up Environment Variables
1. Create a `.env.local` file in the frontend directory:
```bash
cp .env.example .env.local
```

2. Edit `.env.local` and replace the values:
```env
VITE_SUPABASE_URL=https://your-project-id.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key-here
```

### 4. Configure Authentication (Optional)
In your Supabase dashboard:
1. Go to "Authentication" → "Settings"
2. Configure your site URL (e.g., `http://localhost:5173`)
3. Set up email templates if desired

### 5. Test Authentication
1. Start the development server: `npm run dev`
2. Visit `http://localhost:5173`
3. You should be redirected to the login page
4. Try creating an account and signing in

## Features Included

✅ **User Registration** - Sign up with email/password  
✅ **User Login** - Sign in with email/password  
✅ **Protected Routes** - Dashboard, Pets, Appointments, Products require auth  
✅ **Route Guards** - Automatic redirects based on auth state  
✅ **User Session** - Persistent login across browser sessions  
✅ **Sign Out** - Clean logout functionality  
✅ **Loading States** - Smooth loading indicators  
✅ **Form Validation** - Client-side validation with error messages  

## Auth Flow

1. **Unauthenticated users** → Redirected to `/login`
2. **Successful login** → Redirected to dashboard (`/`)
3. **Authenticated users visiting `/login` or `/signup`** → Redirected to dashboard
4. **Sign out** → Redirected to login page

## Architecture

- **Auth Store** (`/stores/auth.ts`) - Pinia store managing auth state
- **Supabase Client** (`/lib/supabase.ts`) - Configured Supabase client
- **Route Guards** (`/router/index.ts`) - Protecting routes based on auth state
- **Auth Views** (`/views/auth/`) - Login and signup pages
- **UI Components** (`/components/ui/`) - Reusable form components

The frontend handles all authentication with Supabase directly - no backend integration needed for auth! 

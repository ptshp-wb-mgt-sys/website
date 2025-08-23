import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

/**
 * Utility for merging CSS classes with Tailwind
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * Format a number into Philippine Peso (PHP) currency.
 * Keeps it chill: just pass a number, get a `₱` string.
 */
export function formatPHP(amount: number): string {
  try {
    return new Intl.NumberFormat('en-PH', { style: 'currency', currency: 'PHP' }).format(amount)
  } catch {
    const safe = isFinite(amount) ? amount : 0
    return `₱${safe.toFixed(2)}`
  }
} 

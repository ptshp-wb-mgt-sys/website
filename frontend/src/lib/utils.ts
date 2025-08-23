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
    // Fallback: manual format with comma separators
    const parts = safe.toFixed(2).split('.')
    parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ',')
    return `₱${parts.join('.')}`
  }
} 

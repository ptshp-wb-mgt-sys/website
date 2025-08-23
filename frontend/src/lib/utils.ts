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

/**
 * formatTimeHM returns HH:mm for a given ISO datetime string.
 * Defaults to interpreting the input as UTC to avoid timezone drift across clients.
 */
export function formatTimeHM(iso: string, useUTC = true): string {
  try {
    const d = new Date(iso)
    return d.toLocaleTimeString([], {
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
      timeZone: useUTC ? 'UTC' : undefined,
    })
  } catch {
    return ''
  }
}

/**
 * formatDateTimeMDYHM returns "M/D/YYYY, HH:mm" for a given ISO datetime string.
 * By default uses the viewer's local timezone; pass useUTC=true to force UTC.
 */
export function formatDateTimeMDYHM(iso: string, useUTC = false): string {
  try {
    const d = new Date(iso)
    return d.toLocaleString([], {
      month: 'numeric',
      day: 'numeric',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
      timeZone: useUTC ? 'UTC' : undefined,
    })
  } catch {
    return ''
  }
}

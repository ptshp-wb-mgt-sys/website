/**
 * Returns the API base URL. Uses VITE_API_URL in prod; falls back to localhost for dev.
 */
export function getApiBaseUrl(): string {
  const fromEnv = import.meta.env.VITE_API_URL as string | undefined
  return fromEnv && fromEnv.length > 0 ? fromEnv : 'http://localhost:3000'
}

/**
 * Thin wrapper around fetch that prefixes paths with the configured API base URL.
 */
export function apiFetch(path: string, init?: RequestInit): Promise<Response> {
  const base = getApiBaseUrl()
  const url = path.startsWith('http') ? path : `${base}${path}`
  return fetch(url, init)
}



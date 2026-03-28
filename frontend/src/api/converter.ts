export type Direction = 'to_numeric' | 'to_string'

export interface ConversionResult {
  input: string
  type: 'ipv4' | 'ipv4_net' | 'ipv6' | 'ipv6_net' | 'error'
  output: string
  error?: string
}

export interface ConversionRequest {
  direction: Direction
  values: string[]
}

export interface ConversionResponse {
  results: ConversionResult[]
}

// In docker-compose, nginx proxies /api → backend. In dev, vite proxies it.
const API_BASE = '/api'

export async function convertValues(req: ConversionRequest): Promise<ConversionResponse> {
  const res = await fetch(`${API_BASE}/convert`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(req),
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`API error ${res.status}: ${text}`)
  }
  return res.json()
}

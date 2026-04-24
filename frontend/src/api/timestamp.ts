export type TsDirection = 'to_human' | 'to_unix'

export interface TsConversionResult {
  input: string
  type: 'timestamp' | 'error'
  output: string
  error?: string
}

export interface TsConversionRequest {
  direction: TsDirection
  values: string[]
}

export interface TsConversionResponse {
  results: TsConversionResult[]
}

const API_BASE = '/api'

export async function convertTimestamps(req: TsConversionRequest): Promise<TsConversionResponse> {
  const res = await fetch(`${API_BASE}/timestamp`, {
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

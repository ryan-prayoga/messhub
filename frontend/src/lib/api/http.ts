import type { ApiEnvelope } from '$lib/api/types';

export type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  token?: string;
  body?: Record<string, unknown>;
};

type ErrorPayloadShape =
  | {
      error?: string | { code?: string; details?: unknown };
      message?: string;
      details?: unknown;
    }
  | null
  | undefined;

export class ApiError extends Error {
  status: number;
  code?: string;
  details?: unknown;
  requestId?: string;
  isNetworkError: boolean;

  constructor(
    status: number,
    message: string,
    options: {
      code?: string;
      details?: unknown;
      requestId?: string;
      isNetworkError?: boolean;
    } = {}
  ) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = options.code;
    this.details = options.details;
    this.requestId = options.requestId;
    this.isNetworkError = options.isNetworkError ?? false;
  }
}

export function buildRequestInit(options: RequestOptions = {}): RequestInit {
  return {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  };
}

export async function parseApiResponse<T>(
  response: Response,
  fallbackMessage = 'Request failed'
): Promise<ApiEnvelope<T>> {
  const requestId = response.headers.get('X-Request-ID') ?? undefined;
  const payload = (await response.json().catch(() => null)) as ApiEnvelope<T> | ErrorPayloadShape;

  if (!response.ok) {
    throw new ApiError(response.status, payload?.message || fallbackMessage, {
      code: extractErrorCode(payload),
      details: extractErrorDetails(payload),
      requestId
    });
  }

  if (!payload || !('data' in payload)) {
    throw new ApiError(502, 'Invalid API response', {
      code: 'invalid_api_response',
      requestId
    });
  }

  return payload as ApiEnvelope<T>;
}

export function wrapNetworkError(error: unknown, fallbackMessage = 'Network request failed'): ApiError {
  if (error instanceof ApiError) {
    return error;
  }

  return new ApiError(503, error instanceof Error ? error.message : fallbackMessage, {
    code: 'network_error',
    isNetworkError: true
  });
}

export function isUnauthorizedError(error: unknown): error is ApiError {
  return error instanceof ApiError && error.status === 401;
}

export function isForbiddenError(error: unknown): error is ApiError {
  return error instanceof ApiError && error.status === 403;
}

function extractErrorCode(payload: ErrorPayloadShape) {
  if (!payload) {
    return undefined;
  }

  if (typeof payload.error === 'string') {
    return payload.error;
  }

  return payload.error?.code;
}

function extractErrorDetails(payload: ErrorPayloadShape) {
  if (!payload) {
    return undefined;
  }

  if (typeof payload.error === 'object' && payload.error?.details !== undefined) {
    return payload.error.details;
  }

  return payload.details;
}

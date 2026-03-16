import { env } from '$env/dynamic/private';
import type {
  ApiEnvelope,
  MemberUser,
  SessionUser,
  WifiBillDetail,
  WifiBillWithSummary,
  WifiMyBill,
  WalletSummary,
  WalletTransactionPage,
  WalletTransactionType
} from '$lib/api/types';

const DEFAULT_PRIVATE_API_BASE_URL = 'http://127.0.0.1:4100/api/v1';

type RequestOptions = {
  method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
  token?: string;
  body?: Record<string, unknown>;
};

export class ApiError extends Error {
  status: number;
  code?: string;

  constructor(status: number, message: string, code?: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
  }
}

function resolveApiUrl(path: string) {
  const baseUrl = env.PRIVATE_API_BASE_URL || DEFAULT_PRIVATE_API_BASE_URL;
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;

  return `${baseUrl}${normalizedPath}`;
}

export async function apiServerRequest<T>(
  fetcher: typeof fetch,
  path: string,
  options: RequestOptions = {}
): Promise<ApiEnvelope<T>> {
  const response = await fetcher(resolveApiUrl(path), {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  });

  const fallbackMessage = 'Request failed';
  const payload = (await response.json().catch(() => null)) as ApiEnvelope<T> | null;

  if (!response.ok) {
    throw new ApiError(
      response.status,
      payload?.message || fallbackMessage,
      payload?.error?.code
    );
  }

  if (!payload) {
    throw new ApiError(502, 'Invalid API response', 'invalid_api_response');
  }

  return payload;
}

export const authServerApi = {
  login: (fetcher: typeof fetch, payload: { email: string; password: string }) =>
    apiServerRequest<{ token: string; user: SessionUser }>(fetcher, '/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<SessionUser>(fetcher, '/auth/me', {
      token
    })
};

export const usersServerApi = {
  list: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<MemberUser[]>(fetcher, '/users', {
      token
    })
};

export const walletServerApi = {
  summary: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WalletSummary>(fetcher, '/wallet', {
      token
    }),
  listTransactions: (
    fetcher: typeof fetch,
    token: string,
    params: {
      page?: number;
      pageSize?: number;
    } = {}
  ) => {
    const searchParams = new URLSearchParams();

    if (params.page) {
      searchParams.set('page', String(params.page));
    }

    if (params.pageSize) {
      searchParams.set('page_size', String(params.pageSize));
    }

    const query = searchParams.toString();

    return apiServerRequest<WalletTransactionPage>(
      fetcher,
      `/wallet/transactions${query ? `?${query}` : ''}`,
      {
        token
      }
    );
  },
  createTransaction: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      type: WalletTransactionType;
      category: string;
      amount: number;
      description: string;
    }
  ) =>
    apiServerRequest(fetcher, '/wallet/transactions', {
      method: 'POST',
      token,
      body: payload
    })
};

export const wifiServerApi = {
  listBills: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiBillWithSummary[]>(fetcher, '/wifi/bills', {
      token
    }),
  getBill: (fetcher: typeof fetch, token: string, billID: string) =>
    apiServerRequest<WifiBillDetail>(fetcher, `/wifi/bills/${billID}`, {
      token
    }),
  getActive: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiBillDetail | null>(fetcher, '/wifi/active', {
      token
    }),
  getMyBills: (fetcher: typeof fetch, token: string) =>
    apiServerRequest<WifiMyBill[]>(fetcher, '/wifi/my', {
      token
    }),
  createBill: (
    fetcher: typeof fetch,
    token: string,
    payload: {
      month: number;
      year: number;
      nominal_per_person?: number;
      deadline_date?: string;
      status?: 'draft' | 'active' | 'closed';
    }
  ) =>
    apiServerRequest<WifiBillDetail>(fetcher, '/wifi/bills', {
      method: 'POST',
      token,
      body: payload
    }),
  submitProof: (
    fetcher: typeof fetch,
    token: string,
    billID: string,
    payload: {
      proof_url: string;
      note?: string;
    }
  ) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/submit`, {
      method: 'POST',
      token,
      body: payload
    }),
  verifyPayment: (fetcher: typeof fetch, token: string, billID: string, memberID: string) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/verify/${memberID}`, {
      method: 'PATCH',
      token
    }),
  rejectPayment: (
    fetcher: typeof fetch,
    token: string,
    billID: string,
    memberID: string,
    payload: {
      reason: string;
    }
  ) =>
    apiServerRequest(fetcher, `/wifi/bills/${billID}/reject/${memberID}`, {
      method: 'PATCH',
      token,
      body: payload
    })
};

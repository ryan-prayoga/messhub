import { API_BASE_URL } from '$lib/config/env';
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

export async function apiRequest<T>(
  path: string,
  options: RequestOptions = {}
): Promise<ApiEnvelope<T>> {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method || 'GET',
    headers: {
      'Content-Type': 'application/json',
      ...(options.token ? { Authorization: `Bearer ${options.token}` } : {})
    },
    body: options.body ? JSON.stringify(options.body) : undefined
  });

  const fallback = { message: 'Request failed' };
  const payload = (await response.json().catch(() => fallback)) as
    | (ApiEnvelope<T> & { message?: string })
    | { message?: string };

  if (!response.ok) {
    throw new ApiError(
      response.status,
      payload.message || fallback.message,
      'error' in payload ? payload.error?.code : undefined
    );
  }

  return payload as ApiEnvelope<T>;
}

export const authApi = {
  login: (payload: { email: string; password: string }) =>
    apiRequest<{ token: string; user: SessionUser }>('/auth/login', {
      method: 'POST',
      body: payload
    }),
  me: (token: string) => apiRequest<SessionUser>('/auth/me', { token })
};

export const usersApi = {
  list: (token: string) => apiRequest<MemberUser[]>('/users', { token })
};

export const walletApi = {
  summary: (token: string) => apiRequest<WalletSummary>('/wallet', { token }),
  listTransactions: (
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

    return apiRequest<WalletTransactionPage>(`/wallet/transactions${query ? `?${query}` : ''}`, {
      token
    });
  },
  createTransaction: (
    token: string,
    payload: {
      type: WalletTransactionType;
      category: string;
      amount: number;
      description: string;
    }
  ) =>
    apiRequest('/wallet/transactions', {
      method: 'POST',
      token,
      body: payload
    })
};

export const wifiApi = {
  listBills: (token: string) => apiRequest<WifiBillWithSummary[]>('/wifi/bills', { token }),
  getBill: (token: string, billID: string) => apiRequest<WifiBillDetail>(`/wifi/bills/${billID}`, { token }),
  getActive: (token: string) => apiRequest<WifiBillDetail | null>('/wifi/active', { token }),
  getMyBills: (token: string) => apiRequest<WifiMyBill[]>('/wifi/my', { token }),
  createBill: (
    token: string,
    payload: {
      month: number;
      year: number;
      nominal_per_person?: number;
      deadline_date?: string;
      status?: 'draft' | 'active' | 'closed';
    }
  ) =>
    apiRequest<WifiBillDetail>('/wifi/bills', {
      method: 'POST',
      token,
      body: payload
    }),
  submitProof: (
    token: string,
    billID: string,
    payload: {
      proof_url: string;
      note?: string;
    }
  ) =>
    apiRequest(`/wifi/bills/${billID}/submit`, {
      method: 'POST',
      token,
      body: payload
    }),
  verifyPayment: (token: string, billID: string, memberID: string) =>
    apiRequest(`/wifi/bills/${billID}/verify/${memberID}`, {
      method: 'PATCH',
      token
    }),
  rejectPayment: (
    token: string,
    billID: string,
    memberID: string,
    payload: {
      reason: string;
    }
  ) =>
    apiRequest(`/wifi/bills/${billID}/reject/${memberID}`, {
      method: 'PATCH',
      token,
      body: payload
    })
};

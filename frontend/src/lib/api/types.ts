export type UserRole = 'admin' | 'treasurer' | 'member';
export type WalletTransactionType = 'income' | 'expense';

export type SessionUser = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
};

export type MemberUser = SessionUser & {
  is_active: boolean;
  joined_at: string;
  left_at: string | null;
  created_at: string;
  updated_at: string;
};

export type WalletSummary = {
  balance: number;
  total_income: number;
  total_expense: number;
};

export type WalletTransaction = {
  id: string;
  type: WalletTransactionType;
  category: string;
  amount: number;
  description: string;
  created_by: string;
  created_by_name?: string;
  created_at: string;
  updated_at: string;
};

export type WalletTransactionPage = {
  items: WalletTransaction[];
  pagination: {
    page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
  };
};

export type ApiEnvelope<T> = {
  message: string;
  data: T;
  error?: {
    code?: string;
  };
};

export type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

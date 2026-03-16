export type UserRole = 'admin' | 'treasurer' | 'member';

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

export type ApiEnvelope<T> = {
  message: string;
  data: T;
  error?: {
    code?: string;
  };
};

export type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

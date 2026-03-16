export type UserRole = 'admin' | 'treasurer' | 'member';
export type ActivityType = 'contribution' | 'food' | 'rice' | 'announcement' | 'other';
export type WalletTransactionType = 'income' | 'expense';
export type WifiBillStatus = 'draft' | 'active' | 'closed';
export type WifiPaymentStatus =
  | 'unpaid'
  | 'pending_verification'
  | 'verified'
  | 'rejected';

export type SessionUser = {
  id: string;
  name: string;
  email: string;
  role: UserRole;
};

export type MemberUser = SessionUser & {
  phone: string | null;
  avatar_url: string | null;
  is_active: boolean;
  joined_at: string | null;
  left_at: string | null;
  created_at: string;
  updated_at: string;
};

export type Profile = MemberUser;

export type MessSettings = {
  id: string;
  mess_name: string;
  wifi_price: number;
  wifi_deadline_day: number;
  bank_account_name: string;
  bank_account_number: string;
  created_at: string;
  updated_at: string;
};

export type SystemStatus = {
  database_status: string;
  server_time: string;
  app_version: string;
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

export type WifiBillSummary = {
  total_members: number;
  verified_count: number;
  pending_count: number;
  unpaid_count: number;
  rejected_count: number;
  total_collected: number;
  total_target: number;
};

export type WifiBill = {
  id: string;
  month: number;
  year: number;
  nominal_per_person: number;
  deadline_date: string;
  status: WifiBillStatus;
  created_by: string;
  created_at: string;
  updated_at: string;
};

export type WifiBillMember = {
  id: string;
  wifi_bill_id: string;
  user_id: string;
  user_name: string;
  user_email: string;
  amount: number;
  payment_status: WifiPaymentStatus;
  proof_url: string | null;
  note: string | null;
  submitted_at: string | null;
  verified_at: string | null;
  verified_by: string | null;
  verified_by_name: string | null;
  rejection_reason: string | null;
  created_at: string;
  updated_at: string;
};

export type WifiBillWithSummary = WifiBill & {
  summary: WifiBillSummary;
};

export type WifiBillDetail = {
  bill: WifiBill;
  summary: WifiBillSummary;
  members: WifiBillMember[];
};

export type WifiMyBill = {
  member_id: string;
  wifi_bill_id: string;
  month: number;
  year: number;
  nominal_per_person: number;
  deadline_date: string;
  bill_status: WifiBillStatus;
  amount: number;
  payment_status: WifiPaymentStatus;
  proof_url: string | null;
  note: string | null;
  submitted_at: string | null;
  verified_at: string | null;
  rejection_reason: string | null;
  verified_by: string | null;
  verified_by_name: string | null;
  created_at: string;
  updated_at: string;
};

export type Activity = {
  id: string;
  type: ActivityType;
  title: string;
  content: string;
  points: number;
  user_id: string;
  user_name: string;
  created_by: string;
  created_by_name: string;
  expires_at: string | null;
  created_at: string;
  updated_at: string;
};

export type ActivityComment = {
  id: string;
  activity_id: string;
  user_id: string;
  user_name: string;
  comment: string;
  created_at: string;
  updated_at: string;
};

export type ActivityReactionSummary = {
  reaction_type: string;
  count: number;
  reacted: boolean;
};

export type FoodClaim = {
  id: string;
  activity_id: string;
  user_id: string;
  user_name: string;
  created_at: string;
};

export type RiceResponse = {
  id: string;
  activity_id: string;
  user_id: string;
  user_name: string;
  created_at: string;
};

export type ActivityFeedItem = {
  activity: Activity;
  comments: ActivityComment[];
  reactions: ActivityReactionSummary[];
  claims: FoodClaim[];
  rice_responses: RiceResponse[];
};

export type ContributionLeaderboardEntry = {
  rank: number;
  user_id: string;
  user_name: string;
  total_points: number;
  total_activities: number;
};

export type Notification = {
  id: string;
  user_id: string;
  title: string;
  message: string;
  type: string;
  entity_id: string | null;
  is_read: boolean;
  created_at: string;
};

export type NotificationList = {
  items: Notification[];
  unread_count: number;
};

export type ApiEnvelope<T> = {
  message: string;
  data: T;
  error?: {
    code?: string;
    details?: unknown;
  };
};

export type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

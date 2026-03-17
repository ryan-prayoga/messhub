export type UserRole = 'admin' | 'treasurer' | 'member';
export type ActivityType = 'contribution' | 'food' | 'rice' | 'announcement' | 'other';
export type WalletTransactionType = 'income' | 'expense';
export type WifiBillStatus = 'draft' | 'active' | 'closed';
export type FeedStatusFilter = 'active' | 'expired' | 'all';
export type WifiPaymentStatus =
  | 'unpaid'
  | 'pending_verification'
  | 'verified'
  | 'rejected';
export type SharedExpenseStatus =
  | 'personal'
  | 'fronted'
  | 'partially_reimbursed'
  | 'reimbursed';
export type ProposalStatus = 'active' | 'closed' | 'approved' | 'rejected';
export type ProposalVoteType = 'agree' | 'disagree';

export type SessionUser = {
  id: string;
  name: string;
  email: string;
  username: string;
  role: UserRole;
};

export type MemberUser = SessionUser & {
  phone: string | null;
  avatar_url: string | null;
  is_active: boolean;
  joined_at: string | null;
  left_at: string | null;
  archived_at: string | null;
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
  status: string;
  database_status: string;
  database_reachable: boolean;
  server_time: string;
  app_version: string;
  uptime_seconds: number;
};

export type WalletSummary = {
  balance: number;
  total_income: number;
  total_expense: number;
};

export type WalletTransaction = {
  id: string;
  transaction_date: string;
  type: WalletTransactionType;
  category: string;
  amount: number;
  description: string;
  proof_url: string | null;
  source: string;
  import_job_id: string | null;
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

export type SharedExpense = {
  id: string;
  expense_date: string;
  category: string;
  description: string;
  amount: number;
  paid_by_user_id: string;
  paid_by_user_name: string;
  status: SharedExpenseStatus;
  notes: string | null;
  proof_url: string | null;
  created_by: string;
  created_by_name: string;
  created_at: string;
  updated_at: string;
};

export type SharedExpenseSummary = {
  total_count: number;
  total_amount: number;
  fronted_count: number;
  outstanding_amount: number;
  this_month_amount: number;
};

export type SharedExpenseList = {
  items: SharedExpense[];
  summary: SharedExpenseSummary;
};

export type Proposal = {
  id: string;
  title: string;
  description: string;
  created_by: string;
  created_by_name: string;
  voting_start: string | null;
  voting_end: string | null;
  status: ProposalStatus;
  final_decision_by: string | null;
  final_decision_by_name: string | null;
  final_decision_note: string | null;
  agree_count: number;
  disagree_count: number;
  total_votes: number;
  current_user_vote: ProposalVoteType | null;
  created_at: string;
  updated_at: string;
};

export type ProposalVote = {
  id: string;
  proposal_id: string;
  user_id: string;
  user_name: string;
  vote_type: ProposalVoteType;
  created_at: string;
};

export type ProposalDetail = {
  proposal: Proposal;
  votes: ProposalVote[];
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
  error?: string | {
    code?: string;
    details?: unknown;
  };
  details?: unknown;
};

export type AuthStatus = 'loading' | 'authenticated' | 'unauthenticated';

export type ImportRowStatus = 'valid' | 'invalid' | 'duplicate';

export type ImportWarning = {
  code: string;
  message: string;
};

export type MemberImportPreviewRow = {
  row_number: number;
  status: ImportRowStatus;
  name: string;
  email: string;
  role: string;
  normalized_role: string;
  is_active: string;
  normalized_is_active?: boolean | null;
  errors: string[];
  warnings: string[];
};

export type MemberImportPreview = {
  job_id: string;
  file_name: string;
  summary: {
    total_rows: number;
    valid_rows: number;
    invalid_rows: number;
    duplicate_rows: number;
    importable_rows: number;
  };
  rows: MemberImportPreviewRow[];
  warnings: ImportWarning[];
  can_commit: boolean;
  requires_temporary_password: boolean;
};

export type WalletImportPreviewRow = {
  row_number: number;
  status: Exclude<ImportRowStatus, 'duplicate'>;
  transaction_date: string;
  normalized_transaction_date?: string | null;
  description: string;
  income: string;
  expense: string;
  type: WalletTransactionType | '';
  amount?: number | null;
  category: string;
  proof: string;
  errors: string[];
  warnings: string[];
};

export type WalletImportPreview = {
  job_id: string;
  file_name: string;
  summary: {
    total_rows: number;
    valid_rows: number;
    invalid_rows: number;
    importable_rows: number;
    total_income: number;
    total_expense: number;
  };
  rows: WalletImportPreviewRow[];
  warnings: ImportWarning[];
  can_commit: boolean;
};

export type ImportCommitResult = {
  job_id: string;
  import_type: 'members' | 'wallet';
  imported_rows: number;
  skipped_rows: number;
  failed_rows: number;
  total_rows: number;
  duplicate_strategy?: 'skip' | 'fail' | null;
  total_income?: number;
  total_expense?: number;
};

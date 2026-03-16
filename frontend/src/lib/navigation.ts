import { APP_NAME } from '$lib/config/env';
import type { SessionUser, UserRole } from '$lib/api/types';

export type NavigationSection = 'primary' | 'workspace' | 'admin';

export type NavigationItem = {
  href: string;
  label: string;
  title: string;
  description: string;
  icon: string;
  section: NavigationSection;
  roles?: UserRole[];
  showInBottomNav?: boolean;
};

type PageMeta = {
  title: string;
  description: string;
};

const navItems: NavigationItem[] = [
  {
    href: '/dashboard',
    label: 'Dashboard',
    title: 'Dashboard',
    description: 'Ringkasan operasional, akses cepat, dan status utama mess.',
    icon: 'lucide:layout-dashboard',
    section: 'primary',
    showInBottomNav: true
  },
  {
    href: '/wallet',
    label: 'Wallet',
    title: 'Wallet',
    description: 'Saldo kas, histori transaksi, dan akses pencatatan uang mess.',
    icon: 'lucide:wallet',
    section: 'primary',
    showInBottomNav: true
  },
  {
    href: '/wifi',
    label: 'Wifi',
    title: 'Wifi',
    description: 'Tagihan bulanan, bukti transfer, dan verifikasi pembayaran wifi.',
    icon: 'lucide:wifi',
    section: 'primary',
    showInBottomNav: true
  },
  {
    href: '/feed',
    label: 'Feed',
    title: 'Feed',
    description: 'Aktivitas harian mess, pengumuman singkat, dan interaksi penghuni.',
    icon: 'lucide:newspaper',
    section: 'primary',
    showInBottomNav: true
  },
  {
    href: '/members',
    label: 'Members',
    title: 'Members',
    description: 'Daftar penghuni, role, dan status aktif anggota mess.',
    icon: 'lucide:users',
    section: 'workspace',
    roles: ['admin', 'treasurer']
  },
  {
    href: '/notifications',
    label: 'Inbox',
    title: 'Inbox',
    description: 'Notifikasi terbaru untuk update wifi, feed, dan aktivitas penting.',
    icon: 'lucide:bell-ring',
    section: 'workspace'
  },
  {
    href: '/contributions',
    label: 'Contributions',
    title: 'Contributions',
    description: 'Papan kontribusi penghuni dan leaderboard aktivitas bersama.',
    icon: 'lucide:trophy',
    section: 'workspace'
  },
  {
    href: '/profile',
    label: 'Profile',
    title: 'Profile',
    description: 'Identitas akun, username, avatar, dan pengaturan sandi pribadi.',
    icon: 'lucide:user-round',
    section: 'workspace',
    showInBottomNav: true
  },
  {
    href: '/settings',
    label: 'Settings',
    title: 'Settings',
    description: 'Pengaturan mess, nominal wifi, dan kesiapan sistem utama.',
    icon: 'lucide:settings-2',
    section: 'admin',
    roles: ['admin']
  },
  {
    href: '/admin/import',
    label: 'Import',
    title: 'Admin',
    description: 'Pusat impor data anggota dan wallet untuk kebutuhan administrasi.',
    icon: 'lucide:download',
    section: 'admin',
    roles: ['admin']
  },
  {
    href: '/shared-expenses',
    label: 'Patungan',
    title: 'Shared Expenses',
    description: 'Pengeluaran non-kas bersama dan status penggantian penghuni.',
    icon: 'lucide:receipt',
    section: 'workspace'
  },
  {
    href: '/proposals',
    label: 'Usulan',
    title: 'Proposals',
    description: 'Usulan bersama dan voting sederhana untuk kebutuhan mess.',
    icon: 'lucide:file-text',
    section: 'workspace'
  }
];

const pageMetaOverrides: Array<{
  matcher: RegExp;
  meta: PageMeta;
}> = [
  {
    matcher: /^\/login$/,
    meta: {
      title: 'Login',
      description: 'Masuk ke MessHub dengan email atau username untuk membuka operasional harian mess.'
    }
  },
  {
    matcher: /^\/offline$/,
    meta: {
      title: 'Offline',
      description: 'Fallback saat koneksi belum tersedia dan halaman belum ada di cache lokal.'
    }
  },
  {
    matcher: /^\/admin\/import\/members(?:\/|$)/,
    meta: {
      title: 'Import Members',
      description: 'Pratinjau dan commit impor anggota dengan validasi yang aman.'
    }
  },
  {
    matcher: /^\/admin\/import\/wallet(?:\/|$)/,
    meta: {
      title: 'Import Wallet',
      description: 'Pratinjau dan commit impor transaksi wallet dengan validasi yang aman.'
    }
  },
  {
    matcher: /^\/wallet\/new(?:\/|$)/,
    meta: {
      title: 'Wallet',
      description: 'Buat transaksi wallet baru dengan detail yang rapi dan mudah ditinjau.'
    }
  }
];

export function buildPageTitle(title: string) {
  if (title === APP_NAME) {
    return APP_NAME;
  }

  return `${title} • ${APP_NAME}`;
}

export function isPathActive(pathname: string, href: string) {
  return pathname === href || pathname.startsWith(`${href}/`);
}

export function getVisibleNavigation(user: SessionUser | null | undefined) {
  const visibleItems = navItems.filter((item) => canAccess(item, user));

  return {
    primary: visibleItems.filter((item) => item.section === 'primary'),
    workspace: visibleItems.filter((item) => item.section === 'workspace'),
    admin: visibleItems.filter((item) => item.section === 'admin'),
    bottom: visibleItems.filter((item) => item.showInBottomNav)
  };
}

export function getCurrentNavigationItem(pathname: string, user: SessionUser | null | undefined) {
  const visibleItems = getVisibleNavigation(user);
  const orderedItems = [...visibleItems.admin, ...visibleItems.workspace, ...visibleItems.primary].sort(
    (left, right) => right.href.length - left.href.length
  );

  return orderedItems.find((item) => isPathActive(pathname, item.href)) ?? null;
}

export function getPageMeta(pathname: string, user: SessionUser | null | undefined): PageMeta {
  const override = pageMetaOverrides.find((entry) => entry.matcher.test(pathname));
  if (override) {
    return override.meta;
  }

  const item = getCurrentNavigationItem(pathname, user);
  if (item) {
    return {
      title: item.title,
      description: item.description
    };
  }

  return {
    title: APP_NAME,
    description: 'MessHub membantu operasional harian mess tetap rapi, hangat, dan mudah dipakai dari mana saja.'
  };
}

function canAccess(item: NavigationItem, user: SessionUser | null | undefined) {
  if (!item.roles || item.roles.length === 0) {
    return true;
  }

  if (!user) {
    return false;
  }

  return item.roles.includes(user.role);
}

import LayoutDashboardIcon from '@lucide/svelte/icons/layout-dashboard';
import BarChartIcon from '@lucide/svelte/icons/bar-chart';
import FileTextIcon from '@lucide/svelte/icons/file-text';
import FolderIcon from '@lucide/svelte/icons/folder';
import CheckSquareIcon from '@lucide/svelte/icons/check-square';
import CalendarIcon from '@lucide/svelte/icons/calendar';
import MessageSquareIcon from '@lucide/svelte/icons/message-square';
import UserIcon from '@lucide/svelte/icons/user';
import CreditCardIcon from '@lucide/svelte/icons/credit-card';
import SettingsIcon from '@lucide/svelte/icons/settings';
import type { Component } from 'svelte';

export interface NavItem {
  title: string;
  href: string;
  icon: Component;
}

export interface NavGroup {
  label: string;
  items: NavItem[];
}

export const navGroups: NavGroup[] = [
  {
    label: 'General',
    items: [
      {
        title: 'Dashboard',
        href: '/admin/dashboard',
        icon: LayoutDashboardIcon,
      },
      { title: 'Analytics', href: '/admin/analytics', icon: BarChartIcon },
      { title: 'Reports', href: '/admin/reports', icon: FileTextIcon },
    ],
  },
  {
    label: 'Manage',
    items: [
      { title: 'Projects', href: '/admin/projects', icon: FolderIcon },
      { title: 'Tasks', href: '/admin/tasks', icon: CheckSquareIcon },
      { title: 'Calendar', href: '/admin/calendar', icon: CalendarIcon },
      { title: 'Messages', href: '/admin/messages', icon: MessageSquareIcon },
    ],
  },
  {
    label: 'Settings',
    items: [
      { title: 'Profile', href: '/admin/profile', icon: UserIcon },
      { title: 'Billing', href: '/admin/billing', icon: CreditCardIcon },
      { title: 'Settings', href: '/admin/settings', icon: SettingsIcon },
    ],
  },
];

export const stats = [
  { label: 'Total Revenue', value: '$45,231', change: '+20.1%', up: true },
  { label: 'Active Users', value: '2,350', change: '+15.3%', up: true },
  { label: 'New Orders', value: '1,247', change: '-4.5%', up: false },
  { label: 'Conversion', value: '3.6%', change: '+1.2%', up: true },
];

export const transactions = [
  {
    name: 'Anna M. Hines',
    desc: 'Commissions',
    amount: '+$120.55',
    date: 'Wed Apr 24',
    status: 'Credit',
  },
  {
    name: 'Candice F. Gilmore',
    desc: 'Affiliates',
    amount: '+$9.68',
    date: 'Thu Dec 06',
    status: 'Credit',
  },
  {
    name: 'Vanessa R. Davis',
    desc: 'Grocery',
    amount: '-$105.22',
    date: 'Sat Apr 20',
    status: 'Debit',
  },
  {
    name: 'Judith H. Fritsche',
    desc: 'Refunds',
    amount: '+$80.59',
    date: 'Thu Apr 18',
    status: 'Credit',
  },
  {
    name: 'Peter T. Smith',
    desc: 'Bill Payments',
    amount: '-$750.95',
    date: 'Thu Apr 18',
    status: 'Debit',
  },
];

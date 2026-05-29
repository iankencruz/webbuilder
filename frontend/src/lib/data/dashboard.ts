export const navGroups = [
  {
    label: 'General',
    items: [
      { icon: '▦', label: 'Dashboard', href: '/admin/dashboard', active: true },
      { icon: '◈', label: 'Analytics', href: '/admin/analytics' },
      { icon: '◎', label: 'Reports', href: '/admin/reports' },
    ],
  },
  {
    label: 'Manage',
    items: [
      { icon: '⊞', label: 'Projects', href: '/admin/projects' },
      { icon: '◻', label: 'Tasks', href: '/admin/tasks' },
      { icon: '◷', label: 'Calendar', href: '/admin/calendar' },
      { icon: '⊜', label: 'Messages', href: '/admin/messages' },
    ],
  },
  {
    label: 'Settings',
    items: [
      { icon: '◉', label: 'Profile', href: '/admin/profile' },
      { icon: '⊟', label: 'Billing', href: '/admin/billing' },
      { icon: '◍', label: 'Settings', href: '/admin/settings' },
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

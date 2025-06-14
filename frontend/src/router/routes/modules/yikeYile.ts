import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const YIKEYILE: AppRouteRecordRaw = {
  path: '/yikeYile',
  name: 'yikeYile',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.yikeYile',
    requiresAuth: true,
    icon: 'icon-apps', // Placeholder icon, you might want to change this later
    order: 6, // Adjust order as needed to fit into the navigation bar
  },
  children: [
    {
      path: 'management',
      name: 'YikeYileManagement',
      component: () => import('@/views/yikeYile/index.vue'),
      meta: {
        locale: 'menu.yikeYile.management',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default YIKEYILE; 
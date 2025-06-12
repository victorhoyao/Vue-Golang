import { DEFAULT_LAYOUT } from'../base';
import { AppRouteRecordRaw } from'../types';

const APPLY: AppRouteRecordRaw = {
  path: '/apply',
  name: 'apply',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.consume',
    requiresAuth: true,
    icon: 'icon-alipay-circle',
    order: 2,
  },
  children: [
    {
      path: 'applylist',
      name: 'Applylist',
      component: () => import('@/views/apply/applylist/index.vue'),
      meta: {
        locale: 'menu.consume.consumelist',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
  ],
};

export default APPLY;

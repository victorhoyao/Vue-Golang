import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const TASKDISTRIBUTION: AppRouteRecordRaw = {
  path: '/taskDistribution',
  name: 'taskDistribution',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.taskDistribution',
    requiresAuth: true,
    icon: 'icon-share-alt',
    order: 8,
  },
  children: [
    {
      path: 'management',
      name: 'TaskDistributionManagement',
      component: () => import('@/views/taskDistribution/index.vue'),
      meta: {
        locale: 'menu.taskDistribution.management',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default TASKDISTRIBUTION; 
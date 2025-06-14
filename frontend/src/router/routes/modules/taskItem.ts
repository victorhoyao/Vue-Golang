import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const TASKITEM: AppRouteRecordRaw = {
  path: '/taskItem',
  name: 'taskItem',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.taskItem',
    requiresAuth: true,
    icon: 'icon-list', // You can change this icon as needed
    order: 7, // Adjust order as needed to fit into the navigation bar
  },
  children: [
    {
      path: 'management',
      name: 'TaskItemManagement',
      component: () => import('@/views/taskItem/index.vue'),
      meta: {
        locale: 'menu.taskItem.management',
        requiresAuth: true,
        roles: ['*'],
      },
    },
  ],
};

export default TASKITEM; 
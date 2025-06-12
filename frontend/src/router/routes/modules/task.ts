import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const TASK: AppRouteRecordRaw = {
  path: '/task',
  name: 'task',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.task',
    requiresAuth: true,
    icon: 'icon-ordered-list',
    order: 3,
  },
  children: [
    {
      path: 'tasklist',
      name: 'Tasklist',
      component: () => import('@/views/task/tasklist/index.vue'),
      meta: {
        locale: 'menu.task.tasklist',
        requiresAuth: true,
        roles: ['user'],
      },
    },
    {
      path: 'tasklog',
      name: 'Tasklog',
      component: () => import('@/views/task/tasklog/index.vue'),
      meta: {
        locale: 'menu.task.tasklog',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
    {
      path: 'taskadmin',
      name: 'Taskadmin',
      component: () => import('@/views/task/taskadmin/index.vue'),
      meta: {
        locale: 'menu.task.taskadmin',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
  ],
};

export default TASK;

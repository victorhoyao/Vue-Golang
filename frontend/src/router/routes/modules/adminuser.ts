import { DEFAULT_LAYOUT } from '../base';
import { AppRouteRecordRaw } from '../types';

const ADMINUSER: AppRouteRecordRaw = {
  path: '/adminuser',
  name: 'adminuser',
  component: DEFAULT_LAYOUT,
  meta: {
    locale: 'menu.adminuser',
    requiresAuth: true,
    icon: 'icon-user-group',
    order: 1,
  },
  children: [
    {
      path: 'userlist',
      name: 'Userlist',
      component: () => import('@/views/adminuser/userlist/index.vue'),
      meta: {
        locale: 'menu.adminuser.userlist',
        requiresAuth: true,
        roles: ['admin'],
      },
    },
    // {
    //   path: 'userworklist',
    //   name: 'UserWorkList',
    //   component: () => import('@/views/adminuser/userworklist/index.vue'),
    //   meta: {
    //     locale: 'menu.adminuser.UserWorkList',
    //     requiresAuth: true,
    //     roles: ['admin'],
    //   },
    // },
  ],
};

export default ADMINUSER;

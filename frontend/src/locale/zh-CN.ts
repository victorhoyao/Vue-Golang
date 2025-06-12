import localeMessageBox from '@/components/message-box/locale/zh-CN';
import localeLogin from '@/views/login/locale/zh-CN';
import localeWorkplace from '@/views/dashboard/workplace/locale/zh-CN';
import localeUserlist from '@/views/adminuser/userlist/locale/zh-CN';
import localeRegister from '@/views/register/locale/zh-CN';
import localeSettings from './zh-CN/settings';

export default {
  'menu.dashboard': '首页',
  'menu.server.dashboard': '仪表盘-服务端',
  'menu.server.workplace': '工作台-服务端',
  'menu.server.monitor': '实时监控-服务端',
  'menu.adminuser': '用户管理',
  'menu.adminuser.userlist': '用户列表',
  'menu.adminuser.UserWorkList':"用户统计列表",
  'menu.consume': '提现管理',
  'menu.consume.consumelist': '提现列表',
  'menu.task': '任务管理',

  'menu.task.tasklist': '任务列表',
  'menu.task.tasklog': '用户领取列表',
  'menu.task.taskadmin': '订单管理列表',
  'menu.user': '个人中心',
  'menu.user.info': '个人信息',


  'menu.list': '列表页',
  'menu.result': '结果页',
  'menu.exception': '异常页',
  'menu.form': '表单页',
  'menu.profile': '详情页',
  'menu.visualization': '数据可视化',
  'menu.arcoWebsite': 'Arco Design',
  'menu.faq': '常见问题',
  'navbar.docs': '文档中心',
  'navbar.action.locale': '切换为中文',
  ...localeSettings,
  ...localeMessageBox,
  ...localeLogin,
  ...localeWorkplace,
  ...localeUserlist,
  ...localeRegister
};

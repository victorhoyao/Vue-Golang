import axios from 'axios';
import { HttpResponse } from './interceptor';

// 用户列表记录接口定义
export interface UserListRecord {
  id: number;                // 用户ID
  userName: string;          // 用户名
  status: number;            // 用户状态：0-正常，1-关注
  authority: number;         // 用户权限级别
  addUserId: number;         // 添加该用户的管理员ID
  userKey: string;           // 用户唯一标识符
  money: number;             // 用户余额
  tranType: string;          // 交易类型
  tranAccount: string;       // 交易账号
  tranName: string;          // 交易名称
  createTime: string;        // 创建时间
}

// 分页参数接口定义
export interface Pagination {
  pageNum: number;           // 当前页码
  pageSize: number;          // 每页显示条数
  total?: number;            // 总记录数
  current?:number,           // 当前页（兼容性参数）
  countDate?:string          // 统计日期
}

// 添加用户数据接口定义
export interface addData {
  userName: string;          // 用户名
  passWord: string;          // 密码
}

/**
 * 获取用户列表
 * @param params 分页参数
 * @returns 用户列表数据
 */
export function getUserList(params: Pagination) {
  return axios.get('/User/userList', {
    params,
  });
}

/**
 * 添加新用户
 * @param data 用户数据（用户名和密码）
 * @returns 添加结果
 */
export function addUser(data: addData): Promise<HttpResponse> {
  return axios.post<HttpResponse>('/User/addUser', data);
}

/**
 * 修改用户密码
 * @param data 包含userId和newPassword的对象
 * @returns 修改结果
 */
export function editPass(data: any): Promise<HttpResponse> {
  return axios.post<HttpResponse>('/User/editPass', data);
}

/**
 * 删除用户
 * @param params 包含userId的对象
 * @returns 删除结果
 */
export function delUser(params:any): Promise<HttpResponse> {
  return axios.get<HttpResponse>('/User/delUser', {params});
}

/**
 * 获取用户工作列表
 * @param params 查询参数
 * @returns 用户工作列表数据
 */
export function getUserWorkList(params:any) {
  return axios.get<HttpResponse>('/userWork/getUserWorkList', {params});
}

/**
 * 根据关键词搜索用户
 * @param keyWord 搜索关键词（用户名）
 * @returns 搜索结果
 */
export function findUser(keyWord: string): Promise<HttpResponse> {
  return axios.get<HttpResponse>('/User/findUser', {
    params: { keyWord }
  });
}

// export function logout() {
//   return axios.post('/api/user/logout');
// }

// export function getUserInfo() {
//   return axios.get('/User/getMyInfo');
// }

// export function getMenuList() {
//   return axios.post<RouteRecordNormalized[]>('/api/user/menu');
// }

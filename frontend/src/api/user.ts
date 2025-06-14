import axios from 'axios';
import type { RouteRecordNormalized } from 'vue-router';
// import { UserState } from '@/store/modules/user/types';

export interface LoginData {
  userName: string;
  passWord: string;
}

export interface RegisterData {
  userName: string;
  passWord: string;
  rootPassWord: string;
}

export interface LoginRes {
  token: string;
  id: number;
  authority: number;
  ukey: string;
}

export interface ChangePassData {
  oldPass: string;
  newPass: string;
}

export function login(data: LoginData) {
  return axios.post<LoginRes>('/login', data);
}

export function register(data: RegisterData) {
  return axios.post<LoginRes>('/register', data);
}

export function changePassword(data: ChangePassData) {
  return axios.post<any>('/User/changePassword', data);
}

// export function logout() {
//   return axios.post('/api/user/logout');
// }

export function getUserInfo() {
  return axios.get('/User/getMyInfo');
}

export function getMenuList() {
  return axios.post<RouteRecordNormalized[]>('/api/user/menu');
}

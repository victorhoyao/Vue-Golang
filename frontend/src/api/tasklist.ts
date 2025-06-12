import axios from 'axios';
import { HttpResponse } from './interceptor';

export interface Pagination {
  pageNum: number;
  pageSize: number;
  total?: number;
  current?: number;
}

export function getMyTaskLogList(params: Pagination) {
  return axios.get('/taskLog/getMyTaskLogList', {
    params,
  });
}

export function getTaskList(params: Pagination) {
  return axios.get('/taskList/getTaskList', {
    params,
  });
}

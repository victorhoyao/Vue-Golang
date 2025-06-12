import axios from 'axios';
import { HttpResponse } from './interceptor';

export interface Pagination {
  pageNum: number;
  pageSize: number;
  total?: number;
  current?: number;
}

export interface TaskLogParams extends Pagination {
  ksAccount?: string;
  preDate?: string;
  nextDate?: string;
}

export function getTaskLogList(params: Pagination) {
  return axios.get('/taskLog/getTaskLogList', {
    params,
  });
}


export function getTaskLogListByKey(params: TaskLogParams) {
  return axios.get('/taskLog/getTaskLogListByKey', {
    params,
  });
}

import axios from 'axios';

export interface Pagination {
  pageNum: number;
  pageSize: number;
  total?: number;
  current?: number;
}

// 添加搜索任务列表的接口
export function getTaskListByKey(params: {
  pageNum: number;
  pageSize: number;
  orderId?: string;
}) {
  return axios.get('/taskList/getTaskListByKey', {
    params,
  });
}

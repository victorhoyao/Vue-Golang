/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-26 13:12:10
 * @FilePath: \ksDiggPlatform\src\api\Price.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from 'axios';

export interface Pagination {
  pageNum: number;
  pageSize: number;
  total?: number;
  current?: number;
}

export function getPrice() {
  return axios.get('/Manager/getPrice');
}

export function setPrice(data:any) {
  return axios.post('/Manager/setPrice',data);
}



export function setTcGl(data:any) {
  return axios.post('/Manager/setTcGl',data);
}

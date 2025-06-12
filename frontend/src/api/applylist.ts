/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-16 15:07:54
 * @FilePath: \ksDiggPlatform\src\api\applylist.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from 'axios';

export interface Pagination {
  pageNum: number;
  pageSize: number;
  total?: number;
  current?: number;
}

export function getApplyList(params: Pagination) {
  return axios.get('/Tran/getApplyList', {
    params,
  });
}

export function SetUpWithdrawalAccount(data:any) {
  return axios.post('/User/setTran', data);
}

export function SubmitWithdrawalApplication(data:any) {
  return axios.post('/Tran/applyTransaction', data);
}


export function getMyTranList(params:Pagination) {
  return axios.get('/Tran/getMyTranList', {
    params
  });
}

export function SubmitDoneTransaction(id:any,data:any) {
  return axios.post(`/Tran/doneTransaction/${id}`, data);
}
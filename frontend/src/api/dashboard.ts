/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-26 12:50:00
 * @FilePath: \ksDiggPlatform\src\api\dashboard.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from 'axios';

export interface PopularRecord {
  key: number;
  clickNumber: string;
  title: string;
  increases: number;
}

export function getTask(params: any) {
  return axios.get('/Task/getTask', {
    params,
  });
}

export function parseSecUidDouYin(params: any) {
  return axios.get('NoAuth/getDyInfo' , {
    params,
  });
}


export function getTaskLogCount(params: any) {
  return axios.get('/taskLog/getTaskLogCount' , {
    params,
  });
}
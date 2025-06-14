/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-16 00:58:44
 * @FilePath: \ksDiggPlatform\src\api\interceptor.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from 'axios';
import type { AxiosRequestConfig, AxiosResponse } from 'axios';
import { Message } from '@arco-design/web-vue';
// import { useUserStore } from '@/store';
import { getToken } from '@/utils/auth';

export interface HttpResponse<T = unknown> {
  msg?: string;
  code?: number;
  data: T;
}

if (import.meta.env.VITE_API_BASE_URL) {
  axios.defaults.baseURL = import.meta.env.VITE_API_BASE_URL;
}

axios.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // eslint-disable-next-line no-console
    console.log('Making request to:', config.url, 'with data:', config.data);
    const token = getToken();
    // eslint-disable-next-line no-console
    console.log('Token from storage:', token);
    if (token) {
      if (!config.headers) {
        config.headers = {};
      }
      config.headers.token = `${token}`;
      // eslint-disable-next-line no-console
      console.log('Added token to headers:', config.headers.token);
    } else {
      // eslint-disable-next-line no-console
      console.log('No token found in storage');
    }
    return config;
  },
  (error) => {
    // eslint-disable-next-line no-console
    console.error('Request error:', error);
    return Promise.reject(error);
  }
);
// add response interceptors
axios.interceptors.response.use(
  (response: AxiosResponse<HttpResponse>) => {
    // eslint-disable-next-line no-console
    console.log('Response received:', response.data);
    const res = response.data;

    // if the custom code is not 20000, it is judged as an error.
    if (res.code !== 200) {
      // eslint-disable-next-line no-console
      console.error('Response error:', res);
      Message.error({
        content: res.msg || 'Error',
        duration: 5 * 1000,
      });

      return Promise.reject(new Error(res.msg || '错误'));
    }
    return res;
  },
  (error) => {
    // eslint-disable-next-line no-console
    console.error('Response error:', error.response || error);
    Message.error({
      content: error.response?.data?.msg || error.message || 'Request 错误',
      duration: 5 * 1000,
    });
    return Promise.reject(error);
  }
);

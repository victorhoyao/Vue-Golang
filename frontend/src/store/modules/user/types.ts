/*
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-12 12:36:11
 * @FilePath: \ksDiggPlatform\src\store\modules\user\types.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export type RoleType = '' | '*' | 'admin' | 'user';
export interface UserState {
  userName?: string;
  id?:  number;
  status?: number;
  authority?: number;
  addUserId?: number;
  userKey?: string;
  role: RoleType;
  money?: number;

}

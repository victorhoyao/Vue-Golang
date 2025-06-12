<template>
  <div class="container">
    <a-row class="wrapper">
      <a-col :span="24">
        <a-tabs default-active-key="1" type="rounded" @tab-click="tabClick">
          <a-tab-pane key="1" title="提现">

            <div style="margin-top: 20px">
              <a-form ref="formRef" :model="UserInfo" class="form" :label-col-props="{ span: 8 }"
                :wrapper-col-props="{ span: 16 }">
                <a-divider orientation="left">提现方式</a-divider>

                <a-form-item field="money" label="提现方式">
                  <a-select placeholder="请选择提现方式" v-model="UserInfo.tranType" :disabled="UserInfo.isTran">
                    <a-option>U提现</a-option>
                  </a-select> </a-form-item>
                <a-form-item field="money" label="提现账号">
                  <a-input placeholder="请输入" :disabled="UserInfo.isTran" v-model="UserInfo.account" />
                </a-form-item>
                <a-form-item field="money" label="姓名">
                  <a-input placeholder="请输入 随便填" :disabled="UserInfo.isTran" v-model="UserInfo.name" />
                </a-form-item>

                <a-form-item>
                  <a-button type="primary" long @click="setTran" :disabled="UserInfo.isTran">设置</a-button>

                </a-form-item>

                <a-divider orientation="left">提现积分</a-divider>


                <a-form-item field="money" label="当前积分">
                  <a-input-number placeholder="请输入" disabled v-model="UserInfo.money" />
                </a-form-item>
                <a-form-item field="applyprice" label="提现积分">
                  <a-input-number placeholder="请输入" v-model="UserInfo.applyprice" />
                </a-form-item>
                <a-form-item>
                  <a-button type="primary" status="success" long @click="WithdrawalBalance">提现</a-button>
                </a-form-item>
              </a-form>
            </div>

          </a-tab-pane>


          <a-tab-pane key="3" title="提现记录">
            <a-row>
              <a-col :span="24">
                <a-table :data="MyTranList" :pagination="false" :columns="MyTranListcolumns" row-key="id" :bordered="{cell:true}">
                  <template #status="{ record }">
                    <a-tag v-if="record.status === 0" color="blue" bordered>
                      {{ `审核中` }}</a-tag>
                    <a-tag v-if="record.status === 1" color="green" bordered>
                      {{ `审核中` }}</a-tag>
                    <a-tag v-if="record.status === 2" bordered>
                      {{ `审核中` }}</a-tag>
                  </template></a-table>
              </a-col>
            </a-row>
          </a-tab-pane>
        </a-tabs>
      </a-col>
    </a-row>
  </div>
</template>

<script lang="ts" setup>
import {
  SetUpWithdrawalAccount,
  SubmitWithdrawalApplication,
  getMyTranList,
} from '@/api/applylist';
import { getUserInfo } from '@/api/user';
import { useUserStore } from '@/store';
import { Message } from '@arco-design/web-vue';
import { ref, onMounted } from 'vue';

const UserInfo = ref({
  money: 0,
  applyprice: 0.01,
  tranType: "",
  account: "",
  name: "",
  isTran: false,

});

onMounted(async () => {
  const userStore:any = useUserStore();
  console.log(userStore.$state);
  if (userStore.$state.tranType !== "") {
    UserInfo.value.money = userStore.$state.money || 0
    UserInfo.value.tranType = userStore.$state.tranType || ""
    UserInfo.value.account = userStore.$state.tranAccount || ""
    UserInfo.value.name = userStore.$state.tranName || ""
    UserInfo.value.isTran = true
  }



});
const setTran = async () => {

  const Withdrawal: any = {
    account: UserInfo.value.account,
    accountType: UserInfo.value.tranType,
    name: UserInfo.value.name,
  }

  const { msg, code }:any= await SetUpWithdrawalAccount(Withdrawal);
  if (code === 200) {
    Message.success(msg);
  } else {
    Message.error(msg);
  }
  setTimeout(() => {
    // 执行刷新屏幕的操作，例如 location.reload()
    window.location.reload();
  }, 1500);


};



const WithdrawalBalance = async () => {
  console.log(UserInfo.value);
  const applyprice = {
    applyprice: UserInfo.value.applyprice
  }
  const data:any = await SubmitWithdrawalApplication(
    applyprice
  );

  if (data.code === 200) {
    Message.success({
      content: data.msg,
      duration: 3000
    });
  }

  setTimeout(() => {
    // 执行刷新屏幕的操作，例如 location.reload()
    window.location.reload();
  }, 1500);
};

const tabClick = (key: number | string) => {
  console.log(key);
  if (key === '3') {
    getMyTran();
  }
};
const basePagination = { pageNum: 1, pageSize: 999 };

const MyTranList = ref([]);
const MyTranListcolumns = [
  {
    title: '提现订单号',
    dataIndex: 'orderId',
  },
  {
    title: '提现申请时间',
    dataIndex: 'createTime',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '提现积分',
    dataIndex: 'applyprice',
  },
  {
    title: '收款app',
    dataIndex: 'tranType',
  },
  {
    title: '审核完成时间',
    dataIndex: 'doneTime',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '收款账号',
    dataIndex: 'tranAccount',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
  },
  {
    title: '审核结果说明',
    dataIndex: 'desc',
  },
];
const getMyTran = async () => {
  const { data } = await getMyTranList(basePagination);

  MyTranList.value = data.result;
};
</script>

<script lang="ts">
export default {
  name: 'Setting',
};
</script>

<style scoped lang="less">
.form {
  width: 540px;
  margin: 0 auto;
}

.container {
  padding: 0 20px 20px 20px;
}

.wrapper {
  padding: 20px 0 0 20px;
  min-height: 580px;
  background-color: var(--color-bg-2);
  border-radius: 4px;
}

:deep(.section-title) {
  margin-top: 0;
  margin-bottom: 16px;
  font-size: 14px;
}
</style>

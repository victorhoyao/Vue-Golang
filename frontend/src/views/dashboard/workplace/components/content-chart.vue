<!--
 * @Author: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @Date: 2024-08-01 23:46:37
 * @LastEditors: error: error: git config user.name & please set dead value or install git && error: git config user.email & please set dead value or install git & please set dead value or install git
 * @LastEditTime: 2024-08-26 11:45:53
 * @FilePath: \ksDiggPlatform\src\views\dashboard\workplace\components\content-chart.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <a-spin style="width: 100%">
    <a-card
      class="general-card"
      :header-style="{ paddingBottom: 0 }"
      :body-style="{
        paddingTop: '20px',
      }"
      title="价格设置"
    >
      <a-space>
        <a-select
          :style="{ width: '320px' }"
          v-model="value"
          placeholder="请选择"
          allow-clear
        >
          <a-option disabled>KS--------</a-option>

          <a-option >快手粉</a-option>
          <a-option >快手赞</a-option>

        </a-select>
        <a-input-number
          :style="{ width: '320px' }"
          placeholder="请输入"
          allow-clear
          v-model="Price"
          :precision="3"
        >
          <template #prepend> 价格 </template>
        </a-input-number>
        <a-button type="primary" @click="subPriceSetting">提交</a-button>
      </a-space>
    </a-card>
  </a-spin>
</template>

<script lang="ts" setup>
  import { setPrice } from '@/api/Price';
  import { ref } from 'vue';
  import { Message } from '@arco-design/web-vue';
import { useRouter } from 'vue-router';

  const value = ref('快手粉');
  const Price = ref(0.001);
  const router = useRouter();

  const subPriceSetting = async () => {
    const res = {
      type: value.value,
      price: Price.value,
    };

    const data = await setPrice(res);

    Message.success(data.msg);
    router.go(0);

  };

</script>

<style scoped lang="less"></style>

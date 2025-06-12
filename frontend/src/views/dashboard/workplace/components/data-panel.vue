<template>
  <a-grid :cols="24" :row-gap="16" class="panel">
    <a-grid-item class="panel-col" :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }">
      <a-space>
        <a-avatar :size="54" class="col-avatar">
          <img alt="avatar"
            src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/288b89194e657603ff40db39e8072640.svg~tplv-49unhts6dw-image.image" />
        </a-avatar>
        <a-statistic title="ks点赞单价" :value="Priceinfo.KSDiggPrice" :precision="3" :value-from="0" animation
          show-group-separator>
          <template #suffix>
            <!-- <span class="unit">{{ $t('workplace.pecs') }}</span> -->
          </template>
        </a-statistic>
      </a-space>
    </a-grid-item>
    <a-grid-item class="panel-col" :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }">
      <a-space>
        <a-avatar :size="54" class="col-avatar">
          <img alt="avatar"
            src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/288b89194e657603ff40db39e8072640.svg~tplv-49unhts6dw-image.image" />
        </a-avatar>
        <a-statistic title="ks关注单价" :value="Priceinfo.KSfenPrice" :precision="3" :value-from="0" animation
          show-group-separator>
          <template #suffix>
            <!-- <span class="unit">{{ $t('workplace.pecs') }}</span> -->
          </template>
        </a-statistic>
      </a-space>
    </a-grid-item>
    <a-grid-item class="panel-col" :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }">
      <a-space>
        <a-avatar :size="54" class="col-avatar">
          <img alt="avatar"
            src="//p3-armor.byteimg.com/tos-cn-i-49unhts6dw/77d74c9a245adeae1ec7fb5d4539738d.svg~tplv-49unhts6dw-image.image" />
        </a-avatar>
        <a-statistic title="账户可用" :value="userStore.$state.money" :value-from="0" :precision="3" animation
          show-group-separator>
          <template #suffix>
            <span class="unit">{{ '积分' }}</span>
            <!-- <icon-caret-up class="up-icon" /> -->
          </template>
        </a-statistic>
      </a-space>
    </a-grid-item>
    <a-grid-item v-permission="['admin']" class="panel-col" :span="{ xs: 12, sm: 12, md: 12, lg: 12, xl: 12, xxl: 6 }">
      <a-space>
        <a-statistic title="审核失败概率" :value="Priceinfo.shTcGl" animation show-group-separator>
          <template #suffix>
            <span class="unit">{{ '%' }}</span>
            <a-slider :default-value="50" :style="{ width: '200px' }" v-model="Priceinfo.shTcGl" :min="0" :max="50"
              @change="onChangeshTcGl"  :show-ticks="true"/>



          </template>
        </a-statistic>
      </a-space>
    </a-grid-item>
    <a-grid-item :span="24">
      <a-divider class="panel-border" />
    </a-grid-item>
  </a-grid>

</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { useUserStore } from '@/store';
import { getPrice, setTcGl } from '@/api/Price';
import { Message } from '@arco-design/web-vue';

const Priceinfo: any = ref({
  KSDiggPrice: 0,
  KSfenPrice: 0,
  shTcGl: 0,
  updateTime: 0,
});
const userStore = useUserStore();


onMounted(async () => {
  const { data } = await getPrice();

  Priceinfo.value = data.result;
  console.log(userStore.$state);
});

const onChangeshTcGl = async (value: number) => {
  console.log(value);
  const { data: { code, msg } } = await setTcGl({ shTcGl: value });
  if (code === 200) {
    Message.success(msg);
  } else {
    Message.error(msg);
  }
}
</script>

<style lang="less" scoped>
.arco-grid.panel {
  margin-bottom: 0;
  padding: 16px 20px 0 20px;
}

.panel-col {
  padding-left: 43px;
  border-right: 1px solid rgb(var(--gray-2));
}

.col-avatar {
  margin-right: 12px;
  background-color: var(--color-fill-2);
}

.up-icon {
  color: rgb(var(--red-6));
}

.unit {
  margin-left: 8px;
  color: rgb(var(--gray-8));
  font-size: 12px;
}

:deep(.panel-border) {
  margin: 4px 0 0 0;
}
</style>

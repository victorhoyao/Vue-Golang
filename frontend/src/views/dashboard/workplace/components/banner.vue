<template>
  <a-col class="banner">

    <a-col  :span="{ xs: 24, sm: 24, md: 24, lg: 12, xl: 12, xxl: 12 }">
      <a-typography-title :heading="5" style="margin-top: 0">
        {{ '欢迎回来！'+userStore.$state.userName }}
      </a-typography-title>
    </a-col>
    <a-divider class="panel-border" />
  </a-col>
</template>

<script lang="ts" setup>
  import { useUserStore } from '@/store';

  const userStore = useUserStore();

</script>

<style scoped lang="less">
  .banner {
    width: 100%;
    padding: 20px 20px 0 20px;
    background-color: var(--color-bg-2);
    border-radius: 4px 4px 0 0;
  }

  :deep(.arco-icon-home) {
    margin-right: 6px;
  }
</style>

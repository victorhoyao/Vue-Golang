<template>
  <div class="container">
      <div class="left-panel">
          <div class="content-panel">
              <Banner />
              <DataPanel />
              <ContentChart v-permission="['admin']" />
          </div>
          <a-grid :cols="24" :col-gap="16" :row-gap="16" class="grid-margin">
              <a-grid-item :span="{ xs: 24, sm: 24, md: 24, lg: 12, xl: 12, xxl: 12 }"></a-grid-item>
              <a-grid-item :span="{ xs: 24, sm: 24, md: 24, lg: 12, xl: 12, xxl: 12 }"></a-grid-item>
          </a-grid>
          <RecentlyVisited v-permission="['user']" />
      </div>
  </div>
</template>

<script lang="ts" setup>
import Banner from './components/banner.vue';
import DataPanel from './components/data-panel.vue';
import ContentChart from './components/content-chart.vue';
import CategoriesPercent from './components/categories-percent.vue';
import RecentlyVisited from './components/recently-visited.vue';
import QuickOperation from './components/quick-operation.vue';
import Announcement from './components/announcement.vue';
import Carousel from './components/carousel.vue';
import Docs from './components/docs.vue';
</script>

<script lang="ts">
export default {
  name: 'Dashboard', // 要使 keep-alive 的 include 属性生效，必须命名组件
};
</script>

<style lang="less" scoped>
.container {
  background-color: var(--color-fill-2);
  padding: 16px 20px 0;
  display: flex;
}

.left-panel {
  flex: 1;
  overflow: auto;
}

.content-panel {
  background-color: var(--color-bg-2);
  border-radius: 4px;
  overflow: auto;
}

:deep(.panel-border) {
  margin-bottom: 0;
  border-bottom: 1px solid rgb(var(--gray-2));
}

.moduler-wrap {
  border-radius: 4px;
  background-color: var(--color-bg-2);

  :deep(.text) {
      font-size: 12px;
      text-align: center;
      color: rgb(var(--gray-8));
  }

  :deep(.wrapper) {
      margin-bottom: 8px;
      text-align: center;
      cursor: pointer;

      &:last-child .text {
          margin-bottom: 0;
      }

      &:hover {

          .icon,
          .text {
              color: rgb(var(--arcoblue-6));
              background-color: #e8f3ff;
          }
      }
  }

  :deep(.icon) {
      display: inline-block;
      width: 32px;
      height: 32px;
      margin-bottom: 4px;
      color: rgb(var(--dark-gray-1));
      line-height: 32px;
      font-size: 16px;
      text-align: center;
      background-color: rgb(var(--gray-1));
      border-radius: 4px;
  }
}

.grid-margin {
  margin-top: 16px;
}
</style>

<style lang="less" scoped>
// 响应式
.mobile {
  .container {
      display: block;
  }

  .right-side {
      width: 100%;
      margin: 16px 0 0 0;
  }
}
</style>
<template>
  <a-card class="general-card" title="数据统计" :header-style="{ paddingBottom: '0' }" :body-style="{ paddingTop: '26px' }">
    <div style="margin-bottom: -1rem">
      <a-row :gutter="8">
        <a-col :span="6">
          <a-date-picker style="width: 300px;" v-model="selectedDate" />
        </a-col>
        <a-col :span="6">
          <a-select placeholder="类型" allow-clear v-model="selectedTaskType">
            <a-option disabled>KS--------</a-option>
            <a-option :value="11">快手粉</a-option>
            <a-option :value="12">快手赞</a-option>
          </a-select>
        </a-col>
        <a-col :span="2">
          <a-button type="primary" @click="fetchTaskStatistics">
            <template #icon>
              <icon-search />
            </template>
            <template #default>搜索</template>
          </a-button>
        </a-col>
      </a-row>
      <br>
      <a-row :gutter="8">
        <a-col :span="24">
          <a-steps type="arrow" :current="2">
            <a-step description="审核中">
              <a-link>{{ taskStatistics.waitCount }}</a-link>
            </a-step>
            <a-step description="审核成功">
              <a-link status="success">{{ taskStatistics.sucessCount }}</a-link>
            </a-step>
            <a-step description="审核失败">
              <a-link status="danger">{{ taskStatistics.failCount }}</a-link>
            </a-step>
          </a-steps>
        </a-col>
      </a-row>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { getTaskLogCount } from '@/api/dashboard';
import Message from "@arco-design/web-vue/es/message";

const selectedTaskType = ref(11) // 任务类型
const selectedDate = ref() // 选择的日期
const taskStatistics = ref({
  failCount: 0,
  sucessCount: 0,
  waitCount: 0
})

// 获取任务统计信息
const fetchTaskStatistics = async () => {
  if (!selectedDate.value) {
    Message.error("请填写时间");
    return;
  }

  const requestPayload = {
    date: selectedDate.value,
    type: selectedTaskType.value
  }

  try {
    const response = await getTaskLogCount(requestPayload);
    taskStatistics.value = response.data;
    Message.success(response.msg);
  } catch (error) {
    console.error(error);
    Message.error("获取数据失败");
  }
}
</script>

<style lang="less" scoped>
:deep(.arco-card-header-title) {
  line-height: inherit;
}
</style>
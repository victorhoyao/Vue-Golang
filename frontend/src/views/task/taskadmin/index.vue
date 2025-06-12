<template>
  <div class="container">
    <a-card class="general-card" title="查询任务">
      <a-row class="search-form">
        <a-col :flex="'auto'">
          <a-form :model="searchForm" layout="inline">
            <a-form-item field="orderId" label="订单编号">
              <a-input v-model="searchForm.orderId" placeholder="请输入订单编号" allow-clear />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button type="primary" @click="search">
                  搜索
                </a-button>
                <a-button @click="reset">
                  重置
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-col>
      </a-row>
      <a-divider style="margin-top: 0" />
      <a-table row-key="id" :loading="loading" :columns="columns" :data="renderData" :bordered="{ cell: true }"
        :pagination="pagination" @page-change="onPageChange">

        <template #UStatus="{ record }">
          <a-tag v-if="record.UStatus === 0" color="blue" bordered>
            {{ `处理中` }}</a-tag>
          <a-tag v-if="record.UStatus === 1" color="arcoblue" bordered>
            {{ `已完成` }}</a-tag>
          <a-tag v-if="record.UStatus === 2" color="gold" bordered>
            {{ `退单` }}</a-tag>
        </template>
        <template #goodsStatus="{ record }">

          <a-tag v-if="record.goodsStatus === 1" color="arcoblue" bordered>
            {{ `已付款` }}</a-tag>
          <a-tag v-if="record.goodsStatus === 3" color="gold" bordered>
            {{ `处理中` }}</a-tag>
          <a-tag v-if="record.goodsStatus === 6" color="blue" bordered>
            {{ `已完成` }}</a-tag>
          <a-tag v-if="record.goodsStatus === 7" color="blue" bordered>
            {{ `已退单` }}</a-tag>
          <a-tag v-if="record.goodsStatus === 8" color="blue" bordered>
            {{ `已退款` }}</a-tag>
        </template>

        <template #sumbmitTime="{ record }">
          <a-typography-text type="secondary" v-if="record.sumbmitTime == ''">
            {{ `未提交` }}
          </a-typography-text>
          <a-typography-text type="secondary" v-if="record.sumbmitTime">
            {{ record.sumbmitTime }}
          </a-typography-text>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import type { Pagination } from '../../../api/applylist';
import { getTaskList } from '../../../api/tasklist';
import { getTaskListByKey } from '../../../api/taskadmin';
import useLoading from '../../../hooks/loading';

const { loading, setLoading } = useLoading();

const searchForm = reactive({
  orderId: '',
});

type Column = TableColumnData & { checked?: true };

const columns = ref<Column[]>([
  {
    title: '编号',
    dataIndex: 'id',
    slotName: 'id',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '上级订单编号',
    dataIndex: 'orderId',
    slotName: 'orderId',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '上级_购买数量',
    dataIndex: 'buyNumber',
    slotName: 'buyNumber',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '商品',
    dataIndex: 'goodsName',
    slotName: 'goodsName',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '领取数量',
    dataIndex: 'collectNum',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '上级_当前进度',
    dataIndex: 'currentNum',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '订单拉取时间',
    dataIndex: 'downTime',
    ellipsis: true,
    tooltip: true,
    slotName: 'downTime',

  },
  {
    title: '订单完成时间',
    dataIndex: 'doneTime',
    ellipsis: true,
    tooltip: true,
    slotName: 'doneTime',

  },


  {
    title: '订单状态',
    dataIndex: 'UStatus',
    slotName: 'UStatus',
  },
  {
    title: '上级_商品状态',
    dataIndex: 'goodsStatus',
    slotName: 'goodsStatus',

  },
  {
    title: '视频短链接',
    dataIndex: 'shortLink',
    slotName: 'shortLink',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '最后一次领取时间',
    dataIndex: 'lastGetTime',
  },


]);
const renderData = ref([]);

const basePagination: Pagination = { pageNum: 1, pageSize: 10, current: 1 };
const pagination = reactive({
  ...basePagination,
});

const fetchData = async (
  params: Pagination = { pageNum: 1, pageSize: 10 },
  searchParams = searchForm
) => {
  setLoading(true);
  let data;

  try {

    if (searchParams.orderId) {
      const response = await getTaskListByKey({
        pageNum: params.pageNum,
        pageSize: params.pageSize,
        orderId: searchParams.orderId,
      });
      data = response.data;
    } else {
      const response = await getTaskList(params);
      data = response.data;
    }

    console.log('任务列表', data);
    renderData.value = data.result;
    pagination.current = params.pageNum;
    pagination.total = data.total;
  } catch (err) {
    // you can report use errorHandler or other
  } finally {
    setLoading(false);
  }
};
fetchData();

const onPageChange = (current: number) => {
  pagination.pageNum = current;
  fetchData(pagination);
};

const search = () => {
  pagination.pageNum = 1;
  fetchData(pagination, searchForm);
};

const reset = () => {
  searchForm.orderId = '';
  pagination.pageNum = 1;
  fetchData(pagination);
};
</script>
<style scoped lang="less">
.container {
  padding: 0 20px 20px 20px;
}

.search-form {
  margin-bottom: 16px;
}

:deep(.arco-table-th) {
  &:last-child {
    .arco-table-th-item-title {
      margin-left: 16px;
    }
  }
}

.action-icon {
  margin-left: 12px;
  cursor: pointer;
}

.active {
  color: #0960bd;
  background-color: #e3f4fc;
}

.setting {
  display: flex;
  align-items: center;
  width: 200px;

  .title {
    margin-left: 12px;
    cursor: pointer;
  }
}
</style>

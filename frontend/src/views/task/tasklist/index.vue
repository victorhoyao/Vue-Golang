<template>
  <div class="container">
    <a-card class="general-card" title="查询任务">
      <a-divider style="margin-top: 0" />
      <a-table row-key="id" :loading="loading" :columns="columns" :data="renderData" :bordered="{cell:true}"
        :pagination="pagination" @page-change="onPageChange">
        <!-- <template #goodsId="{ record }">
          <a-tag v-if="record.goodsId === 53" color="#165dff">
            {{ `ks粉` }}</a-tag>
          <a-tag v-if="record.goodsId === 55" color="#00b42a">
            {{ `ks藏` }}</a-tag>
          <a-tag v-if="record.goodsId === 54" color="#f53f3f">
            {{ `KS赞` }}</a-tag>
        </template> -->
        <template #status="{ record }">
          <a-tag v-if="record.status === 0" color="blue" bordered>
            {{ `未提交` }}</a-tag>
          <a-tag v-if="record.status === 1" color="arcoblue" bordered>
            {{ `等待审核中` }}</a-tag>
          <a-tag v-if="record.status === 2" color="gold" bordered>
            {{ `放弃` }}</a-tag>
          <a-tag v-if="record.status === 3" color="green" bordered>
            {{ `完成` }}</a-tag>
          <a-tag v-if="record.status === 4" color="#f53f3f" bordered>
            {{ `放弃` }}</a-tag>
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
import { getMyTaskLogList } from '../../../api/tasklist';
import useLoading from '../../../hooks/loading';

const { loading, setLoading } = useLoading();

type Column = TableColumnData & { checked?: true };

const columns = ref<Column[]>([
  {
    title: '订单编号',
    dataIndex: 'orderNum',
    slotName: 'orderNum',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '领取账号',
    dataIndex: 'Account',
    ellipsis: true,
    tooltip: true,
  },
  {
      title: '商品',
      dataIndex: 'goodsName',
      ellipsis: true,
      tooltip: true,
      slotName: 'goodsName',
  
    },
  {
    title: 'url',
    dataIndex: 'videoLink',
    ellipsis: true,
    tooltip: true,
  },

  // {
  //   title: '类型',
  //   dataIndex: 'goodsId',
  //   ellipsis: true,
  //   tooltip: true,
  //   slotName: 'goodsId',

  // },
  {
    title: '单价',
    dataIndex: 'price',
    ellipsis: true,
    tooltip: true,

  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
  },
  {
    title: '接单时间',
    dataIndex: 'collectTime',
  },
  {
    title: '提交时间',
    dataIndex: 'sumbmitTime',
    slotName: 'sumbmitTime',
  },
  {
    title: '审核时间',
    dataIndex: 'examineTime',
  },
  // {
  //   title: '领取用户id',
  //   dataIndex: 'userId',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  // {
  //   title: '领取任务secuid',
  //   dataIndex: 'secUid',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  // {
  //   title: '任务平台名称',
  //   dataIndex: 'pingtaiName',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  // {
  //   title: ' 平台',
  //   dataIndex: 'app',
  //   ellipsis: true,
  //     tooltip: true,
  // },
  // {
  //   title: ' 类型',
  //   dataIndex: 'action',
  // },

  // {
  //   title: ' 平台订单号',
  //   dataIndex: 'orderId',
  //   ellipsis: true,
  //     tooltip: true,
  // },
  // {
  //   title: ' 任务id',
  //   dataIndex: 'taskId',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  // {
  //   title: 'taskUrl',
  //   dataIndex: 'taskUrl',
  //   slotName: 'taskUrl',
  //   ellipsis: true,
  //     tooltip: true,
  // },
  // {
  //   title: 'taskUid',
  //   dataIndex: 'taskUid',
  //   slotName: 'taskUid',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  // {
  //   title: 'videoId',
  //   dataIndex: 'videoId',
  //   slotName: 'videoId',
  //   ellipsis: true,
  //     tooltip: true,
  // },

  {
    title: '审核备注',
    dataIndex: 'remark',
    slotName: 'remark',
  },
]);
const renderData = ref([]);

const basePagination: Pagination = { pageNum: 1, pageSize: 10, current: 1 };
const pagination = reactive({
  ...basePagination,
});
const fetchData = async (
  params: Pagination = { pageNum: 1, pageSize: 10 }
) => {
  setLoading(true);
  try {
    const { data } = await getMyTaskLogList(params);
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

// 表格分页
const onPageChange = (current: number) => {
  pagination.pageNum = current;
  fetchData(pagination);
};
</script>
<style scoped lang="less">
.container {
  padding: 0 20px 20px 20px;
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

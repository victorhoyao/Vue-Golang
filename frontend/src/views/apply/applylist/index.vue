<template>
  <div class="container">
    <a-card class="general-card" title="查询提现">
      <a-divider style="margin-top: 0" />
      <a-table row-key="id" :loading="loading" :columns="columns" :data="renderData" :bordered="{cell:true}"
        :pagination="pagination" @page-change="onPageChange">
        <template #status="{ record }">
          <a-tag v-if="record.status === 0" color="blue" bordered>
            {{ `审核中` }}</a-tag>
          <a-tag v-if="record.status === 1" color="green" bordered>
            {{ `审核中` }}</a-tag>
          <a-tag v-if="record.status === 2" bordered>
            {{ `审核中` }}</a-tag>
        </template>
        <template #operations="{ record }">
          <!-- <a-button type="text" size="small" status="success">
            {{ $t('searchTable.columns.operations.view') }}
          </a-button> -->
          <a-button type="text" status="success" size="small" @click="statusEdit(record, true)">
            审核成功
          </a-button>
          <a-button type="text" size="small" status="danger" @click="statusEdit(record, false)">
            审核失败
          </a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import type { Pagination } from '../../../api/applylist';
import { getApplyList, SubmitDoneTransaction } from '../../../api/applylist';
import useLoading from '../../../hooks/loading';


const { loading, setLoading } = useLoading();

type Column = TableColumnData & { checked?: true };

const columns = ref<Column[]>([
  {
    title: '编号',
    dataIndex: 'id',
    slotName: 'id',
  },
  {
    title: '用户id',
    dataIndex: 'userId',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: ' 提现订单号',
    dataIndex: 'orderId',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '提现积分',
    dataIndex: 'applyprice',
    ellipsis: true,
    tooltip: true,
  },

  // {
  //   title: '审核完成时间',
  //   dataIndex: 'doneTime',
  //   ellipsis: true,
  //   tooltip: true,
  // },
  {
    title: '提现申请时间',
    dataIndex: 'createTime',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '收款app',
    dataIndex: 'tranType',
    slotName: 'tranType',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '收款账号',
    dataIndex: 'tranAccount',
    slotName: 'tranAccount',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '收款姓名',
    dataIndex: 'tranName',
    slotName: 'tranName',
    ellipsis: true,
    tooltip: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    slotName: 'status',
    ellipsis: true,
    tooltip: true,
  },
  // {
  //   title: '结果说明',
  //   dataIndex: 'desc',
  //   slotName: 'desc',
  //   ellipsis: true,
  //   tooltip: true,
  // },
  {
    title: '操作',
    dataIndex: 'operations',
    slotName: 'operations',
    width: 250,

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
    const { data } = await getApplyList(params);
    console.log('用户列表数据', data);

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


const statusEdit = async (record: any, status: any) => {
  console.log(record.orderId);

  const result = {
    status: 2,
    desc: "失败联系上级管理员"
  }

  if (status) {
    result.status = 1
    result.desc = "成功"
  }
  const data = await SubmitDoneTransaction(record.id, result)
  if (data.code === 200) {
    Message.success({
      content: data.msg,
      duration: 3000
    });
  }
  fetchData();

}

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

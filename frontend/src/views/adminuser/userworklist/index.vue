<template>
    <div class="container">
        <a-card class="general-card" :title="$t('menu.list.searchTable')">
            <a-divider style="margin-top: 0" />

            <a-row style="margin-bottom: 16px">
                <a-col :span="6">
                    <a-space>
                        <a-date-picker style="width: 300px;" v-model="pickerDate" @change="pickerDateChange" />

                    </a-space>
                </a-col>
                <!-- <a-col :span="6">
                    <a-space>
                        <a-input :style="{ width: '320px' }" placeholder="请输入用户id" allow-clear />

                    </a-space>
                </a-col>

                <a-col :span="12">
                    <a-space>
                        <a-button type="primary">
                            <template #icon>
                                <icon-delete />
                            </template>
<template #default>搜索</template>
</a-button>
</a-space>
</a-col> -->
            </a-row>
            <a-table row-key="userId" :loading="loading" :columns="columns" :data="renderData" :bordered="{ cell: true }" column-resizable
                :pagination="pagination" @page-change="onPageChange"> <template #status="{ record }">
                    <a-tag v-if="record.status === 0" color="#00b42a">
                        {{ `正常` }}</a-tag>
                    <a-tag v-else color="#165dff">
                        {{ `关注` }}</a-tag>

                </template>
            </a-table>
        </a-card>
    </div>




</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';

import type { TableColumnData } from '@arco-design/web-vue/es/table/interface';
import type { UserListRecord, Pagination } from '../../../api/userlist';
import { getUserWorkList } from '../../../api/userlist';
import useLoading from '../../../hooks/loading';

const { loading, setLoading } = useLoading();

type Column = TableColumnData & { checked?: true };

const columns = ref<Column[]>([
    {
        title: "用户编号",
        dataIndex: 'userId',
        slotName: 'userId',
        ellipsis: true,
        tooltip: true,

    },
    {
        title: "用户名",
        dataIndex: 'userName',
        slotName: 'userName',
        ellipsis: true,
        tooltip: true,
    },
    {
        title: "统计日期",
        dataIndex: 'countDate',
        slotName: 'countDate',
        ellipsis: true,
        tooltip: true,
    },
    {
        title: "点赞",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'bliDiggSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'bliDiggFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "三连",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'bliSLSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'bliSLFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "投币",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blTBSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blTBFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "粉丝",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blfenSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blfenFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "会员购",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blhuiyuanGouSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blhuiyuanGouFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "播放",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blbofangSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blbofangFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "高速分享",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blgsfxSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blgsfxFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    },
    {
        title: "高速收藏",
        dataIndex: '',
        slotName: '',
        ellipsis: true,
        tooltip: true,
        children: [{
            title: '成功',
            dataIndex: 'blgsscSucessCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"green"
            }
        }, {
            title: '失败',
            dataIndex: 'blgsscFailCount',
            cellStyle:{
               color:"#ffffff",
                backgroundColor:"red"
            }
        }]
    }
]);
const renderData = ref<UserListRecord[]>([]);

const basePagination: Pagination = { pageNum: 1, pageSize: 10, current: 1, };
const pagination = reactive({
    ...basePagination,
});
const pickerDate = ref("")
const fetchData = async (
    params: Pagination = { pageNum: 1, pageSize: 10 }
) => {
    setLoading(true);

    try {
        const { data } = await getUserWorkList(params);
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
const pickerDateChange = async (value: any) => {
    console.log(value);
    pagination.countDate = value
    try {
        const { data } = await getUserWorkList(pagination);
        console.log('用户列表数据', data);
        renderData.value = data.result;
        // pagination.current = params.pageNum;
        pagination.total = data.total;
    } catch (err) {
        // you can report use errorHandler or other
    } finally {
        setLoading(false);
    }
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
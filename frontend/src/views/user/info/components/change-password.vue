<template>
  <div class="change-password-wrapper">
    <a-form
      ref="passwordFormRef"
      :model="passwordInfo"
      class="password-form"
      layout="vertical"
      @submit="handleChangePassword"
    >
      <a-form-item
        field="oldPassword"
        :rules="[{ required: true, message: '旧密码不能为空' }]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="passwordInfo.oldPassword"
          placeholder="请输入旧密码"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item
        field="newPassword"
        :rules="[
          { required: true, message: '新密码不能为空' },
          { min: 6, message: '新密码不能少于6位' },
        ]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="passwordInfo.newPassword"
          placeholder="请输入新密码"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <a-form-item
        field="confirmPassword"
        :rules="[
          { required: true, message: '请确认新密码' },
          {
            validator: (value: any, callback: any) => {
              if (value !== passwordInfo.newPassword) {
                callback('两次输入的密码不一致');
              } else {
                callback();
              }
            },
            message: '两次输入的密码不一致',
          },
        ]"
        :validate-trigger="['change', 'blur']"
        hide-label
      >
        <a-input-password
          v-model="passwordInfo.confirmPassword"
          placeholder="请再次输入新密码"
          allow-clear
        >
          <template #prefix>
            <icon-lock />
          </template>
        </a-input-password>
      </a-form-item>

      <a-space :size="16" direction="vertical">
        <a-button type="primary" html-type="submit" long :loading="loading">
          修改密码
        </a-button>
      </a-space>
    </a-form>
  </div>
</template>

<script lang="ts" setup>
import { reactive, ref } from 'vue';
import { Message } from '@arco-design/web-vue';
import useLoading from '@/hooks/loading';
import { changePassword } from '@/api/user';
import type { ChangePassData } from '@/api/user';

const passwordFormRef = ref();
const { loading, setLoading } = useLoading();

const passwordInfo = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: '',
});

const handleChangePassword = async ({ errors, values }: any) => {
  if (loading.value) return;
  if (!errors) {
    setLoading(true);
    try {
      const submitData: ChangePassData = {
        oldPass: values.oldPassword,
        newPass: values.newPassword,
      };
      const res: any = await changePassword(submitData);
      if (res.code === 200) {
        Message.success(res.msg || '密码修改成功');
        passwordFormRef.value.resetFields();
      } else {
        Message.error(res.msg || '密码修改失败，请重试');
      }
    } catch (err) {
      Message.error((err as Error).message || '密码修改失败，请重试');
    } finally {
      setLoading(false);
    }
  }
};
</script>

<style lang="less" scoped>
.change-password-wrapper {
  padding: 20px 0;
  width: 320px;
  margin: 0 auto;
}
</style> 
<script setup lang="ts">
import { ref, computed } from 'vue';
import { Modal, message } from 'ant-design-vue';
import PromptEditor from '#/components/PromptEditor/index.vue';
import { promptUpdate } from '#/api/ai/prompt';

const props = defineProps<{
  text?: string;
  row: any;
}>();

const emit = defineEmits(['success']);

const visible = ref(false);
const localText = ref('');
const loading = ref(false);

const displayText = computed(() => {
  const t = props.text || '';
  if (!t) return '点击编辑内容';
  return t.length > 100 ? t.slice(0, 100) + '...' : t;
});

const handleOpen = () => {
    localText.value = props.text || '';
    visible.value = true;
};

const handleOk = async () => {
    if (!props.row?.id) return;
    try {
        loading.value = true;
        // Merge current row data with new text to ensure all required fields are present for PUT
        await promptUpdate({
            ...props.row,
            text: localText.value,
        });
        
        message.success('更新成功');
        visible.value = false;
        
        // Update the row object directly for immediate visual feedback
        // The parent grid might reload, but this ensures no flickers
        props.row.text = localText.value;
        emit('success');
    } catch (error) {
        console.error(error);
    } finally {
        loading.value = false;
    }
};
</script>

<template>
  <div class="cursor-pointer group w-full" @click.stop="handleOpen">
    <div class="text-gray-500 group-hover:text-blue-500 transition-colors break-words whitespace-pre-wrap line-clamp-3" title="点击编辑">
        {{ displayText }}
    </div>
    
    <Modal
      v-model:open="visible"
      title="编辑提示词内容"
      width="900px"
      :confirmLoading="loading"
      @ok="handleOk"
      :bodyStyle="{ padding: '0' }"
      :destroyOnClose="true"
      :maskClosable="false"
      wrapClassName="prompt-preview-modal"
    >
       <div class="h-[550px] flex flex-col bg-gray-50 p-4">
          <PromptEditor v-model:value="localText" class="flex-1 shadow-sm" />
       </div>
    </Modal>
  </div>
</template>

<style scoped>
.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
